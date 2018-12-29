package dhcpv4

import (
	"fmt"
	"net"
)

// This option implements the network time protocol servers option
// https://tools.ietf.org/html/rfc2132

// OptNTPServers represents an option encapsulating the NTP servers.
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
