package dhcpv6

import (
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptClientArchType(t *testing.T) {
	data := []byte{
		0, 6, // EFI_IA32
	}
	opt, err := ParseOptClientArchType(data)
	require.NoError(t, err)
	require.Equal(t, iana.EFI_IA32, opt.ArchTypes[0])
}

func TestParseOptClientArchTypeInvalid(t *testing.T) {
	data := []byte{42}
	_, err := ParseOptClientArchType(data)
	require.Error(t, err)
}

func TestOptClientArchTypeParseAndToBytes(t *testing.T) {
	data := []byte{
		0, 8, // EFI_XSCALE
	}
	expected := []byte{
		0, 61, // OptionClientArchType
		0, 2, // length
		0, 8, // EFI_XSCALE
	}
	opt, err := ParseOptClientArchType(data)
	require.NoError(t, err)
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptClientArchType(t *testing.T) {
	opt := OptClientArchType{
		ArchTypes: []iana.ArchType{iana.EFI_ITANIUM},
	}
	require.Equal(t, 2, opt.Length())
	require.Equal(t, OptionClientArchType, opt.Code())
	require.Contains(t, opt.String(), "archtype=EFI Itanium", "String() should contain the correct ArchType output")
}
