package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptDomainNameServerInterfaceMethods(t *testing.T) {
	servers := []net.IP{
		net.IPv4(192, 168, 0, 10),
		net.IPv4(192, 168, 0, 20),
	}
	o := OptDomainNameServer{NameServers: servers}
	require.Equal(t, OptionDomainNameServer, o.Code(), "Code")
	require.Equal(t, net.IPv4len*len(servers), o.Length(), "Length")
	require.Equal(t, servers, o.NameServers, "NameServers")
}

func TestParseOptDomainNameServer(t *testing.T) {
	data := []byte{
		byte(OptionDomainNameServer),
		8,               // Length
		192, 168, 0, 10, // DNS #1
		192, 168, 0, 20, // DNS #2
	}
	o, err := ParseOptDomainNameServer(data)
	require.NoError(t, err)
	servers := []net.IP{
		net.IPv4(192, 168, 0, 10),
		net.IPv4(192, 168, 0, 20),
	}
	require.Equal(t, &OptDomainNameServer{NameServers: servers}, o)

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

func TestParseOptDomainNameServerNoServers(t *testing.T) {
	// RFC2132 requires that at least one DNS server IP is specified
	data := []byte{
		byte(OptionDomainNameServer),
		0, // Length
	}
	_, err := ParseOptDomainNameServer(data)
	require.Error(t, err)
}

func TestOptDomainNameServerString(t *testing.T) {
	o := OptDomainNameServer{NameServers: []net.IP{net.IPv4(192, 168, 0, 1), net.IPv4(192, 168, 0, 10)}}
	require.Equal(t, "Domain Name Servers -> 192.168.0.1, 192.168.0.10", o.String())
}
