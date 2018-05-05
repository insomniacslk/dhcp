package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptClientArchType(t *testing.T) {
	data := []byte{
		0, 6, // EFI_IA32
	}
	opt, err := ParseOptClientArchType(data)
	require.NoError(t, err)
	require.Equal(t, opt.ArchType, EFI_IA32)
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
		0, 61, // OPTION_CLIENT_ARCH_TYPE
		0, 2, // length
		0, 8, // EFI_XSCALE
	}
	opt, err := ParseOptClientArchType(data)
	require.NoError(t, err)
	require.Equal(t, opt.ToBytes(), expected)
}

func TestOptClientArchType(t *testing.T) {
	opt := OptClientArchType{
		ArchType: EFI_ITANIUM,
	}
	require.Equal(t, opt.Length(), 2)
	require.Equal(t, opt.Code(), OPTION_CLIENT_ARCH_TYPE)
}
