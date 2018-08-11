package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserClassModifier(t *testing.T) {
	d, _ := New()
	userClass := WithUserClass([]byte("linuxboot"), false)
	d = userClass(d)
	expected := []byte{
		77, // OptionUserClass
		9,  // length
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, "User Class Information -> linuxboot", d.options[0].String())
	require.Equal(t, expected, d.options[0].ToBytes())
}

func TestUserClassModifierRFC(t *testing.T) {
	d, _ := New()
	userClass := WithUserClass([]byte("linuxboot"), true)
	d = userClass(d)
	expected := []byte{
		77, // OptionUserClass
		10, // length
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, "User Class Information -> linuxboot", d.options[0].String())
	require.Equal(t, expected, d.options[0].ToBytes())
}

func TestWithNetboot(t *testing.T) {
	d, _ := New()
	d = WithNetboot(d)
	require.Equal(t, "Parameter Request List -> [TFTP Server Name, Bootfile Name]", d.options[0].String())
}

func TestWithNetbootExistingTFTP(t *testing.T) {
	d, _ := New()
	OptParams := &OptParameterRequestList{
		RequestedOpts: []OptionCode{OptionTFTPServerName},
	}
	d.AddOption(OptParams)
	d = WithNetboot(d)
	require.Equal(t, "Parameter Request List -> [TFTP Server Name, Bootfile Name]", d.options[0].String())
}

func TestWithNetbootExistingBootfileName(t *testing.T) {
	d, _ := New()
	OptParams := &OptParameterRequestList{
		RequestedOpts: []OptionCode{OptionBootfileName},
	}
	d.AddOption(OptParams)
	d = WithNetboot(d)
	require.Equal(t, "Parameter Request List -> [Bootfile Name, TFTP Server Name]", d.options[0].String())
}

func TestWithNetbootExistingBoth(t *testing.T) {
	d, _ := New()
	OptParams := &OptParameterRequestList{
		RequestedOpts: []OptionCode{OptionBootfileName, OptionTFTPServerName},
	}
	d.AddOption(OptParams)
	d = WithNetboot(d)
	require.Equal(t, "Parameter Request List -> [Bootfile Name, TFTP Server Name]", d.options[0].String())
}
