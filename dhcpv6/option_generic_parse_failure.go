package dhcpv6

import (
	"fmt"
)

type OptionGenericParseFailure struct {
	OptionCode OptionCode
	OptionData []byte
	Option     Option
	Error      error
}

func (og *OptionGenericParseFailure) Code() OptionCode {
	return og.OptionCode
}

func (og *OptionGenericParseFailure) ToBytes() []byte {
	return og.OptionData
}

func (og *OptionGenericParseFailure) String() string {
	return fmt.Sprintf("GenericParseFailure(%v): %v, Error=%v", og.OptionCode, og.Option, og.Error)
}
