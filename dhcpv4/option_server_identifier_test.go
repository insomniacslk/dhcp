package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptServerIdentifierInterfaceMethods(t *testing.T) {
	ip := net.IP{192, 168, 0, 1}
	o := OptServerIdentifier{ServerID: ip}

	require.Equal(t, OptionServerIdentifier, o.Code(), "Code")

	expectedBytes := []byte{192, 168, 0, 1}
	require.Equal(t, expectedBytes, o.ToBytes(), "ToBytes")

	require.Equal(t, 4, o.Length(), "Length")

	require.Equal(t, "Server Identifier -> 192.168.0.1", o.String(), "String")
}

func TestParseOptServerIdentifier(t *testing.T) {
	var (
		o   *OptServerIdentifier
		err error
	)
	o, err = ParseOptServerIdentifier([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptServerIdentifier([]byte{192, 168, 0})
	require.Error(t, err, "wrong IP length")

	o, err = ParseOptServerIdentifier([]byte{192, 168, 0, 1})
	require.NoError(t, err)
	require.Equal(t, net.IP{192, 168, 0, 1}, o.ServerID)
}
