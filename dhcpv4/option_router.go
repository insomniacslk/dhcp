package dhcpv4

import (
	"fmt"
	"net"
	"strings"

	"github.com/u-root/u-root/pkg/uio"
)

// This option implements the router option
// https://tools.ietf.org/html/rfc2132

// OptRouter represents an option encapsulating the routers.
type OptRouter struct {
	Routers []net.IP
}

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

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptRouter) Length() int {
	return len(o.Routers) * 4
}
