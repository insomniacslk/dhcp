package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptSubnetMaskInterfaceMethods(t *testing.T) {
	mask := net.IPMask{255, 255, 255, 0}
	o := OptSubnetMask{SubnetMask: mask}

	require.Equal(t, OptionSubnetMask, o.Code(), "Code")

	expectedBytes := []byte{1, 4, 255, 255, 255, 0}
	require.Equal(t, expectedBytes, o.ToBytes(), "ToBytes")

	require.Equal(t, 4, o.Length(), "Length")

	require.Equal(t, "Subnet Mask -> ffffff00", o.String(), "String")
}

func TestParseOptSubnetMask(t *testing.T) {
	var (
		o   *OptSubnetMask
		err error
	)
	o, err = ParseOptSubnetMask([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptSubnetMask([]byte{1, 4, 255})
	require.Error(t, err, "short byte stream")

	o, err = ParseOptSubnetMask([]byte{1, 3, 255, 255, 255, 0})
	require.Error(t, err, "wrong IP length")

	o, err = ParseOptSubnetMask([]byte{2, 4, 255, 255, 255})
	require.Error(t, err, "wrong option code")

	o, err = ParseOptSubnetMask([]byte{1, 4, 255, 255, 255, 0})
	require.NoError(t, err)
	require.Equal(t, net.IPMask{255, 255, 255, 0}, o.SubnetMask)
}
