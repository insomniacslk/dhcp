// +build darwin

package bsdp

// Implements Apple's netboot protocol BSDP (Boot Service Discovery Protocol).
// Canonical implementation is defined here:
// http://opensource.apple.com/source/bootp/bootp-198.1/Documentation/BSDP.doc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// MaxDHCPMessageSize is the size set in DHCP option 57 (DHCP Maximum Message Size).
// BSDP includes its own sub-option (12) to indicate to NetBoot servers that the
// client can support larger message sizes, and modern NetBoot servers will
// prefer this BSDP-specific option over the DHCP standard option.
const MaxDHCPMessageSize = 1500

// makeVendorClassIdentifier calls the sysctl syscall on macOS to get the
// platform model.
func makeVendorClassIdentifier() (string, error) {
	// Fetch hardware model for class ID.
	hwModel, err := syscall.Sysctl("hw.model")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("AAPLBSDPC/i386/%s", hwModel), nil
}

// ParseBootImagesFromOption parses data from the BSDPOptionBootImageList
// option and returns a list of BootImages.
func ParseBootImagesFromOption(data []byte) ([]BootImage, error) {
	// Should at least have the # bytes of boot images.
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid length boot image list")
	}

	var (
		readByteCount = 0
		start         = data
		bootImages    []BootImage
	)
	for {
		bootImage, err := BootImageFromBytes(start)
		if err != nil {
			return nil, err
		}
		bootImages = append(bootImages, *bootImage)
		// Read BootImageID + name length + name
		readByteCount += 4 + 1 + len(bootImage.Name)
		if readByteCount+1 >= len(data) {
			break
		}
		start = start[readByteCount:]
	}

	return bootImages, nil
}

// ParseVendorOptionsFromOptions extracts the sub-options list of the vendor-
// specific options from the larger DHCP options list.
// TODO: Implement options.GetOneOption for dhcpv4.
func ParseVendorOptionsFromOptions(options []dhcpv4.Option) []dhcpv4.Option {
	var (
		vendorOpts []dhcpv4.Option
		err        error
	)
	for _, opt := range options {
		if opt.Code() == dhcpv4.OptionVendorSpecificInformation {
			vendorOpts, err = dhcpv4.OptionsFromBytesWithoutMagicCookie(opt.(*dhcpv4.OptionGeneric).Data)
			if err != nil {
				log.Println("Warning: could not parse vendor options in DHCP options")
				return []dhcpv4.Option{}
			}
			break
		}
	}
	return vendorOpts
}

// ParseBootImageListFromAck parses the list of boot images presented in the
// ACK[LIST] packet and returns them as a list of BootImages.
func ParseBootImageListFromAck(ack dhcpv4.DHCPv4) ([]BootImage, error) {
	var bootImages []BootImage
	for _, opt := range ParseVendorOptionsFromOptions(ack.Options()) {
		if opt.Code() == OptionBootImageList {
			images, err := ParseBootImagesFromOption(opt.(*dhcpv4.OptionGeneric).Data)
			if err != nil {
				return nil, err
			}
			bootImages = append(bootImages, images...)
		}
	}

	return bootImages, nil
}

func needsReplyPort(replyPort uint16) bool {
	return replyPort != 0 && replyPort != dhcpv4.ClientPort
}

// NewInformListForInterface creates a new INFORM packet for interface ifname
// with configuration options specified by config.
func NewInformListForInterface(iface string, replyPort uint16) (*dhcpv4.DHCPv4, error) {
	d, err := dhcpv4.NewInformForInterface(iface /* needsBroadcast = */, false)
	if err != nil {
		return nil, err
	}

	// Validate replyPort first
	if needsReplyPort(replyPort) && replyPort >= 1024 {
		return nil, errors.New("replyPort must be a privileged port")
	}

	// These are vendor-specific options used to pass along BSDP information.
	vendorOpts := []dhcpv4.Option{
		&OptMessageType{MessageTypeList},
		&OptVersion{Version1_1},
	}

	if needsReplyPort(replyPort) {
		vendorOpts = append(vendorOpts, &OptReplyPort{replyPort})
	}
	var vendorOptsBytes []byte
	for _, opt := range vendorOpts {
		vendorOptsBytes = append(vendorOptsBytes, opt.ToBytes()...)
	}
	d.AddOption(&dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionVendorSpecificInformation,
		Data:       vendorOptsBytes,
	})

	d.AddOption(&dhcpv4.OptParameterRequestList{
		RequestedOpts: []dhcpv4.OptionCode{
			dhcpv4.OptionVendorSpecificInformation,
			dhcpv4.OptionClassIdentifier,
		},
	})
	d.AddOption(&dhcpv4.OptMaximumDHCPMessageSize{Size: MaxDHCPMessageSize})

	vendorClassID, err := makeVendorClassIdentifier()
	if err != nil {
		return nil, err
	}
	d.AddOption(&dhcpv4.OptClassIdentifier{vendorClassID})
	d.AddOption(&dhcpv4.OptionGeneric{OptionCode: dhcpv4.OptionEnd})
	return d, nil
}

// InformSelectForAck constructs an INFORM[SELECT] packet given an ACK to the
// previously-sent INFORM[LIST] with Config config.
func InformSelectForAck(ack dhcpv4.DHCPv4, replyPort uint16, selectedImage BootImage) (*dhcpv4.DHCPv4, error) {
	d, err := dhcpv4.New()
	if err != nil {
		return nil, err
	}

	if needsReplyPort(replyPort) && replyPort >= 1024 {
		return nil, errors.New("replyPort must be a privilegded port")
	}
	d.SetOpcode(dhcpv4.OpcodeBootRequest)
	d.SetHwType(ack.HwType())
	d.SetHwAddrLen(ack.HwAddrLen())
	clientHwAddr := ack.ClientHwAddr()
	d.SetClientHwAddr(clientHwAddr[:])
	d.SetTransactionID(ack.TransactionID())
	if ack.IsBroadcast() {
		d.SetBroadcast()
	} else {
		d.SetUnicast()
	}

	// Data for OptionSelectedBootImageID
	vendorOpts := []dhcpv4.Option{
		&OptMessageType{MessageTypeSelect},
		&OptVersion{Version1_1},
		&OptSelectedBootImageID{selectedImage.ID},
	}

	// Find server IP address
	var serverIP net.IP
	// TODO replace this loop with `ack.GetOneOption(OptionBootImageList)`
	for _, opt := range ack.Options() {
		if opt.Code() == dhcpv4.OptionServerIdentifier {
			serverIP = net.IP(opt.(*dhcpv4.OptionGeneric).Data)
		}
	}
	if serverIP.To4() == nil {
		return nil, fmt.Errorf("could not parse server identifier from ACK")
	}
	vendorOpts = append(vendorOpts, &dhcpv4.OptServerIdentifier{ServerID: serverIP})

	// Validate replyPort if requested.
	if needsReplyPort(replyPort) {
		vendorOpts = append(vendorOpts, &OptReplyPort{replyPort})
	}

	vendorClassID, err := makeVendorClassIdentifier()
	if err != nil {
		return nil, err
	}
	d.AddOption(&dhcpv4.OptClassIdentifier{vendorClassID})
	d.AddOption(&dhcpv4.OptParameterRequestList{
		[]dhcpv4.OptionCode{
			dhcpv4.OptionSubnetMask,
			dhcpv4.OptionRouter,
			dhcpv4.OptionBootfileName,
			dhcpv4.OptionVendorSpecificInformation,
			dhcpv4.OptionClassIdentifier,
		},
	})
	d.AddOption(&dhcpv4.OptMessageType{dhcpv4.MessageTypeInform})
	var vendorOptsBytes []byte
	for _, opt := range vendorOpts {
		vendorOptsBytes = append(vendorOptsBytes, opt.ToBytes()...)
	}
	d.AddOption(&dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionVendorSpecificInformation,
		Data:       vendorOptsBytes,
	})
	d.AddOption(&dhcpv4.OptionGeneric{OptionCode: dhcpv4.OptionEnd})
	return d, nil
}
