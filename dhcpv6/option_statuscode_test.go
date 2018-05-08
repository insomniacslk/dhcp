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

func TestOptStatusCodeToBytes(t *testing.T) {
	expected := []byte{
		0, 13, // OPTION_STATUS_CODE
		0, 9, // length
		0, 0, // StatusSuccess
		's', 'u', 'c', 'c', 'e', 's', 's',
	}
	opt := OptStatusCode{
		iana.StatusSuccess,
		[]byte("success"),
	}
	actual := opt.ToBytes()
	require.Equal(t, expected, actual)
}
