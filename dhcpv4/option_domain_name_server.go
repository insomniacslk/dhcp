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
	return IPsToBytes(o.NameServers)
}

// String returns a human-readable string.
func (o *OptDomainNameServer) String() string {
	return fmt.Sprintf("Domain Name Servers -> %s", IPsToString(o.NameServers))
}
