package dhcpv6

// This module defines the OptDNSRecursiveNameServer structure.
// https://www.ietf.org/rfc/rfc3646.txt

import (
	"encoding/binary"
	"fmt"
	"net"
)

// OptDNSRecursiveNameServer represents a DNS_RECURSIVE_NAME_SERVER option
type OptDNSRecursiveNameServer struct {
	NameServers []net.IP
}

// Code returns the option code
func (op *OptDNSRecursiveNameServer) Code() OptionCode {
	return DNS_RECURSIVE_NAME_SERVER
}

// ToBytes returns the option serialized to bytes, including option code and
// length
func (op *OptDNSRecursiveNameServer) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(DNS_RECURSIVE_NAME_SERVER))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	for _, ns := range op.NameServers {
		buf = append(buf, ns...)
	}
	return buf
}

// Length returns the option length
func (op *OptDNSRecursiveNameServer) Length() int {
	return len(op.NameServers) * net.IPv6len
}

func (op *OptDNSRecursiveNameServer) String() string {
	return fmt.Sprintf("OptDNSRecursiveNameServer{nameservers=%v}", op.NameServers)
}

// ParseOptDNSRecursiveNameServer builds an OptDNSRecursiveNameServer structure
// from a sequence of bytes. The input data does not include option code and length
// bytes.
func ParseOptDNSRecursiveNameServer(data []byte) (*OptDNSRecursiveNameServer, error) {
	if len(data)%net.IPv6len != 0 {
		return nil, fmt.Errorf("Invalid OptDNSRecursiveNameServer data: length is not a multiple of %d", net.IPv6len)
	}
	opt := OptDNSRecursiveNameServer{}
	var nameServers []net.IP
	for i := 0; i < len(data); i += net.IPv6len {
		nameServers = append(nameServers, data[i:i+net.IPv6len])
	}
	opt.NameServers = nameServers
	return &opt, nil
}
