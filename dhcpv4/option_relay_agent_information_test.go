package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptRelayAgentInformation(t *testing.T) {
	data := []byte{
		1, 5, 'l', 'i', 'n', 'u', 'x',
		2, 4, 'b', 'o', 'o', 't',
	}

	// short sub-option bytes
	opt, err := ParseOptRelayAgentInformation([]byte{1, 0, 1})
	require.Error(t, err)

	// short sub-option length
	opt, err = ParseOptRelayAgentInformation([]byte{1, 1})
	require.Error(t, err)

	opt, err = ParseOptRelayAgentInformation(data)
	require.NoError(t, err)
	require.Equal(t, len(opt.Options), 2)
	circuit := opt.Options.GetOne(1).(*OptionGeneric)
	require.NoError(t, err)
	remote := opt.Options.GetOne(2).(*OptionGeneric)
	require.NoError(t, err)
	require.Equal(t, circuit.Data, []byte("linux"))
	require.Equal(t, remote.Data, []byte("boot"))
}

func TestParseOptRelayAgentInformationToBytes(t *testing.T) {
	opt := OptRelayAgentInformation{
		Options: Options{
			&OptionGeneric{OptionCode: 1, Data: []byte("linux")},
			&OptionGeneric{OptionCode: 2, Data: []byte("boot")},
		},
	}
	data := opt.ToBytes()
	expected := []byte{
		1, 5, 'l', 'i', 'n', 'u', 'x',
		2, 4, 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptRelayAgentInformationToBytesString(t *testing.T) {
	o := OptRelayAgentInformation{}
	require.Equal(t, "Relay Agent Information -> []", o.String())
}
