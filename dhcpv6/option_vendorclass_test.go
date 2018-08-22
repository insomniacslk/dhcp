package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptVendorClass(t *testing.T) {
	data := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // EntID
		0, 10, 'H', 'T', 'T', 'P', 'C', 'l', 'i', 'e', 'n', 't',
		0, 4, 't', 'e', 's', 't',
	}
	opt, err := ParseOptVendorClass(data)
	require.NoError(t, err)
	require.Equal(t, len(opt.Data), 2)
	require.Equal(t, opt.EntID, uint32(0xaabbccdd))
	require.Equal(t, []byte("HTTPClient"), opt.Data[0])
	require.Equal(t, []byte("test"), opt.Data[1])
}

func TestOptVendorClassToBytesMultiple(t *testing.T) {
	opt := OptVendorClass{
		EntID: uint32(0xaabbccdd),
		Data: [][]byte{
			[]byte("HTTPClient"),
			[]byte("test"),
		},
	}
	data := opt.ToBytes()
	expected := []byte{
		0, 16, // OptionVendorClass
		0, 22, // length
		0xaa, 0xbb, 0xcc, 0xdd, //EntID
		0, 10, 'H', 'T', 'T', 'P', 'C', 'l', 'i', 'e', 'n', 't',
		0, 4, 't', 'e', 's', 't',
	}
	require.Equal(t, expected, data)
}
