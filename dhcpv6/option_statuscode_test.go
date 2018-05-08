package dhcpv6

import (
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptStatusCode(t *testing.T) {
	data := []byte{
		0, 5, // StatusUseMulticast
		'u', 's', 'e', ' ', 'm', 'u', 'l', 't', 'i', 'c', 'a', 's', 't',
	}
	opt, err := ParseOptStatusCode(data)
	require.NoError(t, err)
	require.Equal(t, opt.StatusCode, iana.StatusUseMulticast)
	require.Equal(t, opt.StatusMessage, []byte("use multicast"))
}
