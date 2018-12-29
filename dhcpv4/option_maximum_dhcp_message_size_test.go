package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptMaximumDHCPMessageSize(t *testing.T) {
	o := OptMaxMessageSize(1500)
	require.Equal(t, OptionMaximumDHCPMessageSize, o.Code, "Code")
	require.Equal(t, []byte{5, 220}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Maximum DHCP Message Size: 1500", o.String())
}

func TestGetMaximumDHCPMessageSize(t *testing.T) {
	options := Options{OptionMaximumDHCPMessageSize.Code(): []byte{5, 220}}
	o, err := GetMaxMessageSize(options)
	require.NoError(t, err)
	require.Equal(t, uint16(1500), o)

	// Short byte stream
	options = Options{OptionMaximumDHCPMessageSize.Code(): []byte{2}}
	_, err = GetMaxMessageSize(options)
	require.Error(t, err, "should get error from short byte stream")

	// Bad length
	options = Options{OptionMaximumDHCPMessageSize.Code(): []byte{1, 1, 1}}
	_, err = GetMaxMessageSize(options)
	require.Error(t, err, "should get error from bad length")
}
