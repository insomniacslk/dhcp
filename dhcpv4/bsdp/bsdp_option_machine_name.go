// +build darwin

package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Implements the BSDP option machine name, which gives the Netboot server's
// machine name.

// OptMachineName represents a BSDP message type.
type OptMachineName struct {
	Name string
}

// ParseOptMachineName constructs an OptMachineName struct from a sequence of
// bytes and returns it, or an error.
func ParseOptMachineName(data []byte) (*OptMachineName, error) {
	if len(data) < 2 {
		return nil, dhcpv4.ErrShortByteStream
	}
	code := dhcpv4.OptionCode(data[0])
	if code != OptionMachineName {
		return nil, fmt.Errorf("expected option %v, got %v instead", OptionMachineName, code)
	}
	length := int(data[1])
	if len(data) < length+2 {
		return nil, fmt.Errorf("expected length %d, got %d instead", length, len(data))
	}
	return &OptMachineName{Name: string(data[2 : length+2])}, nil
}

// Code returns the option code.
func (o *OptMachineName) Code() dhcpv4.OptionCode {
	return OptionMachineName
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptMachineName) ToBytes() []byte {
	return append([]byte{byte(o.Code()), byte(o.Length())}, []byte(o.Name)...)
}

// String returns a human-readable string for this option.
func (o *OptMachineName) String() string {
	return "BSDP Machine Name -> " + o.Name
}

// Length returns the length of the data portion of this option.
func (o *OptMachineName) Length() int {
	return len(o.Name)
}
