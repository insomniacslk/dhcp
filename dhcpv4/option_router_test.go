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
	require.Equal(t, net.IPv4len*len(routers), o.Length(), "Length")
	require.Equal(t, routers, o.Routers, "Routers")
}

func TestParseOptRouter(t *testing.T) {
	data := []byte{
		byte(OptionRouter),
		8,               // Length
		192, 168, 0, 10, // Router #1
		192, 168, 0, 20, // Router #2
	}
	o, err := ParseOptRouter(data)
	require.NoError(t, err)
	routers := []net.IP{
		net.IPv4(192, 168, 0, 10),
		net.IPv4(192, 168, 0, 20),
	}
	require.Equal(t, &OptRouter{Routers: routers}, o)

	// Short byte stream
	data = []byte{byte(OptionDomainNameServer)}
	_, err = ParseOptDomainNameServer(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 2, 1, 1}
	_, err = ParseOptDomainNameServer(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{byte(OptionDomainNameServer), 6, 1, 1, 1}
	_, err = ParseOptDomainNameServer(data)
	require.Error(t, err, "should get error from bad length")
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
	o := OptRouter{Routers: []net.IP{net.IPv4(192, 168, 0, 1), net.IPv4(192, 168, 0, 10)}}
	require.Equal(t, "Routers -> 192.168.0.1, 192.168.0.10", o.String())
}
