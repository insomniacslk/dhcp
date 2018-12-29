package dhcpv4

import (
	"errors"
	"fmt"
)

// OptionGeneric is an option that only contains the option code and associated
// data. Every option that does not have a specific implementation will fall
// back to this option.
type OptionGeneric struct {
	OptionCode OptionCode
	Data       []byte
}

// ParseOptionGeneric parses a bytestream and creates a new OptionGeneric from
// it, or an error.
func ParseOptionGeneric(code OptionCode, data []byte) (Option, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid zero-length bytestream")
	}
	return &OptionGeneric{OptionCode: code, Data: data}, nil
}

// Code returns the generic option code.
func (o OptionGeneric) Code() OptionCode {
	return o.OptionCode
}

// ToBytes returns a serialized generic option as a slice of bytes.
func (o OptionGeneric) ToBytes() []byte {
	return o.Data
}

// String returns a human-readable representation of a generic option.
func (o OptionGeneric) String() string {
	return fmt.Sprintf("%v -> %v", o.OptionCode.String(), o.Data)
}

// Length returns the number of bytes comprising the data section of the option.
func (o OptionGeneric) Length() int {
	return len(o.Data)
}
