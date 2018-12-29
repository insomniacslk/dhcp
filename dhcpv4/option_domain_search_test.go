package dhcpv4

import (
	"testing"

	"github.com/insomniacslk/dhcp/rfc1035label"
	"github.com/stretchr/testify/require"
)

func TestParseOptDomainSearch(t *testing.T) {
	data := []byte{
		119, // OptionDNSDomainSearchList
		33,  // length
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
	}
	opt, err := ParseOptDomainSearch(data[2:])
	require.NoError(t, err)
	require.Equal(t, 2, len(opt.DomainSearch.Labels))
	require.Equal(t, data[2:], opt.DomainSearch.ToBytes())
	require.Equal(t, len(data[2:]), opt.DomainSearch.Length())
	require.Equal(t, opt.DomainSearch.Labels[0], "example.com")
	require.Equal(t, opt.DomainSearch.Labels[1], "subnet.example.org")
}

func TestOptDomainSearchToBytes(t *testing.T) {
	expected := []byte{
		119, // OptionDNSDomainSearchList
		33,  // length
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
	}
	opt := OptDomainSearch{
		DomainSearch: &rfc1035label.Labels{
			Labels: []string{
				"example.com",
				"subnet.example.org",
			},
		},
	}
	require.Equal(t, opt.ToBytes(), expected)
}
