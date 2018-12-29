package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptTFTPServerNameCode(t *testing.T) {
	opt := OptTFTPServerName{}
	require.Equal(t, OptionTFTPServerName, opt.Code())
}

func TestOptTFTPServerNameToBytes(t *testing.T) {
	opt := OptTFTPServerName{
		TFTPServerName: "linuxboot",
	}
	data := opt.ToBytes()
	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestParseOptTFTPServerName(t *testing.T) {
	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptTFTPServerName(expected)
	require.NoError(t, err)
	require.Equal(t, "linuxboot", string(opt.TFTPServerName))
}

func TestOptTFTPServerNameString(t *testing.T) {
	o := OptTFTPServerName{TFTPServerName: "testy test"}
	require.Equal(t, "TFTP Server Name -> testy test", o.String())
}
