package dhcpv4

import (
	"fmt"
)

// OptBootfileName implements the bootfile name option described in RFC 2132,
// Section 9.5.
type OptBootfileName struct {
	BootfileName string
}

// Code returns the option code
func (op *OptBootfileName) Code() OptionCode {
	return OptionBootfileName
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptBootfileName) ToBytes() []byte {
	return []byte(op.BootfileName)
}

func (op *OptBootfileName) String() string {
	return fmt.Sprintf("Bootfile Name -> %s", op.BootfileName)
}

// ParseOptBootfileName returns a new OptBootfile from a byte stream or error if any
func ParseOptBootfileName(data []byte) (*OptBootfileName, error) {
	return &OptBootfileName{BootfileName: string(data)}, nil
}
