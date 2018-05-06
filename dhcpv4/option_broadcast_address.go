package dhcpv4

import "fmt"

// This option implements the broadcast address option
// https://tools.ietf.org/html/rfc2132

// OptBroadcastAddress represents an option encapsulating the broadcast address.
type OptBroadcastAddress struct {
	OptGenericIP
}

// ParseOptBroadcastAddress returns a new OptBroadcastAddress from a byte
// stream, or error if any.
func ParseOptBroadcastAddress(data []byte) (*OptBroadcastAddress, error) {
	opt, err := ParseOptGenericIP(OptionBroadcastAddress, data)
	if err != nil {
		return nil, err
	}
	return &OptBroadcastAddress{OptGenericIP: *opt}, nil
}

// Code returns the option code.
func (o *OptBroadcastAddress) Code() OptionCode {
	return OptionBroadcastAddress
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptBroadcastAddress) ToBytes() []byte {
	return o.OptGenericIP.ToBytes(OptionBroadcastAddress)
}

// String returns a human-readable string.
func (o *OptBroadcastAddress) String() string {
	return fmt.Sprintf("Broadcast Address -> %v", o.IP.String())
}
