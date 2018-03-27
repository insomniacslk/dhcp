package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Implements the BSDP option default boot image ID, which tells the client
// which image is the default boot image if one is not selected.

// OptDefaultBootImageID contains the selected boot image ID.
type OptDefaultBootImageID struct {
	ID BootImageID
}

// ParseOptDefaultBootImageID constructs an OptDefaultBootImageID struct from a sequence of
// bytes and returns it, or an error.
func ParseOptDefaultBootImageID(data []byte) (*OptDefaultBootImageID, error) {
	if len(data) < 6 {
		return nil, dhcpv4.ErrShortByteStream
	}
	code := dhcpv4.OptionCode(data[0])
	if code != OptionDefaultBootImageID {
		return nil, fmt.Errorf("expected option %v, got %v instead", OptionDefaultBootImageID, code)
	}
	length := int(data[1])
	if length != 4 {
		return nil, fmt.Errorf("expected length 4, got %d instead", length)
	}
	id, err := BootImageIDFromBytes(data[2:6])
	if err != nil {
		return nil, err
	}
	return &OptDefaultBootImageID{*id}, nil
}

// Code returns the option code.
func (o *OptDefaultBootImageID) Code() dhcpv4.OptionCode {
	return OptionDefaultBootImageID
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptDefaultBootImageID) ToBytes() []byte {
	serializedID := o.ID.ToBytes()
	return append([]byte{byte(o.Code()), byte(len(serializedID))}, serializedID...)
}

// String returns a human-readable string for this option.
func (o *OptDefaultBootImageID) String() string {
	return fmt.Sprintf("BSDP Default Boot Image ID -> %s", o.ID.String())
}

// Length returns the length of the data portion of this option.
func (o *OptDefaultBootImageID) Length() int {
	return 4
}
