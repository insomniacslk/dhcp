package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptNTPServersInterfaceMethods(t *testing.T) {
	ntpServers := []net.IP{
		net.IPv4(192, 168, 0, 10),
		net.IPv4(192, 168, 0, 20),
	}
	o := OptNTPServers{NTPServers: ntpServers}
	require.Equal(t, OptionNTPServers, o.Code(), "Code")
	require.Equal(t, net.IPv4len*len(ntpServers), o.Length(), "Length")
	require.Equal(t, ntpServers, o.NTPServers, "NTPServers")
}

func TestParseOptNTPServers(t *testing.T) {
	data := []byte{
		byte(OptionNTPServers),
		8,               // Length
		192, 168, 0, 10, // NTP server #1
		192, 168, 0, 20, // NTP server #2
	}
	o, err := ParseOptNTPServers(data)
	require.NoError(t, err)
	ntpServers := []net.IP{
		net.IPv4(192, 168, 0, 10),
		net.IPv4(192, 168, 0, 20),
	}
	require.Equal(t, &OptNTPServers{NTPServers: ntpServers}, o)

	// Short byte stream
	data = []byte{byte(OptionNTPServers)}
	_, err = ParseOptNTPServers(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 2, 1, 1}
	_, err = ParseOptNTPServers(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{byte(OptionNTPServers), 6, 1, 1, 1}
	_, err = ParseOptNTPServers(data)
	require.Error(t, err, "should get error from bad length")
}

func TestParseOptNTPserversNoNTPServers(t *testing.T) {
	// RFC2132 requires that at least one NTP server IP is specified
	data := []byte{
		byte(OptionNTPServers),
		0, // Length
	}
	_, err := ParseOptNTPServers(data)
	require.Error(t, err)
}

func TestOptNTPServersString(t *testing.T) {
	o := OptNTPServers{NTPServers: []net.IP{net.IPv4(192, 168, 0, 1), net.IPv4(192, 168, 0, 10)}}
	require.Equal(t, "NTP Servers -> 192.168.0.1, 192.168.0.10", o.String())
}
