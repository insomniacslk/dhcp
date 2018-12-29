package bsdp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
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
	buf := uio.NewBigEndianBuffer(data)

	var bootImages []BootImage
	for buf.Has(5) {
		var image BootImage
		if err := (&image).Unmarshal(buf); err != nil {
			return nil, err
		}
		bootImages = append(bootImages, image)
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
