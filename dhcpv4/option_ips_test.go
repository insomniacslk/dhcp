package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptRoutersInterfaceMethods(t *testing.T) {
	routers := []net.IP{
		net.IPv4(192, 168, 0, 10),
		net.IPv4(192, 168, 0, 20),
	}
	o := OptRouter{Routers: routers}
	require.Equal(t, OptionRouter, o.Code(), "Code")
	require.Equal(t, routers, o.Routers, "Routers")
}

func TestParseOptRouter(t *testing.T) {
	data := []byte{
		byte(OptionRouter),
		8,               // Length
		192, 168, 0, 10, // Router #1
		192, 168, 0, 20, // Router #2
	}
	o, err := ParseOptRouter(data[2:])
	require.NoError(t, err)
	routers := []net.IP{
		net.IP{192, 168, 0, 10},
		net.IP{192, 168, 0, 20},
	}
	require.Equal(t, &OptRouter{Routers: routers}, o)

	// Short byte stream
	data = []byte{byte(OptionRouter)}
	_, err = ParseOptRouter(data)
	require.Error(t, err, "should get error from short byte stream")
}

func TestParseOptRouterNoRouters(t *testing.T) {
	// RFC2132 requires that at least one Router IP is specified
	data := []byte{
		byte(OptionRouter),
		0, // Length
	}
	_, err := ParseOptRouter(data)
	require.Error(t, err)
}

func TestOptRouterString(t *testing.T) {
	o := OptRouter{Routers: []net.IP{net.IP{192, 168, 0, 1}, net.IP{192, 168, 0, 10}}}
	require.Equal(t, "Routers -> 192.168.0.1, 192.168.0.10", o.String())
}

func TestOptDomainNameServerInterfaceMethods(t *testing.T) {
	servers := []net.IP{
		net.IPv4(192, 168, 0, 10),
		net.IPv4(192, 168, 0, 20),
	}
	o := OptDomainNameServer{NameServers: servers}
	require.Equal(t, OptionDomainNameServer, o.Code(), "Code")
	require.Equal(t, servers, o.NameServers, "NameServers")
}

func TestParseOptDomainNameServer(t *testing.T) {
	data := []byte{
		byte(OptionDomainNameServer),
		8,               // Length
		192, 168, 0, 10, // DNS #1
		192, 168, 0, 20, // DNS #2
	}
	o, err := ParseOptDomainNameServer(data[2:])
	require.NoError(t, err)
	servers := []net.IP{
		net.IP{192, 168, 0, 10},
		net.IP{192, 168, 0, 20},
	}
	require.Equal(t, &OptDomainNameServer{NameServers: servers}, o)

	// Bad length
	data = []byte{1, 1, 1}
	_, err = ParseOptDomainNameServer(data)
	require.Error(t, err, "should get error from bad length")
}

func TestParseOptDomainNameServerNoServers(t *testing.T) {
	// RFC2132 requires that at least one DNS server IP is specified
	_, err := ParseOptDomainNameServer([]byte{})
	require.Error(t, err)
}

func TestOptDomainNameServerString(t *testing.T) {
	o := OptDomainNameServer{NameServers: []net.IP{net.IPv4(192, 168, 0, 1), net.IPv4(192, 168, 0, 10)}}
	require.Equal(t, "Domain Name Servers -> 192.168.0.1, 192.168.0.10", o.String())
}

func TestOptNTPServersInterfaceMethods(t *testing.T) {
	ntpServers := []net.IP{
		net.IPv4(192, 168, 0, 10),
		net.IPv4(192, 168, 0, 20),
	}
	o := OptNTPServers{NTPServers: ntpServers}
	require.Equal(t, OptionNTPServers, o.Code(), "Code")
	require.Equal(t, ntpServers, o.NTPServers, "NTPServers")
}

func TestParseOptNTPServers(t *testing.T) {
	data := []byte{
		byte(OptionNTPServers),
		8,               // Length
		192, 168, 0, 10, // NTP server #1
		192, 168, 0, 20, // NTP server #2
	}
	o, err := ParseOptNTPServers(data[2:])
	require.NoError(t, err)
	ntpServers := []net.IP{
		net.IP{192, 168, 0, 10},
		net.IP{192, 168, 0, 20},
	}
	require.Equal(t, &OptNTPServers{NTPServers: ntpServers}, o)

	// Bad length
	data = []byte{1, 1, 1}
	_, err = ParseOptNTPServers(data)
	require.Error(t, err, "should get error from bad length")
}

func TestParseOptNTPserversNoNTPServers(t *testing.T) {
	// RFC2132 requires that at least one NTP server IP is specified
	_, err := ParseOptNTPServers([]byte{})
	require.Error(t, err)
}

func TestOptNTPServersString(t *testing.T) {
	o := OptNTPServers{NTPServers: []net.IP{net.IPv4(192, 168, 0, 1), net.IPv4(192, 168, 0, 10)}}
	require.Equal(t, "NTP Servers -> 192.168.0.1, 192.168.0.10", o.String())
}
