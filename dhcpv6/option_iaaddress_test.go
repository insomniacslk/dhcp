package dhcpv6

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptIAAddressParse(t *testing.T) {
	ipaddr := []byte{0x24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	data := append(ipaddr, []byte{
		0xa, 0xb, 0xc, 0xd, // preferred lifetime
		0xe, 0xf, 0x1, 0x2, // valid lifetime
		0, 8, 0, 2, 0xaa, 0xbb, // options
	}...)
	opt, err := ParseOptIAAddress(data)
	require.NoError(t, err)
	require.Equal(t, 30, opt.Length())
	require.Equal(t, net.IP(ipaddr), opt.IPv6Addr)
	require.Equal(t, uint32(0x0a0b0c0d), opt.PreferredLifetime)
	require.Equal(t, uint32(0x0e0f0102), opt.ValidLifetime)
}

func TestOptIAAddressParseInvalidTooShort(t *testing.T) {
	data := []byte{
		0x24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0xa, 0xb, 0xc, 0xd, // preferred lifetime
		// truncated here
	}
	_, err := ParseOptIAAddress(data)
	require.Error(t, err)
}

func TestOptIAAddressParseInvalidBrokenOptions(t *testing.T) {
	data := []byte{
		0x24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0xa, 0xb, 0xc, 0xd, // preferred lifetime
		0xe, 0xf, 0x1, 0x2, // valid lifetime
		0, 8, 0, 2, 0xaa, // broken options
	}
	_, err := ParseOptIAAddress(data)
	require.Error(t, err)
}

func TestOptIAAddressToBytes(t *testing.T) {
	expected := []byte{
		0, 5, // OptionIAAddr
		0, 30, // length
	}
	ipBytes := []byte{0x24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	expected = append(expected, ipBytes...)
	expected = append(expected, []byte{
		0xa, 0xb, 0xc, 0xd, // preferred lifetime
		0xe, 0xf, 0x1, 0x2, // valid lifetime
		0, 8, 0, 2, 0xaa, 0xbb, // options
	}...)
	opt := OptIAAddress{
		IPv6Addr:          net.IP(ipBytes),
		PreferredLifetime: 0x0a0b0c0d,
		ValidLifetime:     0x0e0f0102,
		Options: []Option{
			&OptElapsedTime{
				ElapsedTime: 0xaabb,
			},
		},
	}
	require.Equal(t, expected, opt.ToBytes())
}
