package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsRequetsed(t *testing.T) {
	pkt, err := New()
	require.NoError(t, err)
	require.False(t, IsRequested(pkt, OptionDomainNameServer))

	optprl := OptParameterRequestList{RequestedOpts: []OptionCode{OptionDomainNameServer}}
	pkt.AddOption(&optprl)
	require.True(t, IsRequested(pkt, OptionDomainNameServer))
}
