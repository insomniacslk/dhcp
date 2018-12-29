package dhcpv4

import "fmt"

// This option implements the host name option
// https://tools.ietf.org/html/rfc2132

// OptHostName represents an option encapsulating the host name.
type OptHostName struct {
	HostName string
}

// ParseOptHostName returns a new OptHostName from a byte stream, or error if
// any.
func ParseOptHostName(data []byte) (*OptHostName, error) {
	return &OptHostName{HostName: string(data)}, nil
}

// Code returns the option code.
func (o *OptHostName) Code() OptionCode {
	return OptionHostName
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptHostName) ToBytes() []byte {
	return []byte(o.HostName)
}

// String returns a human-readable string.
func (o *OptHostName) String() string {
	return fmt.Sprintf("Host Name -> %v", o.HostName)
}
