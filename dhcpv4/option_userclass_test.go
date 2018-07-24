package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptUserClassToBytes(t *testing.T) {
	opt := OptUserClass{
		UserClasses: [][]byte{[]byte("linuxboot")},
	}
	data := opt.ToBytes()
	expected := []byte{
		77, // OPTION_USER_CLASS
		10, // length
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestParseOptUserClassMultiple(t *testing.T) {
	expected := []byte{
		77, 15,
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

func TestParseOptUserClass(t *testing.T) {
	expected := []byte{
		77, 10, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
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
	}
	data := opt.ToBytes()
	expected := []byte{
		77, // OPTION_USER_CLASS
		15, // length
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		4, 't', 'e', 's', 't',
	}
	require.Equal(t, expected, data)
}
