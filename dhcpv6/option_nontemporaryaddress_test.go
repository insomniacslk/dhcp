package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptIANAParseOptIANA(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, 0, 0, 0xb2, 0x7a, 0, 0, 0xc0, 0x8a, // options
	}
	opt, err := ParseOptIANA(data)
	require.NoError(t, err)
	require.Equal(t, len(data), opt.Length())
}

func TestOptIANAParseOptIANAInvalidLength(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		// truncated from here
	}
	_, err := ParseOptIANA(data)
	require.Error(t, err)
}

func TestOptIANAParseOptIANAInvalidOptions(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, 0, 0, 0xb2, 0x7a, // truncated options
	}
	_, err := ParseOptIANA(data)
	require.Error(t, err)
}

func TestOptIANADelOption(t *testing.T) {
	optiana := OptIANA{}
	optiaadr := OptIAAddress{}
	optiana.Options = append(optiana.Options, &optiaadr)
	optiana.DelOption(OPTION_IAADDR)
	opt := getOption(optiana.Options, OPTION_IAADDR)
	require.Nil(t, opt)
}

func TestOptIANAToBytes(t *testing.T) {
	opt := OptIANA{
		IaId: [4]byte{1, 2, 3, 4},
		T1:   12345,
		T2:   54321,
		Options: []Option{
			&OptElapsedTime{
				ElapsedTime: 0xaabb,
			},
		},
	}
	expected := []byte{
		0, 3, // OPTION_IA_NA
		0, 18, // length
		1, 2, 3, 4, // IA ID
		0, 0, 0x30, 0x39, // T1 = 12345
		0, 0, 0xd4, 0x31, // T2 = 54321
		0, 8, 0, 2, 0xaa, 0xbb,
	}
	require.Equal(t, expected, opt.ToBytes())
}
