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
	if len(data) < 2 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionNTPServers {
		return nil, fmt.Errorf("expected code %v, got %v", OptionNTPServers, code)
	}
	length := int(data[1])
	if length == 0 || length%4 != 0 {
		return nil, fmt.Errorf("Invalid length: expected multiple of 4 larger than 4, got %v", length)
	}
	if len(data) < 2+length {
		return nil, ErrShortByteStream
	}
	ntpServers := make([]net.IP, 0, length%4)
	for idx := 0; idx < length; idx += 4 {
		b := data[2+idx : 2+idx+4]
		ntpServers = append(ntpServers, net.IPv4(b[0], b[1], b[2], b[3]))
	}
	return &OptNTPServers{NTPServers: ntpServers}, nil
}

// Code returns the option code.
func (o *OptNTPServers) Code() OptionCode {
	return OptionNTPServers
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptNTPServers) ToBytes() []byte {
	ret := []byte{byte(o.Code()), byte(o.Length())}
	for _, ntp := range o.NTPServers {
		ret = append(ret, ntp.To4()...)
	}
	return ret
}

// String returns a human-readable string.
func (o *OptNTPServers) String() string {
	var ntpServers string
	for idx, ntp := range o.NTPServers {
		ntpServers += ntp.String()
		if idx < len(o.NTPServers)-1 {
			ntpServers += ", "
		}
	}
	return fmt.Sprintf("NTP Servers -> %v", ntpServers)
}

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptNTPServers) Length() int {
	return len(o.NTPServers) * 4
}
