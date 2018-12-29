package dhcpv4

import (
	"testing"

	"github.com/insomniacslk/dhcp/rfc1035label"
	"github.com/stretchr/testify/require"
)

func TestGetDomainSearch(t *testing.T) {
	data := []byte{
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
	}
	o := Options{
		OptionDNSDomainSearchList.Code(): data,
	}
	labels := GetDomainSearch(o)
	require.NotNil(t, labels)
	require.Equal(t, 2, len(labels.Labels))
	require.Equal(t, data, labels.ToBytes())
	require.Equal(t, labels.Labels[0], "example.com")
	require.Equal(t, labels.Labels[1], "subnet.example.org")
}

func TestOptDomainSearchToBytes(t *testing.T) {
	expected := []byte{
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
	}
	opt := OptDomainSearch(&rfc1035label.Labels{
		Labels: []string{
			"example.com",
			"subnet.example.org",
		},
	},
	)
	require.Equal(t, opt.Value.ToBytes(), expected)
}
