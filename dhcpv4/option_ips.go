package dhcpv4

import (
	"fmt"
	"net"
	"strings"

	"github.com/u-root/u-root/pkg/uio"
)

// ParseIPs parses an IPv4 address from a DHCP packet as used and specified by
// options in RFC 2132, Sections 3.5 through 3.13, 8.2, 8.3, 8.5, 8.6, 8.9, and
// 8.10.
func ParseIPs(data []byte) ([]net.IP, error) {
	buf := uio.NewBigEndianBuffer(data)

	if buf.Len() == 0 {
		return nil, fmt.Errorf("IP DHCP options must always list at least one IP")
	}

	ips := make([]net.IP, 0, buf.Len()/net.IPv4len)
	for buf.Has(net.IPv4len) {
		ips = append(ips, net.IP(buf.CopyN(net.IPv4len)))
	}
	return ips, buf.FinError()
}

// IPsToBytes marshals an IPv4 address to a DHCP packet as specified by RFC
// 2132, Section 3.5 et al.
func IPsToBytes(i []net.IP) []byte {
	buf := uio.NewBigEndianBuffer(nil)

	for _, ip := range i {
		buf.WriteBytes(ip.To4())
	}
	return buf.Data()
}

// IPsToString returns a human-readable representation of a list of IPs.
func IPsToString(i []net.IP) string {
	s := make([]string, 0, len(i))
	for _, ip := range i {
		s = append(s, ip.String())
	}
	return strings.Join(s, ", ")
}

// OptRouter implements the router option described by RFC 2132, Section 3.5.
type OptRouter struct {
	Routers []net.IP
}

// ParseOptRouter returns a new OptRouter from a byte stream, or error if any.
func ParseOptRouter(data []byte) (*OptRouter, error) {
	ips, err := ParseIPs(data)
	if err != nil {
		return nil, err
	}
	return &OptRouter{Routers: ips}, nil
}

// Code returns the option code.
func (o *OptRouter) Code() OptionCode {
	return OptionRouter
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptRouter) ToBytes() []byte {
	return IPsToBytes(o.Routers)
}

// String returns a human-readable string.
func (o *OptRouter) String() string {
	return fmt.Sprintf("Routers -> %s", IPsToString(o.Routers))
}

// OptNTPServers implements the NTP servers option described by RFC 2132,
// Section 8.3.
type OptNTPServers struct {
	NTPServers []net.IP
}

// ParseOptNTPServers returns a new OptNTPServers from a byte stream, or error if any.
func ParseOptNTPServers(data []byte) (*OptNTPServers, error) {
	ips, err := ParseIPs(data)
	if err != nil {
		return nil, err
	}
	return &OptNTPServers{NTPServers: ips}, nil
}

// Code returns the option code.
func (o *OptNTPServers) Code() OptionCode {
	return OptionNTPServers
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptNTPServers) ToBytes() []byte {
	return IPsToBytes(o.NTPServers)
}

// String returns a human-readable string.
func (o *OptNTPServers) String() string {
	return fmt.Sprintf("NTP Servers -> %v", IPsToString(o.NTPServers))
}

// OptDomainNameServer implements the DNS server option described by RFC 2132,
// Section 3.8.
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
