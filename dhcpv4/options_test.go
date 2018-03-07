package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOption(t *testing.T) {
	option := []byte{5, 4, 192, 168, 1, 254} // DNS option
	opt, err := ParseOption(option)
	require.NoError(t, err, "should not get error from parsing option")
	require.Equal(t, OptionNameServer, opt.Code, "opt should have the same opcode")
	require.Equal(t, option[2:], opt.Data, "opt should have the same data")
}

func TestParseOptionPad(t *testing.T) {
	option := []byte{0}
	opt, err := ParseOption(option)
	require.NoError(t, err, "should not get error from parsing option")
	require.Equal(t, OptionPad, opt.Code, "should get pad option code")
	require.Empty(t, opt.Data, "should get empty data with pad option")
}

func TestParseOptionZeroLength(t *testing.T) {
	option := []byte{}
	_, err := ParseOption(option)
	require.Error(t, err, "should get error from zero-length options")
}

func TestParseOptionShortOption(t *testing.T) {
	option := []byte{53, 1}
	_, err := ParseOption(option)
	require.Error(t, err, "should get error from short options")
}

func TestOptionsFromBytes(t *testing.T) {
	options := []byte{
		99, 130, 83, 99, // Magic Cookie
		5, 4, 192, 168, 1, 1, // DNS
		255,     // end
		0, 0, 0, //padding
	}
	opts, err := OptionsFromBytes(options)
	require.NoError(t, err)
	require.Equal(t, []Option{
		Option{
			Code: OptionNameServer,
			Data: []byte{192, 168, 1, 1},
		},
		Option{Code: OptionEnd, Data: []byte{}},
		Option{Code: OptionPad, Data: []byte{}},
		Option{Code: OptionPad, Data: []byte{}},
		Option{Code: OptionPad, Data: []byte{}},
	}, opts)
}

func TestOptionsFromBytesZeroLength(t *testing.T) {
	options := []byte{}
	_, err := OptionsFromBytes(options)
	require.Error(t, err)
}

func TestOptionsFromBytesBadMagicCookie(t *testing.T) {
	options := []byte{1, 2, 3, 4}
	_, err := OptionsFromBytes(options)
	require.Error(t, err)
}

func TestOptionsToBytes(t *testing.T) {
	originalOptions := []byte{
		99, 130, 83, 99, // Magic Cookie
		5, 4, 192, 168, 1, 1, // DNS
		255,     // end
		0, 0, 0, //padding
	}
	options, err := OptionsFromBytes(originalOptions)
	require.NoError(t, err)
	finalOptions := OptionsToBytes(options)
	require.Equal(t, originalOptions, finalOptions)
}

func TestOptionsToBytesEmpty(t *testing.T) {
	originalOptions := []byte{99, 130, 83, 99}
	options, err := OptionsFromBytes(originalOptions)
	require.NoError(t, err)
	finalOptions := OptionsToBytes(options)
	require.Equal(t, originalOptions, finalOptions)
}

func TestOptionsToStringPad(t *testing.T) {
	option := []byte{0}
	opt, err := ParseOption(option)
	require.NoError(t, err)
	stropt := opt.String()
	require.Equal(t, "Pad -> []", stropt)
}

func TestOptionsToStringDHCPMessageType(t *testing.T) {
	option := []byte{53, 1, 5}
	opt, err := ParseOption(option)
	require.NoError(t, err)
	stropt := opt.String()
	require.Equal(t, "DHCP Message Type -> [5]", stropt)
}
