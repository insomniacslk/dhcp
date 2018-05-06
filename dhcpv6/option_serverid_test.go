package dhcpv6

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptServerId(t *testing.T) {
	data := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		0, 1, 2, 3, 4, 5, // hw addr
	}
	opt, err := ParseOptServerId(data)
	require.NoError(t, err)
	require.Equal(t, opt.Sid.Type, DUID_LL)
	require.Equal(t, opt.Sid.HwType, iana.HwTypeEthernet)
	require.Equal(t, opt.Sid.LinkLayerAddr, net.HardwareAddr([]byte{0, 1, 2, 3, 4, 5}))
}

func TestOptServerIdToBytes(t *testing.T) {
	opt := OptServerId{
		Sid: Duid{
			Type:          DUID_LL,
			HwType:        iana.HwTypeEthernet,
			LinkLayerAddr: net.HardwareAddr([]byte{5, 4, 3, 2, 1, 0}),
		},
	}
	expected := []byte{
		0, 2, // OPTION_SERVERID
		0, 10, // length
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		5, 4, 3, 2, 1, 0, // hw addr
	}
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptServerIdDecodeEncode(t *testing.T) {
	data := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		5, 4, 3, 2, 1, 0, // hw addr
	}
	expected := append([]byte{
		0, 2, // OPTION_SERVERID
		0, 10, // length
	}, data...)
	opt, err := ParseOptServerId(data)
	require.NoError(t, err)
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptionServerId(t *testing.T) {
	opt := OptServerId{
		Sid: Duid{
			Type:          DUID_LL,
			HwType:        iana.HwTypeEthernet,
			LinkLayerAddr: net.HardwareAddr([]byte{0xde, 0xad, 0, 0, 0xbe, 0xef}),
		},
	}
	require.Equal(t, opt.Length(), 10)
	require.Equal(t, opt.Code(), OPTION_SERVERID)
}
