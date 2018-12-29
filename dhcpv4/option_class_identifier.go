package dhcpv4

import (
	"fmt"
)

// This option implements the Class Identifier option
// https://tools.ietf.org/html/rfc2132

// OptClassIdentifier represents the DHCP message type option.
type OptClassIdentifier struct {
	Identifier string
}

// ParseOptClassIdentifier constructs an OptClassIdentifier struct from a sequence of
// bytes and returns it, or an error.
func ParseOptClassIdentifier(data []byte) (*OptClassIdentifier, error) {
	return &OptClassIdentifier{Identifier: string(data)}, nil
}

// Code returns the option code.
func (o *OptClassIdentifier) Code() OptionCode {
	return OptionClassIdentifier
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptClassIdentifier) ToBytes() []byte {
	return []byte(o.Identifier)
}

// String returns a human-readable string for this option.
func (o *OptClassIdentifier) String() string {
	return fmt.Sprintf("Class Identifier -> %v", o.Identifier)
}

// Length returns the length of the data portion (excluding option code and byte
// for length, if any).
func (o *OptClassIdentifier) Length() int {
	return len(o.Identifier)
}
