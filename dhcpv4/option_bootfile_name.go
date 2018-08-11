package dhcpv4

import (
	"fmt"
)

// This option implements the Bootfile name Option.
// https://tools.ietf.org/html/rfc2132

// OptBootfileName implements the BootFile Name option
type OptBootfileName struct {
	BootfileName []byte
}

// Code returns the option code
func (op *OptBootfileName) Code() OptionCode {
	return OptionBootfileName
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptBootfileName) ToBytes() []byte {
	return append([]byte{byte(op.Code()), byte(op.Length())}, op.BootfileName...)
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
	if len(data) < 3 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionBootfileName {
		return nil, fmt.Errorf("ParseOptBootfileName: invalid code: %v; want %v", code, OptionBootfileName)
	}
	length := int(data[1])
	if length < 1 {
		return nil, fmt.Errorf("Bootfile name has invalid length of %d", length)
	}
	bootFileNameData := data[2:]
	if len(bootFileNameData) < length {
		return nil, fmt.Errorf("ParseOptBootfileName: short data: %d bytes; want %d",
			len(bootFileNameData), length)
	}
	return &OptBootfileName{BootfileName: bootFileNameData[:length]}, nil
}
