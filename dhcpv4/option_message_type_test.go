package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptMessageTypeInterfaceMethods(t *testing.T) {
	o := OptMessageType{MessageType: MessageTypeDiscover}
	require.Equal(t, OptionDHCPMessageType, o.Code(), "Code")
	require.Equal(t, []byte{1}, o.ToBytes(), "ToBytes")
}

func TestOptMessageTypeNew(t *testing.T) {
	o := OptMessageType{MessageType: MessageTypeDiscover}
	require.Equal(t, OptionDHCPMessageType, o.Code())
	require.Equal(t, MessageTypeDiscover, o.MessageType)
}

func TestParseOptMessageType(t *testing.T) {
	data := []byte{1} // DISCOVER
	o, err := ParseOptMessageType(data)
	require.NoError(t, err)
	require.Equal(t, &OptMessageType{MessageType: MessageTypeDiscover}, o)

	// Bad length
	data = []byte{1, 2}
	_, err = ParseOptMessageType(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptMessageTypeString(t *testing.T) {
	// known
	o := OptMessageType{MessageType: MessageTypeDiscover}
	require.Equal(t, "DHCP Message Type -> DISCOVER", o.String())

	// unknown
	o = OptMessageType{MessageType: 99}
	require.Equal(t, "DHCP Message Type -> Unknown", o.String())
}
