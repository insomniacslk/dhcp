package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptDomainSearchList(t *testing.T) {
	data := []byte{
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
	}
	opt, err := ParseOptDomainSearchList(data)
	require.NoError(t, err)
	require.Equal(t, len(opt.DomainSearchList), 2)
	require.Equal(t, opt.DomainSearchList[0], "example.com")
	require.Equal(t, opt.DomainSearchList[1], "subnet.example.org")
}

func TestOptDomainSearchListToBytes(t *testing.T) {
	expected := []byte{
		0, 24, // OptionDomainSearchList
		0, 33, // length
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
	}
	opt := OptDomainSearchList{
		DomainSearchList: []string{
			"example.com",
			"subnet.example.org",
		},
	}
	require.Equal(t, opt.ToBytes(), expected)
}
