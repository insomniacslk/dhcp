// +build darwin

package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Implements the BSDP option listing the boot images.

// OptBootImageList contains the list of boot images presented by a netboot
// server.
type OptBootImageList struct {
	Images []BootImage
}

// ParseOptBootImageList constructs an OptBootImageList struct from a sequence
// of bytes and returns it, or an error.
func ParseOptBootImageList(data []byte) (*OptBootImageList, error) {
	// Should have at least code + length
	if len(data) < 2 {
		return nil, dhcpv4.ErrShortByteStream
	}
	code := dhcpv4.OptionCode(data[0])
	if code != OptionBootImageList {
		return nil, fmt.Errorf("expected option %v, got %v instead", OptionBootImageList, code)
	}
	length := int(data[1])
	if len(data) < length+2 {
		return nil, fmt.Errorf("expected length %d, got %d instead", length, len(data))
	}

	// Offset from code + length byte
	var bootImages []BootImage
	idx := 2
	for {
		if idx >= len(data) {
			break
		}
		image, err := BootImageFromBytes(data[idx:])
		if err != nil {
			return nil, fmt.Errorf("parsing bytes stream: %v", err)
		}
		bootImages = append(bootImages, *image)

		// 4 bytes of BootImageID, 1 byte of name length, name
		idx += 4 + 1 + len(image.Name)
	}

	return &OptBootImageList{bootImages}, nil
}

// Code returns the option code.
func (o *OptBootImageList) Code() dhcpv4.OptionCode {
	return OptionBootImageList
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptBootImageList) ToBytes() []byte {
	bs := make([]byte, 0, 2+o.Length())
	bs = append(bs, []byte{byte(o.Code()), byte(o.Length())}...)
	for _, image := range o.Images {
		bs = append(bs, image.ToBytes()...)
	}
	return bs
}

// String returns a human-readable string for this option.
func (o *OptBootImageList) String() string {
	s := "BSDP Boot Image List ->"
	for _, image := range o.Images {
		s += "\n  " + image.String()
	}
	return s
}

// Length returns the length of the data portion of this option.
func (o *OptBootImageList) Length() int {
	length := 0
	for _, image := range o.Images {
		length += 4 + 1 + len(image.Name)
	}
	return length
}
