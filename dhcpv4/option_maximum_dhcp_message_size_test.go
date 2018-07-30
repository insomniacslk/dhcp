package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptMaximumDHCPMessageSizeInterfaceMethods(t *testing.T) {
	o := OptMaximumDHCPMessageSize{Size: 1500}
	require.Equal(t, OptionMaximumDHCPMessageSize, o.Code(), "Code")
	require.Equal(t, 2, o.Length(), "Length")
	require.Equal(t, []byte{57, 2, 5, 220}, o.ToBytes(), "ToBytes")
}

func TestParseOptMaximumDHCPMessageSize(t *testing.T) {
	data := []byte{57, 2, 5, 220}
	o, err := ParseOptMaximumDHCPMessageSize(data)
	require.NoError(t, err)
	require.Equal(t, &OptMaximumDHCPMessageSize{Size: 1500}, o)

	// Short byte stream
	data = []byte{57, 2}
	_, err = ParseOptMaximumDHCPMessageSize(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 2, 1, 1}
	_, err = ParseOptMaximumDHCPMessageSize(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{57, 3, 1, 1, 1}
	_, err = ParseOptMaximumDHCPMessageSize(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptMaximumDHCPMessageSizeString(t *testing.T) {
	o := OptMaximumDHCPMessageSize{Size: 1500}
	require.Equal(t, "Maximum DHCP Message Size -> 1500", o.String())
}
