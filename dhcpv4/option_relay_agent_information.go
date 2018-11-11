package dhcpv4

import "fmt"

// This option implements the relay agent information option
// https://tools.ietf.org/html/rfc3046

// OptRelayAgentInformation is a "container" option for specific agent-supplied
// sub-options.
type OptRelayAgentInformation struct {
	Options []Option
}

// ParseOptRelayAgentInformation returns a new OptRelayAgentInformation from a
// byte stream, or error if any.
func ParseOptRelayAgentInformation(data []byte) (*OptRelayAgentInformation, error) {
	if len(data) < 4 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionRelayAgentInformation {
		return nil, fmt.Errorf("expected code %v, got %v", OptionRelayAgentInformation, code)
	}
	length := int(data[1])
	if len(data) < 2+length {
		return nil, ErrShortByteStream
	}
	options, err := OptionsFromBytesWithParser(data[2:length+2], relayParseOption)
	if err != nil {
		return nil, err
	}
	return &OptRelayAgentInformation{Options: options}, nil
}

func relayParseOption(data []byte) (Option, error) {
	if len(data) < 2 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	length := int(data[1])
	if len(data) < 2+length {
		return nil, ErrShortByteStream
	}
	return &OptionGeneric{OptionCode: code, Data: data[2:length+2]}, nil
}

// Code returns the option code.
func (o *OptRelayAgentInformation) Code() OptionCode {
	return OptionRelayAgentInformation
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptRelayAgentInformation) ToBytes() []byte {
	ret := []byte{byte(o.Code()), byte(o.Length())}
	for _, opt := range o.Options {
		ret = append(ret, opt.ToBytes()...)
	}
	return ret
}

// String returns a human-readable string for this option.
func (o *OptRelayAgentInformation) String() string {
	return fmt.Sprintf("Relay Agent Information -> [%v]", o.Options)
}

// Length returns the length of the data portion (excluding option code and byte
// for length, if any).
func (o *OptRelayAgentInformation) Length() int {
	l := 0
	for _, opt := range o.Options {
		l += 2 + opt.Length()
	}
	return l
}
