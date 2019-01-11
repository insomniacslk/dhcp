package dhcpv4

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptRelayAgentInformation implements the relay agent info option described by
// RFC 3046.
type OptRelayAgentInformation struct {
	Options Options
}

// ParseOptRelayAgentInformation returns a new OptRelayAgentInformation from a
// byte stream, or error if any.
func ParseOptRelayAgentInformation(data []byte) (*OptRelayAgentInformation, error) {
	options, err := OptionsFromBytesWithParser(data, codeGetter, ParseOptionGeneric, false /* don't check for OptionEnd tag */)
	if err != nil {
		return nil, err
	}
	return &OptRelayAgentInformation{Options: options}, nil
}

// Code returns the option code.
func (o *OptRelayAgentInformation) Code() OptionCode {
	return OptionRelayAgentInformation
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptRelayAgentInformation) ToBytes() []byte {
	return uio.ToBigEndian(o.Options)
}

// String returns a human-readable string for this option.
func (o *OptRelayAgentInformation) String() string {
	return fmt.Sprintf("Relay Agent Information -> %v", o.Options)
}
