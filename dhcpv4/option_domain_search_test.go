package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptDomainSearch(t *testing.T) {
	data := []byte{
		119, // OptionDNSDomainSearchList
		33,  // length
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
	}
	opt, err := ParseOptDomainSearch(data)
	require.NoError(t, err)
	require.Equal(t, len(opt.DomainSearch), 2)
	require.Equal(t, opt.DomainSearch[0], "example.com")
	require.Equal(t, opt.DomainSearch[1], "subnet.example.org")
}

func TestOptDomainSearchToBytes(t *testing.T) {
	expected := []byte{
		119, // OptionDNSDomainSearchList
		33,  // length
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
	}
	opt := OptDomainSearch{
		DomainSearch: []string{
			"example.com",
			"subnet.example.org",
		},
	}
	require.Equal(t, opt.ToBytes(), expected)
}
