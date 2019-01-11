package dhcpv4

import (
	"fmt"

	"github.com/insomniacslk/dhcp/rfc1035label"
)

// OptDomainSearch implements the domain search list option described by RFC
// 3397, Section 2.
//
// FIXME: rename OptDomainSearch to OptDomainSearchList, and DomainSearch to
// SearchList, for consistency with the equivalent v6 option
type OptDomainSearch struct {
	DomainSearch *rfc1035label.Labels
}

// Code returns the option code.
func (op *OptDomainSearch) Code() OptionCode {
	return OptionDNSDomainSearchList
}

// ToBytes returns a serialized stream of bytes for this option.
func (op *OptDomainSearch) ToBytes() []byte {
	return op.DomainSearch.ToBytes()
}

// String returns a human-readable string.
func (op *OptDomainSearch) String() string {
	return fmt.Sprintf("DNS Domain Search List -> %v", op.DomainSearch.Labels)
}

// ParseOptDomainSearch returns a new OptDomainSearch from a byte stream, or
// error if any.
func ParseOptDomainSearch(data []byte) (*OptDomainSearch, error) {
	labels, err := rfc1035label.FromBytes(data)
	if err != nil {
		return nil, err
	}
	return &OptDomainSearch{DomainSearch: labels}, nil
}
