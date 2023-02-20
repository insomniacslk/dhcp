package dhcpv6

import (
	"net"
	"reflect"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseMessageOptionsWithClientID(t *testing.T) {
	buf := []byte{
		0, 1, // Client ID option
		0, 10, // length
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		0, 1, 2, 3, 4, 5, // HW addr
	}

	want := &DUIDLL{HWType: iana.HWTypeEthernet, LinkLayerAddr: net.HardwareAddr{0, 1, 2, 3, 4, 5}}
	var mo MessageOptions
	if err := mo.FromBytes(buf); err != nil {
		t.Errorf("FromBytes = %v", err)
	} else if got := mo.ClientID(); !reflect.DeepEqual(got, want) {
		t.Errorf("ClientID = %v, want %v", got, want)
	}
}

func TestParseOptClientID(t *testing.T) {
	data := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		0, 1, 2, 3, 4, 5, // hw addr
	}
	var opt optClientID
	err := opt.FromBytes(data)
	require.NoError(t, err)
	want := OptClientID(
		&DUIDLL{
			HWType:        iana.HWTypeEthernet,
			LinkLayerAddr: net.HardwareAddr([]byte{0, 1, 2, 3, 4, 5}),
		},
	)
	require.Equal(t, want, &opt)
}

func TestOptClientIdToBytes(t *testing.T) {
	opt := OptClientID(
		&DUIDLL{
			HWType:        iana.HWTypeEthernet,
			LinkLayerAddr: net.HardwareAddr([]byte{5, 4, 3, 2, 1, 0}),
		},
	)
	expected := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		5, 4, 3, 2, 1, 0, // hw addr
	}
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptClientIdDecodeEncode(t *testing.T) {
	data := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		5, 4, 3, 2, 1, 0, // hw addr
	}
	var opt optClientID
	err := opt.FromBytes(data)
	require.NoError(t, err)
	require.Equal(t, data, opt.ToBytes())
}

func TestOptionClientId(t *testing.T) {
	opt := OptClientID(
		&DUIDLL{
			HWType:        iana.HWTypeEthernet,
			LinkLayerAddr: net.HardwareAddr([]byte{0xde, 0xad, 0, 0, 0xbe, 0xef}),
		},
	)
	require.Equal(t, OptionClientID, opt.Code())
	require.Contains(
		t,
		opt.String(),
		"Client ID: DUID-LL{HWType=Ethernet HWAddr=de:ad:00:00:be:ef}",
		"String() should contain the correct cid output",
	)
}

func TestOptClientIdparseOptClientIDBogusDUID(t *testing.T) {
	data := []byte{
		0, 4, // DUID_UUID
		1, 2, 3, 4, 5, 6, 7, 8, 9, // a UUID should be 18 bytes not 17
		10, 11, 12, 13, 14, 15, 16, 17,
	}
	var opt optClientID
	err := opt.FromBytes(data)
	require.Error(t, err, "A truncated OptClientId DUID should return an error")
}

func TestOptClientIdparseOptClientIDInvalidTooShort(t *testing.T) {
	data := []byte{
		0, // truncated: DUIDs are at least 2 bytes
	}
	var opt optClientID
	err := opt.FromBytes(data)
	require.Error(t, err, "A truncated OptClientId should return an error")
}
