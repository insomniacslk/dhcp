package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRelayPort(t *testing.T) {
	var opt optRelayPort
	err := opt.FromBytes([]byte{0x12, 0x32})
	require.NoError(t, err)
	require.Equal(t, optRelayPort{DownstreamSourcePort: 0x1232}, opt)
}

func TestRelayPortToBytes(t *testing.T) {
	op := OptRelayPort(0x3845)
	require.Equal(t, []byte{0x38, 0x45}, op.ToBytes())
}
