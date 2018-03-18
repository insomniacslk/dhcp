package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptBroadcastAddressInterfaceMethods(t *testing.T) {
	ip := net.IP{192, 168, 0, 1}
	o := OptBroadcastAddress{OptGenericIP{IP: ip}}

	require.Equal(t, OptionBroadcastAddress, o.Code(), "Code")

	expectedBytes := []byte{byte(OptionBroadcastAddress), 4, 192, 168, 0, 1}
	require.Equal(t, expectedBytes, o.ToBytes(), "ToBytes")

	require.Equal(t, 4, o.Length(), "Length")

	require.Equal(t, "Broadcast Address -> 192.168.0.1", o.String(), "String")
}

func TestParseOptBroadcastAddress(t *testing.T) {
	var (
		o   *OptBroadcastAddress
		err error
	)
	o, err = ParseOptBroadcastAddress([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptBroadcastAddress([]byte{byte(OptionBroadcastAddress), 4, 192})
	require.Error(t, err, "short byte stream")

	o, err = ParseOptBroadcastAddress([]byte{byte(OptionBroadcastAddress), 3, 192, 168, 0, 1})
	require.Error(t, err, "wrong IP length")

	o, err = ParseOptBroadcastAddress([]byte{53, 4, 192, 168, 1})
	require.Error(t, err, "wrong option code")

	o, err = ParseOptBroadcastAddress([]byte{byte(OptionBroadcastAddress), 4, 192, 168, 0, 1})
	require.NoError(t, err)
	require.Equal(t, net.IP{192, 168, 0, 1}, o.IP)
}
