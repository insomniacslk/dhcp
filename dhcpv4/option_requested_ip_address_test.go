package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptRequestedIPAddressInterfaceMethods(t *testing.T) {
	ip := net.IP{192, 168, 0, 1}
	o := OptRequestedIPAddress{RequestedAddr: ip}

	require.Equal(t, OptionRequestedIPAddress, o.Code(), "Code")

	expectedBytes := []byte{192, 168, 0, 1}
	require.Equal(t, expectedBytes, o.ToBytes(), "ToBytes")

	require.Equal(t, 4, o.Length(), "Length")

	require.Equal(t, "Requested IP Address -> 192.168.0.1", o.String(), "String")
}

func TestParseOptRequestedIPAddress(t *testing.T) {
	var (
		o   *OptRequestedIPAddress
		err error
	)
	o, err = ParseOptRequestedIPAddress([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptRequestedIPAddress([]byte{192})
	require.Error(t, err, "wrong IP length")

	o, err = ParseOptRequestedIPAddress([]byte{192, 168, 0, 1})
	require.NoError(t, err)
	require.Equal(t, net.IP{192, 168, 0, 1}, o.RequestedAddr)
}
