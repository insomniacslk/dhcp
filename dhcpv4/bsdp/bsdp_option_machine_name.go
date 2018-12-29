package bsdp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
)

// OptMachineName represents a BSDP message type.
//
// Implements the BSDP option machine name, which gives the Netboot server's
// machine name.
type OptMachineName struct {
	Name string
}

// ParseOptMachineName constructs an OptMachineName struct from a sequence of
// bytes and returns it, or an error.
func ParseOptMachineName(data []byte) (*OptMachineName, error) {
	return &OptMachineName{Name: string(data)}, nil
}

// Code returns the option code.
func (o *OptMachineName) Code() dhcpv4.OptionCode {
	return OptionMachineName
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptMachineName) ToBytes() []byte {
	return []byte(o.Name)
}

// String returns a human-readable string for this option.
func (o *OptMachineName) String() string {
	return "BSDP Machine Name -> " + o.Name
}

// Length returns the length of the data portion of this option.
func (o *OptMachineName) Length() int {
	return len(o.Name)
}
