package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOption(t *testing.T) {
	option := []byte{5, 4, 192, 168, 1, 254} // DNS option
	opt, err := ParseOption(option)
	require.NoError(t, err)
	generic := opt.(*OptionGeneric)
	require.Equal(t, OptionNameServer, generic.Code())
	require.Equal(t, []byte{192, 168, 1, 254}, generic.Data)
	require.Equal(t, 4, generic.Length())
	require.Equal(t, "Name Server -> [192 168 1 254]", generic.String())
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
	require.Equal(t, 5, len(opts))
	require.Equal(t, opts[0].(*OptionGeneric), &OptionGeneric{OptionCode: OptionNameServer, Data: []byte{192, 168, 1, 1}})
	require.Equal(t, opts[1].(*OptionGeneric), &OptionGeneric{OptionCode: OptionEnd})
	require.Equal(t, opts[2].(*OptionGeneric), &OptionGeneric{OptionCode: OptionPad})
	require.Equal(t, opts[3].(*OptionGeneric), &OptionGeneric{OptionCode: OptionPad})
	require.Equal(t, opts[4].(*OptionGeneric), &OptionGeneric{OptionCode: OptionPad})
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

func TestOptionsFromBytesShortOption(t *testing.T) {
	options := []byte{
		99, 130, 83, 99, // Magic Cookie
		5, 4, 192, 168, // DNS
	}
	_, err := OptionsFromBytes(options)
	require.Error(t, err)
}
