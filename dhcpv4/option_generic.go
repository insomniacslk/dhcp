package dhcpv4

import (
	"fmt"
)

// OptionGeneric is an option that only contains the option code and associated
// data. Every option that does not have a specific implementation will fall
// back to this option.
type OptionGeneric struct {
	OptionCode OptionCode
	Data       []byte
}

// Code returns the generic option code.
func (o OptionGeneric) Code() OptionCode {
	return o.OptionCode
}

// ToBytes returns a serialized generic option as a slice of bytes.
func (o OptionGeneric) ToBytes() []byte {
	ret := []byte{byte(o.OptionCode)}
	if o.OptionCode == OptionEnd || o.OptionCode == OptionPad {
		return ret
	}
	ret = append(ret, byte(o.Length()))
	ret = append(ret, o.Data...)
	return ret
}

// String returns a human-readable representation of a generic option.
func (o OptionGeneric) String() string {
	code, ok := OptionCodeToString[o.OptionCode]
	if !ok {
		code = "Unknown"
	}
	return fmt.Sprintf("%v -> %v", code, o.Data)
}

// Length returns the number of bytes comprising the data section of the option.
func (o OptionGeneric) Length() int {
	return len(o.Data)
}
