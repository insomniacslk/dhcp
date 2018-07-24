package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptBootfileNameCode(t *testing.T) {
	opt := OptBootfileName{}
	require.Equal(t, OptionBootfileName, opt.Code())
}

func TestOptBootfileNameToBytes(t *testing.T) {
	opt := OptBootfileName{
		BootfileName: []byte("linuxboot"),
	}
	data := opt.ToBytes()
	expected := []byte{
		67, // OptionBootfileName
		9,  // length
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestParseOptBootfileName(t *testing.T) {
	expected := []byte{
		67, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptBootfileName(expected)
	require.NoError(t, err)
	require.Equal(t, 9, opt.Length())
	require.Equal(t, "linuxboot", string(opt.BootfileName))
}

func TestParseOptBootfileNameZeroLength(t *testing.T) {
	expected := []byte{
		67, 0,
	}
	_, err := ParseOptBootfileName(expected)
	require.Error(t, err)
}

func TestParseOptBootfileNameInvalidLength(t *testing.T) {
	expected := []byte{
		67, 9, 'l', 'i', 'n', 'u', 'x', 'b',
	}
	_, err := ParseOptBootfileName(expected)
	require.Error(t, err)
}

func TestParseOptBootfileNameShortLength(t *testing.T) {
	expected := []byte{
		67, 4, 'l', 'i', 'n', 'u', 'x',
	}
	opt, err := ParseOptBootfileName(expected)
	require.NoError(t, err)
	require.Equal(t, []byte("linu"), opt.BootfileName)
}
