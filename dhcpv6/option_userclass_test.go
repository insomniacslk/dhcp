package dhcpv6

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptUserClass(t *testing.T) {
	expected := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptUserClass(expected)
	require.Nil(t, err)
	log.Printf("%+v", opt)
	require.Equal(t, len(opt.UserClasses), 1)
	require.Equal(t, opt.UserClasses[0], []byte("linuxboot"))
}
