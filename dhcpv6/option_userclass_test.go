package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptUserClass(t *testing.T) {
	expected := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptUserClass(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(opt.UserClasses))
	require.Equal(t, []byte("linuxboot"), opt.UserClasses[0])
}

func TestParseOptUserClassMultiple(t *testing.T) {
	expected := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e', 's', 't',
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

func TestOptUserClassToBytes(t *testing.T) {
	opt := OptUserClass{
		UserClasses: [][]byte{[]byte("linuxboot")},
	}
	data := opt.ToBytes()
	expected := []byte{
		0, 15, // OptionUserClass
		0, 11, // length
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptUserClassToBytesMultiple(t *testing.T) {
	opt := OptUserClass{
		UserClasses: [][]byte{
			[]byte("linuxboot"),
			[]byte("test"),
		},
	}
	data := opt.ToBytes()
	expected := []byte{
		0, 15, // OptionUserClass
		0, 17, // length
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e', 's', 't',
	}
	require.Equal(t, expected, data)
}
