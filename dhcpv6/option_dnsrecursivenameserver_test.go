package dhcpv6

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptDNSRecursiveNameServer(t *testing.T) {
	data := []byte{
		0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35,
	}
	expected := []net.IP{
		net.IP(data),
	}
	opt, err := ParseOptDNSRecursiveNameServer(data)
	require.NoError(t, err)
	require.Equal(t, opt.NameServers, expected)
	require.Equal(t, opt.Length(), 16)
}

func TestOptDNSRecursiveNameServerToBytes(t *testing.T) {
	ns1 := net.ParseIP("2a03:2880:fffe:c:face:b00c:0:35")
	ns2 := net.ParseIP("2001:4860:4860::8888")
	nameservers := []net.IP{ns1, ns2}
	expected := []byte{
		0, 23, // DNS_RECURSIVE_NAME_SERVER
		0, 32, // length
	}
	expected = append(expected, []byte(ns1)...)
	expected = append(expected, []byte(ns2)...)
	opt := OptDNSRecursiveNameServer{NameServers: nameservers}
	require.Equal(t, opt.ToBytes(), expected)
}
