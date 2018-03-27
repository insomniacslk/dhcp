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
func ParseOptGeneric(data []byte) (*OptGeneric, error) {
	if len(data) == 0 {
		return nil, dhcpv4.ErrZeroLengthByteStream
	}
	var (
		length     int
		optionData []byte
	)
	code := dhcpv4.OptionCode(data[0])
	length = int(data[1])
	if len(data) < length+2 {
		return nil, fmt.Errorf("invalid data length: declared %v, actual %v", length, len(data))
	}
	optionData = data[2 : length+2]
	return &OptGeneric{OptionCode: code, Data: optionData}, nil
}

// Code returns the generic option code.
func (o OptGeneric) Code() dhcpv4.OptionCode {
	return o.OptionCode
}

// ToBytes returns a serialized generic option as a slice of bytes.
func (o OptGeneric) ToBytes() []byte {
	return append([]byte{byte(o.Code()), byte(o.Length())}, o.Data...)
}

// String returns a human-readable representation of a generic option.
func (o OptGeneric) String() string {
	code, ok := OptionCodeToString[o.Code()]
	if !ok {
		code = "Unknown"
	}
	return fmt.Sprintf("%v -> %v", code, o.Data)
}

// Length returns the number of bytes comprising the data section of the option.
func (o OptGeneric) Length() int {
	return len(o.Data)
}
