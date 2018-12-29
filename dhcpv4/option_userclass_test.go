package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptUserClassToBytes(t *testing.T) {
	opt := OptUserClass{
		UserClasses: [][]byte{[]byte("linuxboot")},
		Rfc3004:     true,
	}
	data := opt.ToBytes()
	expected := []byte{
		77, // OptionUserClass
		10, // length
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptUserClassMicrosoftToBytes(t *testing.T) {
	opt := OptUserClass{
		UserClasses: [][]byte{[]byte("linuxboot")},
	}
	data := opt.ToBytes()
	expected := []byte{
		77, // OptionUserClass
		9,  // length
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestParseOptUserClassMultiple(t *testing.T) {
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		4, 't', 'e', 's', 't',
	}
	opt, err := ParseOptUserClass(expected)
	require.NoError(t, err)
	require.Equal(t, len(opt.UserClasses), 2)
	require.Equal(t, []byte("linuxboot"), opt.UserClasses[0])
	require.Equal(t, []byte("test"), opt.UserClasses[1])
}

func TestParseOptUserClassNone(t *testing.T) {
	expected := []byte{}
	_, err := ParseOptUserClass(expected)
	require.Error(t, err)
}

func TestParseOptUserClassMicrosoft(t *testing.T) {
	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptUserClass(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(opt.UserClasses))
	require.Equal(t, []byte("linuxboot"), opt.UserClasses[0])
}

func TestParseOptUserClassMicrosoftShort(t *testing.T) {
	expected := []byte{
		'l',
	}
	opt, err := ParseOptUserClass(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(opt.UserClasses))
	require.Equal(t, []byte("l"), opt.UserClasses[0])
}

func TestParseOptUserClass(t *testing.T) {
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptUserClass(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(opt.UserClasses))
	require.Equal(t, []byte("linuxboot"), opt.UserClasses[0])
}

func TestOptUserClassToBytesMultiple(t *testing.T) {
	opt := OptUserClass{
		UserClasses: [][]byte{
			[]byte("linuxboot"),
			[]byte("test"),
		},
		Rfc3004: true,
	}
	data := opt.ToBytes()
	expected := []byte{
		77, // OptionUserClass
		15, // length
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		4, 't', 'e', 's', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptUserClassLength(t *testing.T) {
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't', 'X',
	}
	opt, err := ParseOptUserClass(expected)
	require.NoError(t, err)
	require.Equal(t, 11, opt.Length())
}

func TestParseOptUserClassZeroLength(t *testing.T) {
	expected := []byte{
		0, 0,
	}
	_, err := ParseOptUserClass(expected)
	require.Error(t, err)
}

func TestOptUserClassCode(t *testing.T) {
	opt := OptUserClass{}
	require.Equal(t, OptionUserClassInformation, opt.Code())
}
