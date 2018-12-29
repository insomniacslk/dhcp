package dhcpv4

import "fmt"

// This option implements the domain name option
// https://tools.ietf.org/html/rfc2132

// OptDomainName represents an option encapsulating the domain name.
type OptDomainName struct {
	DomainName string
}

// ParseOptDomainName returns a new OptDomainName from a byte
// stream, or error if any.
func ParseOptDomainName(data []byte) (*OptDomainName, error) {
	return &OptDomainName{DomainName: string(data)}, nil
}

// Code returns the option code.
func (o *OptDomainName) Code() OptionCode {
	return OptionDomainName
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptDomainName) ToBytes() []byte {
	return []byte(o.DomainName)
}

// String returns a human-readable string.
func (o *OptDomainName) String() string {
	return fmt.Sprintf("Domain Name -> %v", o.DomainName)
}
