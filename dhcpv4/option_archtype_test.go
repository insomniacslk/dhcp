package dhcpv4

import (
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptClientArchType(t *testing.T) {
	o := Options{OptionClientSystemArchitectureType.Code(): []byte{
		0, 6, // EFI_IA32
	}}
	archs := GetClientArch(o)
	require.NotNil(t, archs)
	require.Equal(t, archs[0], iana.EFI_IA32)
}

func TestParseOptClientArchTypeMultiple(t *testing.T) {
	o := Options{OptionClientSystemArchitectureType.Code(): []byte{
		0, 6, // EFI_IA32
		0, 2, // EFI_ITANIUM
	}}
	archs := GetClientArch(o)
	require.NotNil(t, archs)
	require.Equal(t, archs[0], iana.EFI_IA32)
	require.Equal(t, archs[1], iana.EFI_ITANIUM)
}

func TestParseOptClientArchTypeInvalid(t *testing.T) {
	o := Options{OptionClientSystemArchitectureType.Code(): []byte{42}}
	archs := GetClientArch(o)
	require.Nil(t, archs)
}

func TestGetClientArchEmpty(t *testing.T) {
	require.Nil(t, GetClientArch(Options{}))
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
