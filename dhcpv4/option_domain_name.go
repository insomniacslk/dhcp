package dhcpv4

import "fmt"

// OptDomainName implements the domain name option described in RFC 2132,
// Section 3.17.
type OptDomainName struct {
	DomainName string
}

// ParseOptDomainName returns a new OptDomainName from a byte stream, or error
// if any.
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
