package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptVendorClass(t *testing.T) {
	data := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // EnterpriseNumber
		0, 10, 'H', 'T', 'T', 'P', 'C', 'l', 'i', 'e', 'n', 't',
		0, 4, 't', 'e', 's', 't',
	}
	var opt OptVendorClass
	err := opt.FromBytes(data)
	require.NoError(t, err)
	require.Equal(t, OptionVendorClass, opt.Code())
	require.Equal(t, 2, len(opt.Data))
	require.Equal(t, uint32(0xaabbccdd), opt.EnterpriseNumber)
	require.Equal(t, []byte("HTTPClient"), opt.Data[0])
	require.Equal(t, []byte("test"), opt.Data[1])
}

func TestOptVendorClassToBytes(t *testing.T) {
	opt := OptVendorClass{
		EnterpriseNumber: uint32(0xaabbccdd),
		Data: [][]byte{
			[]byte("HTTPClient"),
			[]byte("test"),
		},
	}
	data := opt.ToBytes()
	expected := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // EnterpriseNumber
		0, 10, 'H', 'T', 'T', 'P', 'C', 'l', 'i', 'e', 'n', 't',
		0, 4, 't', 'e', 's', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptVendorClassParseOptVendorClassMalformed(t *testing.T) {
	buf := []byte{
		0xaa, 0xbb, // truncated EnterpriseNumber
	}
	var opt OptVendorClass
	err := opt.FromBytes(buf)
	require.Error(t, err, "ParseOptVendorClass() should error if given truncated EnterpriseNumber")

	buf = []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // EnterpriseNumber
	}
	err = opt.FromBytes(buf)
	require.Error(t, err, "ParseOptVendorClass() should error if given no vendor classes")

	buf = []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // EnterpriseNumber
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e',
	}
	err = opt.FromBytes(buf)
	require.Error(t, err, "ParseOptVendorClass() should error if given truncated vendor classes")

	buf = []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // EnterpriseNumber
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0,
	}
	err = opt.FromBytes(buf)
	require.Error(t, err, "ParseOptVendorClass() should error if given a truncated length")
}

func TestOptVendorClassString(t *testing.T) {
	data := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // EnterpriseNumber
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e', 's', 't',
	}
	var opt OptVendorClass
	err := opt.FromBytes(data)
	require.NoError(t, err)

	str := opt.String()
	require.Contains(
		t,
		str,
		"EnterpriseNumber=2864434397",
		"String() should contain the enterprisenum",
	)
	require.Contains(
		t,
		str,
		"Data=[linuxboot, test]",
		"String() should contain the list of vendor classes",
	)
}
