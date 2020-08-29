package dhcpv6

import (
	"fmt"
)

// OptionGenericParseFailure represents an option that failed to be parsed correctly
type OptionGenericParseFailure struct {
	OptionCode OptionCode
	OptionData []byte
	Option     Option
	Error      error
}

// Code returns the option's code
func (og *OptionGenericParseFailure) Code() OptionCode {
	return og.OptionCode
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (og *OptionGenericParseFailure) ToBytes() []byte {
	return og.OptionData
}

func (og *OptionGenericParseFailure) String() string {
	return fmt.Sprintf("GenericParseFailure(%v): %v, Error=%v", og.OptionCode, og.Option, og.Error)
}
