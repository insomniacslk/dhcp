package dhcpv4

import (
	"fmt"
	"net"
)

// This option implements the domain name server option
// https://tools.ietf.org/html/rfc2132

// OptDomainNameServer represents an option encapsulating the domain name
// servers.
type OptDomainNameServer struct {
	NameServers []net.IP
}

// ParseOptDomainNameServer returns a new OptDomainNameServer from a byte
// stream, or error if any.
func ParseOptDomainNameServer(data []byte) (*OptDomainNameServer, error) {
	ips, err := ParseIPs(data)
	if err != nil {
		return nil, err
	}
	return &OptDomainNameServer{NameServers: ips}, nil
}

// Code returns the option code.
func (o *OptDomainNameServer) Code() OptionCode {
	return OptionDomainNameServer
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptDomainNameServer) ToBytes() []byte {
	ret := []byte{byte(o.Code()), byte(o.Length())}
	for _, ns := range o.NameServers {
		ret = append(ret, ns.To4()...)
	}
	return ret
}

// String returns a human-readable string.
func (o *OptDomainNameServer) String() string {
	var servers string
	for idx, ns := range o.NameServers {
		servers += ns.String()
		if idx < len(o.NameServers)-1 {
			servers += ", "
		}
	}
	return fmt.Sprintf("Domain Name Servers -> %v", servers)
}

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptDomainNameServer) Length() int {
	return len(o.NameServers) * 4
}
