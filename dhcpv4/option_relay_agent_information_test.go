package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRelayAgentInformation(t *testing.T) {
	m, _ := New(WithGeneric(OptionRelayAgentInformation, []byte{
		1, 5, 'l', 'i', 'n', 'u', 'x',
		2, 4, 'b', 'o', 'o', 't',
	}))

	opt := m.RelayAgentInfo()
	require.NotNil(t, opt)
	require.Equal(t, len(opt.Options), 2)

	circuit := opt.Get(GenericOptionCode(1))
	remote := opt.Get(GenericOptionCode(2))
	require.Equal(t, circuit, []byte("linux"))
	require.Equal(t, remote, []byte("boot"))

	// Empty.
	m, _ = New()
	require.Nil(t, m.RelayAgentInfo())

	// Invalid contents.
	m, _ = New(WithGeneric(OptionRelayAgentInformation, []byte{
		1, 7, 'l', 'i', 'n', 'u', 'x',
	}))
	require.Nil(t, m.RelayAgentInfo())
}

func TestOptRelayAgentInfo(t *testing.T) {
	opt := OptRelayAgentInfo(
		OptGeneric(GenericOptionCode(1), []byte("linux")),
		OptGeneric(GenericOptionCode(2), []byte("boot")),
		OptGeneric(GenericOptionCode(LinkSelectionSubOption), []byte{192, 0, 2, 1}),
	)
	wantBytes := []byte{
		1, 5, 'l', 'i', 'n', 'u', 'x',
		2, 4, 'b', 'o', 'o', 't',
		5, 4, 192, 0, 2, 1,
	}
	wantString := "Relay Agent Information:\n\n    Agent Circuit ID Sub-option: linux ([108 105 110 117 120])\n    Agent Remote ID Sub-option: boot ([98 111 111 116])\n    Link Selection Sub-option: 192.0.2.1\n"
	require.Equal(t, wantBytes, opt.Value.ToBytes())
	require.Equal(t, OptionRelayAgentInformation, opt.Code)
	require.Equal(t, wantString, opt.String())
}
