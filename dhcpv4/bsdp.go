// +build darwin

package dhcpv4

// Implements Apple's netboot protocol BSDP (Boot Service Discovery Protocol).
// Canonical implementation is defined here:
// http://opensource.apple.com/source/bootp/bootp-198.1/Documentation/BSDP.doc

import (
	"encoding/binary"
	"fmt"
	"syscall"
)

// Options (occur as sub-options of DHCP option 43).
const (
	BSDPOptionMessageType OptionCode = iota + 1
	BSDPOptionVersion
	BSDPOptionServerIdentifier
	BSDPOptionServerPriority
	BSDPOptionReplyPort
	BSDPOptionBootImageListPath // Not used
	BSDPOptionDefaultBootImageID
	BSDPOptionSelectedBootImageID
	BSDPOptionBootImageList
	BSDPOptionNetboot1_0Firmware
	BSDPOptionBootImageAttributesFilterList
)

// Versions (seen so far)
var (
	BSDPVersion1_0 = []byte{1, 0}
	BSDPVersion1_1 = []byte{1, 1}
)

// BSDP message types
const (
	BSDPMessageTypeList byte = iota + 1
	BSDPMessageTypeSelect
	BSDPMessageTypeFailed
)

// Boot image kinds
const (
	BSDPBootImageMacOS9 byte = iota
	BSDPBootImageMacOSX
	BSDPBootImageMacOSXServer
	BSDPBootImageHardwareDiagnostics
	// 0x4 - 0x7f are reserved for future use.
)

// BootImageID describes a boot image ID - whether it's an install image and
// what kind of boot image (e.g. OS 9, macOS, hardware diagnostics)
type BootImageID struct {
	isInstall bool
	imageKind byte
	index     uint16
}

// toBytes serializes a BootImageID to network-order bytes.
func (b BootImageID) toBytes() (bytes []byte) {
	bytes = make([]byte, 4)
	// Attributes.
	if b.isInstall {
		bytes[0] |= 0x80
	}
	bytes[0] |= b.imageKind

	// Index
	binary.BigEndian.PutUint16(bytes[2:], b.index)
	return
}

// BootImageIDFromBytes deserializes a collection of 4 bytes to a BootImageID.
func bootImageIDFromBytes(bytes []byte) BootImageID {
	return BootImageID{
		isInstall: bytes[0]&0x80 != 0,
		imageKind: bytes[0] & 0x7f,
		index:     binary.BigEndian.Uint16(bytes[2:]),
	}
}

// BootImage describes a boot image - contains the boot image ID and the name.
type BootImage struct {
	ID BootImageID
	// This is a utf-8 string.
	Name string
}

// toBytes converts a BootImage to a slice of bytes.
func (b *BootImage) toBytes() (bytes []byte) {
	idBytes := b.ID.toBytes()
	bytes = append(bytes, idBytes[:]...)
	bytes = append(bytes, byte(len(b.Name)))
	bytes = append(bytes, []byte(b.Name)...)
	return
}

// BootImageFromBytes returns a deserialized BootImage struct from bytes as well
// as the number of bytes read from the slice.
func bootImageFromBytes(bytes []byte) (*BootImage, int, error) {
	// If less than length of boot image ID and count, it's probably invalid.
	if len(bytes) < 5 {
		return nil, 0, fmt.Errorf("not enough bytes for BootImage")
	}
	imageID := bootImageIDFromBytes(bytes[:4])
	nameLength := int(bytes[4])
	if 5+nameLength > len(bytes) {
		return nil, 0, fmt.Errorf("not enough bytes for BootImage")
	}
	name := string(bytes[5 : 5+nameLength])
	return &BootImage{ID: imageID, Name: name}, 5 + nameLength, nil
}

// makeVendorClassIdentifier calls the sysctl syscall on macOS to get the
// platform model.
func makeVendorClassIdentifier() (string, error) {
	// Fetch hardware model for class ID.
	hwModel, err := syscall.Sysctl("hw.model")
	if err != nil {
		return "", err
	}
	vendorClassID := fmt.Sprintf("AAPLBSDPC/i386/%s", hwModel)
	return vendorClassID, nil
}

// parseBootImagesFromBSDPOption parses data from the BSDPOptionBootImageList
// option and returns a list of BootImages.
func parseBootImagesFromBSDPOption(data []byte) ([]BootImage, error) {
	// Should at least have the # bytes of boot images.
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid length boot image list")
	}

	readByteCount := 0
	start := data
	var bootImages []BootImage
	for {
		bootImage, readBytes, err := bootImageFromBytes(start)
		if err != nil {
			return nil, err
		}
		bootImages = append(bootImages, *bootImage)
		readByteCount += readBytes
		if readByteCount+1 >= len(data) {
			break
		}
		start = start[readByteCount:]
	}

	return bootImages, nil
}

// parseVendorOptionsFromOptions extracts the sub-options list of the vendor-
// specific options from the larger DHCP options list.
func parseVendorOptionsFromOptions(options []Option) []Option {
	var vendorOpts []Option
	var err error
	for _, opt := range options {
		if opt.Code == OptionVendorSpecificInformation {
			vendorOpts, err = OptionsFromBytes(opt.Data)
			if err != nil {
				return []Option{}
			}
			break
		}
	}
	return vendorOpts
}

// ParseBootImageListFromAck parses the list of boot images presented in the
// ACK[LIST] packet and returns them as a list of BootImages.
func ParseBootImageListFromAck(ack DHCPv4) ([]BootImage, error) {
	var bootImages []BootImage
	vendorOpts := parseVendorOptionsFromOptions(ack.options)
	for _, opt := range vendorOpts {
		if opt.Code == BSDPOptionBootImageList {
			images, err := parseBootImagesFromBSDPOption(opt.Data)
			if err != nil {
				return nil, err
			}
			bootImages = append(bootImages, images...)
		}
	}

	return bootImages, nil
}

// NewInformListForInterface creates a new INFORM packet for interface ifname
// with configuration options specified by config.
func NewInformListForInterface(iface string, replyPort uint16) (*DHCPv4, error) {
	d, err := NewInformForInterface(iface /* needsBroadcast */, false)
	if err != nil {
		return nil, err
	}

	// These are vendor-specific options used to pass along BSDP information.
	vendorOpts := []Option{
		Option{
			Code: BSDPOptionMessageType,
			Data: []byte{BSDPMessageTypeList},
		},
		Option{
			Code: BSDPOptionVersion,
			Data: BSDPVersion1_1,
		},
	}

	// If specified, replyPort MUST be a priviledged port.
	if replyPort != 0 && replyPort != ClientPort {
		if replyPort >= 1024 {
			return nil, fmt.Errorf("replyPort must be a priviledged port (< 1024)")
		}
		bytes := make([]byte, 3)
		bytes[0] = 2
		binary.BigEndian.PutUint16(bytes[1:], replyPort)
		d.AddOption(Option{
			Code: BSDPOptionReplyPort,
			Data: bytes,
		})
	}
	d.AddOption(Option{
		Code: OptionVendorSpecificInformation,
		Data: OptionsToBytes(vendorOpts),
	})

	d.AddOption(Option{
		Code: OptionParameterRequestList,
		Data: []byte{OptionVendorSpecificInformation, OptionClassIdentifier},
	})

	u16 := make([]byte, 2)
	binary.BigEndian.PutUint16(u16, 1500)
	d.AddOption(Option{
		Code: OptionMaximumDHCPMessageSize,
		Data: u16,
	})

	vendorClassID, err := makeVendorClassIdentifier()
	if err != nil {
		return nil, err
	}
	d.AddOption(Option{
		Code: OptionClassIdentifier,
		Data: []byte(vendorClassID),
	})

	d.AddOption(Option{Code: OptionEnd})
	return d, nil
}

// InformSelectForAck constructs an INFORM[SELECT] packet given an ACK to the
// previously-sent INFORM[LIST] with BSDPConfig config.
func InformSelectForAck(ack DHCPv4, replyPort uint16, selectedImage BootImage) (*DHCPv4, error) {
	d, err := New()
	if err != nil {
		return nil, err
	}
	d.SetOpcode(OpcodeBootRequest)
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

	// Data for BSDPOptionSelectedBootImageID
	vendorOpts := []Option{
		Option{
			Code: BSDPOptionMessageType,
			Data: []byte{BSDPMessageTypeSelect},
		},
		Option{
			Code: BSDPOptionVersion,
			Data: BSDPVersion1_1,
		},
		Option{
			Code: BSDPOptionSelectedBootImageID,
			Data: append([]byte{4}, selectedImage.ID.toBytes()...),
		},
	}

	// Find server IP address
	var serverIP []byte
	for _, opt := range ack.options {
		if opt.Code == OptionServerIdentifier {
			serverIP = make([]byte, 4)
			copy(serverIP, opt.Data)
		}
	}
	if len(serverIP) == 0 {
		return nil, fmt.Errorf("could not parse server identifier from ACK")
	}
	vendorOpts = append(vendorOpts, Option{
		Code: BSDPOptionServerIdentifier,
		Data: serverIP,
	})

	// Validate replyPort if requested.
	if replyPort != 0 && replyPort != ClientPort {
		// replyPort MUST be a priviledged port.
		if replyPort >= 1024 {
			return nil, fmt.Errorf("replyPort must be a priviledged port")
		}
		bytes := make([]byte, 3)
		bytes[0] = 2
		binary.BigEndian.PutUint16(bytes[1:], replyPort)
		vendorOpts = append(vendorOpts, Option{
			Code: BSDPOptionReplyPort,
			Data: bytes,
		})
	}

	d.AddOption(Option{
		Code: OptionDHCPMessageType,
		Data: []byte{MessageTypeInform},
	})
	d.AddOption(Option{
		Code: OptionVendorSpecificInformation,
		Data: OptionsToBytes(vendorOpts),
	})
	d.AddOption(Option{Code: OptionEnd})
	return d, nil
}
