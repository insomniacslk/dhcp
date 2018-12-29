package dhcpv4

import (
	"fmt"
)

// This option implements the Bootfile name Option.
// https://tools.ietf.org/html/rfc2132

// OptBootfileName implements the BootFile Name option
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

// Length returns the option length in bytes
func (op *OptBootfileName) Length() int {
	return len(op.BootfileName)
}

func (op *OptBootfileName) String() string {
	return fmt.Sprintf("Bootfile Name -> %s", op.BootfileName)

}

// ParseOptBootfileName returns a new OptBootfile from a byte stream or error if any
func ParseOptBootfileName(data []byte) (*OptBootfileName, error) {
	return &OptBootfileName{BootfileName: string(data)}, nil
}
