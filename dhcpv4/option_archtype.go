package dhcpv4

import (
	"github.com/insomniacslk/dhcp/iana"
)

// OptClientArch returns a new Client System Architecture Type option.
func OptClientArch(archs ...iana.Arch) Option {
	return Option{Code: OptionClientSystemArchitectureType, Value: iana.Archs(archs)}
}

// GetClientArch returns the Client System Architecture Type option.
func GetClientArch(o Options) []iana.Arch {
	v := o.Get(OptionClientSystemArchitectureType)
	if v == nil {
		return nil
	}
	var archs iana.Archs
	if err := archs.FromBytes(v); err != nil {
		return nil
	}
	return archs
}
