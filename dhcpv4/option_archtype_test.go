package dhcpv4

import (
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptClientArchType(t *testing.T) {
	m, _ := New(WithGeneric(OptionClientSystemArchitectureType, []byte{
		0, 6, // EFI_IA32
	}))
	archs := m.ClientArch()
	require.NotNil(t, archs)
	require.Equal(t, archs[0], iana.EFI_IA32)
}

func TestParseOptClientArchTypeMultiple(t *testing.T) {
	m, _ := New(WithGeneric(OptionClientSystemArchitectureType, []byte{
		0, 6, // EFI_IA32
		0, 2, // EFI_ITANIUM
	}))
	archs := m.ClientArch()
	require.NotNil(t, archs)
	require.Equal(t, archs[0], iana.EFI_IA32)
	require.Equal(t, archs[1], iana.EFI_ITANIUM)
}

func TestParseOptClientArchTypeInvalid(t *testing.T) {
	m, _ := New(WithGeneric(OptionClientSystemArchitectureType, []byte{42}))
	archs := m.ClientArch()
	require.Nil(t, archs)
}

func TestGetClientArchEmpty(t *testing.T) {
	m, _ := New()
	require.Nil(t, m.ClientArch())
}

func TestOptClientArchTypeParseAndToBytesMultiple(t *testing.T) {
	data := []byte{
		0, 6, // EFI_IA32
		0, 8, // EFI_XSCALE
	}
	opt := OptClientArch(iana.EFI_IA32, iana.EFI_XSCALE)
	require.Equal(t, opt.Value.ToBytes(), data)
	require.Equal(t, opt.Code, OptionClientSystemArchitectureType)
	require.Equal(t, opt.String(), "Client System Architecture Type: EFI IA32, EFI Xscale")
}
