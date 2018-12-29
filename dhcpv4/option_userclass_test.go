package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptUserClassToBytes(t *testing.T) {
	opt := OptRFC3004UserClass([][]byte{[]byte("linuxboot")})
	data := opt.Value.ToBytes()
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptUserClassMicrosoftToBytes(t *testing.T) {
	opt := OptUserClass([]byte("linuxboot"))
	data := opt.Value.ToBytes()
	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestParseOptUserClassMultiple(t *testing.T) {
	var opt UserClass
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		4, 't', 'e', 's', 't',
	}
	err := opt.FromBytes(expected)
	require.NoError(t, err)
	require.Equal(t, len(opt.UserClasses), 2)
	require.Equal(t, []byte("linuxboot"), opt.UserClasses[0])
	require.Equal(t, []byte("test"), opt.UserClasses[1])
}

func TestParseOptUserClassNone(t *testing.T) {
	var opt UserClass
	expected := []byte{}
	err := opt.FromBytes(expected)
	require.Error(t, err)
}

func TestParseOptUserClassMicrosoft(t *testing.T) {
	var opt UserClass
	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	err := opt.FromBytes(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(opt.UserClasses))
	require.Equal(t, []byte("linuxboot"), opt.UserClasses[0])
}

func TestParseOptUserClassMicrosoftShort(t *testing.T) {
	var opt UserClass
	expected := []byte{
		'l',
	}
	err := opt.FromBytes(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(opt.UserClasses))
	require.Equal(t, []byte("l"), opt.UserClasses[0])
}

func TestParseOptUserClass(t *testing.T) {
	var opt UserClass
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	err := opt.FromBytes(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(opt.UserClasses))
	require.Equal(t, []byte("linuxboot"), opt.UserClasses[0])
}

func TestOptUserClassToBytesMultiple(t *testing.T) {
	opt := OptRFC3004UserClass(
		[][]byte{
			[]byte("linuxboot"),
			[]byte("test"),
		},
	)
	data := opt.Value.ToBytes()
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		4, 't', 'e', 's', 't',
	}
	require.Equal(t, expected, data)
}

func TestParseOptUserClassZeroLength(t *testing.T) {
	var opt UserClass
	err := opt.FromBytes([]byte{0, 0})
	require.Error(t, err)
}
