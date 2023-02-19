package dhcpv6

import (
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMessageOptionsWithDHCP4oDHCP6Server(t *testing.T) {
	ip := net.IP{0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35}
	data := append([]byte{
		0, 88, // DHCP4oDHCP6 option.
		0, 16, // length
	}, ip...)

	want := []net.IP{ip}
	var mo MessageOptions
	if err := mo.FromBytes(data); err != nil {
		t.Errorf("FromBytes = %v", err)
	} else if got := mo.DHCP4oDHCP6Server(); !reflect.DeepEqual(got.DHCP4oDHCP6Servers, want) {
		t.Errorf("FromBytes = %v, want %v", got.DHCP4oDHCP6Servers, want)
	}
}

func TestParseOptDHCP4oDHCP6Server(t *testing.T) {
	data := []byte{
		0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35,
	}
	expected := []net.IP{
		net.IP(data),
	}
	var opt OptDHCP4oDHCP6Server
	err := opt.FromBytes(data)
	require.NoError(t, err)
	require.Equal(t, expected, opt.DHCP4oDHCP6Servers)
	require.Equal(t, OptionDHCP4oDHCP6Server, opt.Code())
	require.Contains(t, opt.String(), "[2a03:2880:fffe:c:face:b00c:0:35]", "String() should contain the correct DHCP4-over-DHCP6 server output")
}

func TestOptDHCP4oDHCP6ServerToBytes(t *testing.T) {
	ip1 := net.ParseIP("2a03:2880:fffe:c:face:b00c:0:35")
	ip2 := net.ParseIP("2001:4860:4860::8888")
	opt := OptDHCP4oDHCP6Server{DHCP4oDHCP6Servers: []net.IP{ip1, ip2}}

	want := []byte(append(ip1, ip2...))
	require.Equal(t, want, opt.ToBytes())
}

func TestParseOptDHCP4oDHCP6ServerParseNoAddr(t *testing.T) {
	data := []byte{}
	var expected []net.IP
	var opt OptDHCP4oDHCP6Server
	err := opt.FromBytes(data)
	require.NoError(t, err)
	require.Equal(t, expected, opt.DHCP4oDHCP6Servers)
}

func TestOptDHCP4oDHCP6ServerToBytesNoAddr(t *testing.T) {
	expected := []byte(nil)
	opt := OptDHCP4oDHCP6Server{}
	require.Equal(t, expected, opt.ToBytes())
}

func TestParseOptDHCP4oDHCP6ServerParseBogus(t *testing.T) {
	data := []byte{
		0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, // invalid IPv6 address
	}
	var opt OptDHCP4oDHCP6Server
	err := opt.FromBytes(data)
	require.Error(t, err, "An invalid IPv6 address should return an error")
}
