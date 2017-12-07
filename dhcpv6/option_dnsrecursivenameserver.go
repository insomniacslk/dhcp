package dhcpv6

// This module defines the OptDNSRecursiveNameServer structure.
// https://www.ietf.org/rfc/rfc3646.txt

import (
	"encoding/binary"
	"fmt"
	"net"
)

type OptDNSRecursiveNameServer struct {
	nameServers []net.IP
}

func (op *OptDNSRecursiveNameServer) Code() OptionCode {
	return DNS_RECURSIVE_NAME_SERVER
}

func (op *OptDNSRecursiveNameServer) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(DNS_RECURSIVE_NAME_SERVER))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	for _, ns := range op.nameServers {
		buf = append(buf, ns...)
	}
	return buf
}

func (op *OptDNSRecursiveNameServer) NameServers() []net.IP {
	return op.nameServers
}

func (op *OptDNSRecursiveNameServer) SetNameServers(ns []net.IP) {
	op.nameServers = ns
}

func (op *OptDNSRecursiveNameServer) Length() int {
	return len(op.nameServers) * net.IPv6len
}

func (op *OptDNSRecursiveNameServer) String() string {
	return fmt.Sprintf("OptDNSRecursiveNameServer{nameservers=%v}", op.nameServers)
}

// build an OptDNSRecursiveNameServer structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptDNSRecursiveNameServer(data []byte) (*OptDNSRecursiveNameServer, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("Invalid OptDNSRecursiveNameServer data: length is not a multiple of 2")
	}
	opt := OptDNSRecursiveNameServer{}
	var nameServers []net.IP
	for i := 0; i < len(data); i += net.IPv6len {
		nameServers = append(nameServers, data[i:i+net.IPv6len])
	}
	opt.nameServers = nameServers
	return &opt, nil
}
