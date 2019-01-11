package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// OptGeneric is an option that only contains the option code and associated
// data. Every option that does not have a specific implementation will fall
// back to this option.
type OptGeneric struct {
	OptionCode dhcpv4.OptionCode
	Data       []byte
}

// ParseOptGeneric parses a bytestream and creates a new OptGeneric from it,
// or an error.
func ParseOptGeneric(code dhcpv4.OptionCode, data []byte) (*OptGeneric, error) {
	return &OptGeneric{OptionCode: code, Data: data}, nil
}

// Code returns the generic option code.
func (o OptGeneric) Code() dhcpv4.OptionCode {
	return o.OptionCode
}

// ToBytes returns a serialized generic option as a slice of bytes.
func (o OptGeneric) ToBytes() []byte {
	return o.Data
}

// String returns a human-readable representation of a generic option.
func (o OptGeneric) String() string {
	return fmt.Sprintf("%s -> %v", o.OptionCode, o.Data)
}
