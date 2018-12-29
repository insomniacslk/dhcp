package dhcpv4

import (
	"github.com/insomniacslk/dhcp/rfc1035label"
)

// OptDomainSearch returns a new domain search option.
//
// The domain search option is described by RFC 3397, Section 2.
func OptDomainSearch(labels *rfc1035label.Labels) Option {
	return Option{Code: OptionDNSDomainSearchList, Value: labels}
}

// GetDomainSearch returns the domain search list in o, if present.
//
// The domain search option is described by RFC 3397, Section 2.
func GetDomainSearch(o Options) *rfc1035label.Labels {
	v := o.Get(OptionDNSDomainSearchList)
	if v == nil {
		return nil
	}
	labels, err := rfc1035label.FromBytes(v)
	if err != nil {
		return nil
	}
	return labels
}
