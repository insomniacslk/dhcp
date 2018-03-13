// +build darwin

package bsdp

// Implements Apple's netboot protocol BSDP (Boot Service Discovery Protocol).
// Canonical implementation is defined here:
// http://opensource.apple.com/source/bootp/bootp-198.1/Documentation/BSDP.doc

import (
	"encoding/binary"
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

// BootImageID describes a boot image ID - whether it's an install image and
// what kind of boot image (e.g. OS 9, macOS, hardware diagnostics)
type BootImageID struct {
	IsInstall bool
	ImageType BootImageType
	Index     uint16
}

// ToBytes serializes a BootImageID to network-order bytes.
func (b BootImageID) ToBytes() []byte {
	bytes := make([]byte, 4)
	if b.IsInstall {
		bytes[0] |= 0x80
	}
	bytes[0] |= byte(b.ImageType)
	binary.BigEndian.PutUint16(bytes[2:], b.Index)
	return bytes
}

// BootImageIDFromBytes deserializes a collection of 4 bytes to a BootImageID.
func BootImageIDFromBytes(bytes []byte) (*BootImageID, error) {
	if len(bytes) < 4 {
		return nil, fmt.Errorf("not enough bytes to serialize BootImageID")
	}
	return &BootImageID{
		IsInstall: bytes[0]&0x80 != 0,
		ImageType: BootImageType(bytes[0] & 0x7f),
		Index:     binary.BigEndian.Uint16(bytes[2:]),
	}, nil
}

// BootImage describes a boot image - contains the boot image ID and the name.
type BootImage struct {
	ID   BootImageID
	Name string
}

// ToBytes converts a BootImage to a slice of bytes.
func (b *BootImage) ToBytes() []byte {
	bytes := b.ID.ToBytes()
	bytes = append(bytes, byte(len(b.Name)))
	bytes = append(bytes, []byte(b.Name)...)
	return bytes
}

// BootImageFromBytes returns a deserialized BootImage struct from bytes.
func BootImageFromBytes(bytes []byte) (*BootImage, error) {
	// Should at least contain 4 bytes of BootImageID + byte for length of
	// boot image name.
	if len(bytes) < 5 {
		return nil, fmt.Errorf("not enough bytes to serialize BootImage")
	}
	imageID, err := BootImageIDFromBytes(bytes[:4])
	if err != nil {
		return nil, err
	}
	nameLength := int(bytes[4])
	if 5+nameLength > len(bytes) {
		return nil, fmt.Errorf("not enough bytes for BootImage")
	}
	name := string(bytes[5 : 5+nameLength])
	return &BootImage{ID: *imageID, Name: name}, nil
}

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

func serializeReplyPort(replyPort uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, replyPort)
	return bytes
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
		dhcpv4.OptionGeneric{
			OptionCode: OptionMessageType,
			Data:       []byte{byte(MessageTypeList)},
		},
		dhcpv4.OptionGeneric{
			OptionCode: OptionVersion,
			Data:       Version1_1,
		},
	}

	if needsReplyPort(replyPort) {
		vendorOpts = append(vendorOpts,
			dhcpv4.OptionGeneric{
				OptionCode: OptionReplyPort,
				Data:       serializeReplyPort(replyPort),
			},
		)
	}
	var vendorOptsBytes []byte
	for _, opt := range vendorOpts {
		vendorOptsBytes = append(vendorOptsBytes, opt.ToBytes()...)
	}
	d.AddOption(dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionVendorSpecificInformation,
		Data:       vendorOptsBytes,
	})

	d.AddOption(dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionParameterRequestList,
		Data: []byte{
			byte(dhcpv4.OptionVendorSpecificInformation),
			byte(dhcpv4.OptionClassIdentifier),
		},
	})

	u16 := make([]byte, 2)
	binary.BigEndian.PutUint16(u16, MaxDHCPMessageSize)
	d.AddOption(dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionMaximumDHCPMessageSize,
		Data:       u16,
	})

	vendorClassID, err := makeVendorClassIdentifier()
	if err != nil {
		return nil, err
	}
	d.AddOption(dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionClassIdentifier,
		Data:       []byte(vendorClassID),
	})

	d.AddOption(dhcpv4.OptionGeneric{OptionCode: dhcpv4.OptionEnd})
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
		dhcpv4.OptionGeneric{
			OptionCode: OptionMessageType,
			Data:       []byte{byte(MessageTypeSelect)},
		},
		dhcpv4.OptionGeneric{
			OptionCode: OptionVersion,
			Data:       Version1_1,
		},
		dhcpv4.OptionGeneric{
			OptionCode: OptionSelectedBootImageID,
			Data:       selectedImage.ID.ToBytes(),
		},
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
	vendorOpts = append(vendorOpts, dhcpv4.OptionGeneric{
		OptionCode: OptionServerIdentifier,
		Data:       serverIP,
	})

	// Validate replyPort if requested.
	if needsReplyPort(replyPort) {
		vendorOpts = append(vendorOpts, dhcpv4.OptionGeneric{
			OptionCode: OptionReplyPort,
			Data:       serializeReplyPort(replyPort),
		})
	}

	vendorClassID, err := makeVendorClassIdentifier()
	if err != nil {
		return nil, err
	}
	d.AddOption(dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionClassIdentifier,
		Data:       []byte(vendorClassID),
	})
	d.AddOption(dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionParameterRequestList,
		Data: []byte{
			byte(dhcpv4.OptionSubnetMask),
			byte(dhcpv4.OptionRouter),
			byte(dhcpv4.OptionBootfileName),
			byte(dhcpv4.OptionVendorSpecificInformation),
			byte(dhcpv4.OptionClassIdentifier),
		},
	})
	d.AddOption(dhcpv4.NewOptMessageType(dhcpv4.MessageTypeInform))
	var vendorOptsBytes []byte
	for _, opt := range vendorOpts {
		vendorOptsBytes = append(vendorOptsBytes, opt.ToBytes()...)
	}
	d.AddOption(dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionVendorSpecificInformation,
		Data:       vendorOptsBytes,
	})
	d.AddOption(dhcpv4.OptionGeneric{OptionCode: dhcpv4.OptionEnd})
	return d, nil
}
