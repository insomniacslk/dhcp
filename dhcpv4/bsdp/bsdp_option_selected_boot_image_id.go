// +build darwin

package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Implements the BSDP option selected boot image ID, which tells the server
// which boot image has been selected by the client.

// OptSelectedBootImageID contains the selected boot image ID.
type OptSelectedBootImageID struct {
	ID BootImageID
}

// ParseOptSelectedBootImageID constructs an OptSelectedBootImageID struct from a sequence of
// bytes and returns it, or an error.
func ParseOptSelectedBootImageID(data []byte) (*OptSelectedBootImageID, error) {
	if len(data) < 6 {
		return nil, dhcpv4.ErrShortByteStream
	}
	code := dhcpv4.OptionCode(data[0])
	if code != OptionSelectedBootImageID {
		return nil, fmt.Errorf("expected option %v, got %v instead", OptionSelectedBootImageID, code)
	}
	length := int(data[1])
	if length != 4 {
		return nil, fmt.Errorf("expected length 4, got %d instead", length)
	}
	id, err := BootImageIDFromBytes(data[2:6])
	if err != nil {
		return nil, err
	}
	return &OptSelectedBootImageID{*id}, nil
}

// Code returns the option code.
func (o *OptSelectedBootImageID) Code() dhcpv4.OptionCode {
	return OptionSelectedBootImageID
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptSelectedBootImageID) ToBytes() []byte {
	serializedID := o.ID.ToBytes()
	return append([]byte{byte(o.Code()), byte(len(serializedID))}, serializedID...)
}

// String returns a human-readable string for this option.
func (o *OptSelectedBootImageID) String() string {
	return fmt.Sprintf("BSDP Selected Boot Image ID -> %s", o.ID.String())
}

// Length returns the length of the data portion of this option.
func (o *OptSelectedBootImageID) Length() int {
	return 4
}
