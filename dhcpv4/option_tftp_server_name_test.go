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
		TFTPServerName: []byte("linuxboot"),
	}
	data := opt.ToBytes()
	expected := []byte{
		66, // OptionTFTPServerName
		9,  // length
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestParseOptTFTPServerName(t *testing.T) {
	expected := []byte{
		66, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptTFTPServerName(expected)
	require.NoError(t, err)
	require.Equal(t, 9, opt.Length())
	require.Equal(t, "linuxboot", string(opt.TFTPServerName))
}

func TestParseOptTFTPServerNameZeroLength(t *testing.T) {
	expected := []byte{
		66, 0,
	}
	_, err := ParseOptTFTPServerName(expected)
	require.Error(t, err)
}

func TestParseOptTFTPServerNameInvalidLength(t *testing.T) {
	expected := []byte{
		66, 9, 'l', 'i', 'n', 'u', 'x', 'b',
	}
	_, err := ParseOptTFTPServerName(expected)
	require.Error(t, err)
}

func TestParseOptTFTPServerNameShortLength(t *testing.T) {
	expected := []byte{
		66, 4, 'l', 'i', 'n', 'u', 'x',
	}
	opt, err := ParseOptTFTPServerName(expected)
	require.NoError(t, err)
	require.Equal(t, []byte("linu"), opt.TFTPServerName)
}

func TestOptTFTPServerNameString(t *testing.T) {
	o := OptTFTPServerName{TFTPServerName: []byte("testy test")}
	require.Equal(t, "TFTP Server Name -> testy test", o.String())
}
