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
	require.Nil(t, err)
	require.Equal(t, len(opt.UserClasses), 1)
	require.Equal(t, opt.UserClasses[0], []byte("linuxboot"))
}

func TestParseOptUserClassMultiple(t *testing.T) {
	expected := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e', 's', 't',
	}
	opt, err := ParseOptUserClass(expected)
	require.Nil(t, err)
	require.Equal(t, len(opt.UserClasses), 2)
	require.Equal(t, opt.UserClasses[0], []byte("linuxboot"))
	require.Equal(t, opt.UserClasses[1], []byte("test"))
}

func TestParseOptUserClassNone(t *testing.T) {
	expected := []byte{}
	opt, err := ParseOptUserClass(expected)
	require.Nil(t, err)
	require.Equal(t, len(opt.UserClasses), 0)
}

func TestOptUserClassToBytes(t *testing.T) {
	opt := OptUserClass{
		UserClasses: [][]byte{[]byte("linuxboot")},
	}
	data := opt.ToBytes()
	expected := []byte{
		0, 15, // OPTION_USER_CLASS
		0, 11, // length
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, data, expected)
}
