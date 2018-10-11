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
	require.Equal(t, expected, opt.NameServers)
	require.Equal(t, OptionDNSRecursiveNameServer, opt.Code())
	require.Equal(t, 16, opt.Length())
	require.Contains(t, opt.String(), "nameservers=[2a03:2880:fffe:c:face:b00c:0:35]", "String() should contain the correct nameservers output")
}

func TestOptDNSRecursiveNameServerToBytes(t *testing.T) {
	ns1 := net.ParseIP("2a03:2880:fffe:c:face:b00c:0:35")
	ns2 := net.ParseIP("2001:4860:4860::8888")
	nameservers := []net.IP{ns1, ns2}
	expected := []byte{
		0, 23, // OptionDNSRecursiveNameServer
		0, 32, // length
	}
	expected = append(expected, []byte(ns1)...)
	expected = append(expected, []byte(ns2)...)
	opt := OptDNSRecursiveNameServer{NameServers: nameservers}
	require.Equal(t, expected, opt.ToBytes())
}

func TestParseOptDNSRecursiveNameServerParseBogusNameserver(t *testing.T) {
	data := []byte{
		0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, // invalid IPv6 address
	}
	_, err := ParseOptDNSRecursiveNameServer(data)
	require.Error(t, err, "An invalid nameserver IPv6 address should return an error")
}
