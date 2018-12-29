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
		BootfileName: "linuxboot",
	}
	data := opt.ToBytes()
	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestParseOptBootfileName(t *testing.T) {
	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptBootfileName(expected)
	require.NoError(t, err)
	require.Equal(t, 9, opt.Length())
	require.Equal(t, "linuxboot", opt.BootfileName)
}

func TestOptBootfileNameString(t *testing.T) {
	o := OptBootfileName{BootfileName: "testy test"}
	require.Equal(t, "Bootfile Name -> testy test", o.String())
}
