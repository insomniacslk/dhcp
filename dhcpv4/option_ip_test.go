package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptBroadcastAddress(t *testing.T) {
	o := OptBroadcastAddress{BroadcastAddress: net.IP{192, 168, 0, 1}}

	require.Equal(t, OptionBroadcastAddress, o.Code(), "Code")
	require.Equal(t, []byte{192, 168, 0, 1}, o.ToBytes(), "ToBytes")
	require.Equal(t, "Broadcast Address -> 192.168.0.1", o.String(), "String")
}

func TestParseOptBroadcastAddress(t *testing.T) {
	o, err := ParseOptBroadcastAddress([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptBroadcastAddress([]byte{192, 168, 0})
	require.Error(t, err, "wrong IP length")

	o, err = ParseOptBroadcastAddress([]byte{192, 168, 0, 1})
	require.NoError(t, err)
	require.Equal(t, net.IP{192, 168, 0, 1}, o.BroadcastAddress)
}

func TestOptRequestedIPAddress(t *testing.T) {
	o := OptRequestedIPAddress{RequestedAddr: net.IP{192, 168, 0, 1}}

	require.Equal(t, OptionRequestedIPAddress, o.Code(), "Code")
	require.Equal(t, []byte{192, 168, 0, 1}, o.ToBytes(), "ToBytes")
	require.Equal(t, "Requested IP Address -> 192.168.0.1", o.String(), "String")
}

func TestParseOptRequestedIPAddress(t *testing.T) {
	o, err := ParseOptRequestedIPAddress([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptRequestedIPAddress([]byte{192})
	require.Error(t, err, "wrong IP length")

	o, err = ParseOptRequestedIPAddress([]byte{192, 168, 0, 1})
	require.NoError(t, err)
	require.Equal(t, net.IP{192, 168, 0, 1}, o.RequestedAddr)
}

func TestOptServerIdentifierInterfaceMethods(t *testing.T) {
	o := OptServerIdentifier{ServerID: net.IP{192, 168, 0, 1}}

	require.Equal(t, OptionServerIdentifier, o.Code(), "Code")
	require.Equal(t, []byte{192, 168, 0, 1}, o.ToBytes(), "ToBytes")
	require.Equal(t, "Server Identifier -> 192.168.0.1", o.String(), "String")
}

func TestParseOptServerIdentifier(t *testing.T) {
	o, err := ParseOptServerIdentifier([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptServerIdentifier([]byte{192, 168, 0})
	require.Error(t, err, "wrong IP length")

	o, err = ParseOptServerIdentifier([]byte{192, 168, 0, 1})
	require.NoError(t, err)
	require.Equal(t, net.IP{192, 168, 0, 1}, o.ServerID)
}
