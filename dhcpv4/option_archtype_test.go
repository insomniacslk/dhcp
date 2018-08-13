package dhcpv4

import (
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptClientArchType(t *testing.T) {
	data := []byte{
		93,   // OptionClientSystemArchitectureType
		2,    // Length
		0, 6, // EFI_IA32
	}
	opt, err := ParseOptClientArchType(data)
	require.NoError(t, err)
	require.Equal(t, opt.ArchTypes[0], iana.EFI_IA32)
}

func TestParseOptClientArchTypeMultiple(t *testing.T) {
	data := []byte{
		93,   // OptionClientSystemArchitectureType
		4,    // Length
		0, 6, // EFI_IA32
		0, 2, // EFI_ITANIUM
	}
	opt, err := ParseOptClientArchType(data)
	require.NoError(t, err)
	require.Equal(t, opt.ArchTypes[0], iana.EFI_IA32)
	require.Equal(t, opt.ArchTypes[1], iana.EFI_ITANIUM)
}

func TestParseOptClientArchTypeInvalid(t *testing.T) {
	data := []byte{42}
	_, err := ParseOptClientArchType(data)
	require.Error(t, err)
}

func TestOptClientArchTypeParseAndToBytes(t *testing.T) {
	data := []byte{
		93,   // OptionClientSystemArchitectureType
		2,    // Length
		0, 8, // EFI_XSCALE
	}
	opt, err := ParseOptClientArchType(data)
	require.NoError(t, err)
	require.Equal(t, opt.ToBytes(), data)
}

func TestOptClientArchTypeParseAndToBytesMultiple(t *testing.T) {
	data := []byte{
		93,   // OptionClientSystemArchitectureType
		4,    // Length
		0, 8, // EFI_XSCALE
		0, 6, // EFI_IA32
	}
	opt, err := ParseOptClientArchType(data)
	require.NoError(t, err)
	require.Equal(t, opt.ToBytes(), data)
}

func TestOptClientArchType(t *testing.T) {
	opt := OptClientArchType{
		ArchTypes: []iana.ArchType{iana.EFI_ITANIUM},
	}
	require.Equal(t, opt.Length(), 2)
	require.Equal(t, opt.Code(), OptionClientSystemArchitectureType)
}
