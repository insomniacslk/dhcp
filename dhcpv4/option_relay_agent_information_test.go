package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptRelayAgentInformation(t *testing.T) {
	data := []byte{
		byte(OptionRelayAgentInformation),
		13,
		1, 5, 'l', 'i', 'n', 'u', 'x',
		2, 4, 'b', 'o', 'o', 't',
	}

	// short sub-option bytes
	opt, err := ParseOptRelayAgentInformation([]byte{1, 0, 1})
	require.Error(t, err)

	// short sub-option length
	opt, err = ParseOptRelayAgentInformation([]byte{1, 1})
	require.Error(t, err)

	opt, err = ParseOptRelayAgentInformation(data[2:])
	require.NoError(t, err)
	require.Equal(t, len(opt.Options), 2)
	circuit := opt.Options.GetOneOption(1).(*OptionGeneric)
	require.NoError(t, err)
	remote := opt.Options.GetOneOption(2).(*OptionGeneric)
	require.NoError(t, err)
	require.Equal(t, circuit.Data, []byte("linux"))
	require.Equal(t, remote.Data, []byte("boot"))
}

func TestParseOptRelayAgentInformationToBytes(t *testing.T) {
	opt := OptRelayAgentInformation{}
	opt1 := &OptionGeneric{OptionCode: 1, Data: []byte("linux")}
	opt.Options = append(opt.Options, opt1)
	opt2 := &OptionGeneric{OptionCode: 2, Data: []byte("boot")}
	opt.Options = append(opt.Options, opt2)
	data := opt.ToBytes()
	expected := []byte{
		byte(OptionRelayAgentInformation),
		13,
		1, 5, 'l', 'i', 'n', 'u', 'x',
		2, 4, 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptRelayAgentInformationToBytesString(t *testing.T) {
	o := OptRelayAgentInformation{}
	require.Equal(t, "Relay Agent Information -> []", o.String())
}
