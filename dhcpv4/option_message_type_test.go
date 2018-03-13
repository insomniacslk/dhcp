package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptMessageTypeInterfaceMethods(t *testing.T) {
	o := OptMessageType{MessageType: MessageTypeDiscover}
	require.Equal(t, OptionDHCPMessageType, o.Code(), "Code")
	require.Equal(t, 1, o.Length(), "Length")
	require.Equal(t, []byte{53, 1, 1}, o.ToBytes(), "ToBytes")
}

func TestOptMessageTypeNew(t *testing.T) {
	o := OptMessageType{MessageType: MessageTypeDiscover}
	require.Equal(t, OptionDHCPMessageType, o.Code())
	require.Equal(t, MessageTypeDiscover, o.MessageType)
}

func TestParseOptMessageType(t *testing.T) {
	data := []byte{53, 1, 1} // DISCOVER
	o, err := ParseOptMessageType(data)
	require.NoError(t, err)
	require.Equal(t, &OptMessageType{MessageType: MessageTypeDiscover}, o)

	// Short byte stream
	data = []byte{53, 1}
	_, err = ParseOptMessageType(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 1, 1}
	_, err = ParseOptMessageType(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{53, 5, 1}
	_, err = ParseOptMessageType(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptMessageTypeString(t *testing.T) {
	var o OptMessageType
	o = OptMessageType{MessageType: MessageTypeDiscover}
	require.Equal(t, "DHCP Message Type -> DISCOVER", o.String())
	o = OptMessageType{MessageType: MessageTypeOffer}
	require.Equal(t, "DHCP Message Type -> OFFER", o.String())
	o = OptMessageType{MessageType: MessageTypeRequest}
	require.Equal(t, "DHCP Message Type -> REQUEST", o.String())
	o = OptMessageType{MessageType: MessageTypeDecline}
	require.Equal(t, "DHCP Message Type -> DECLINE", o.String())
	o = OptMessageType{MessageType: MessageTypeAck}
	require.Equal(t, "DHCP Message Type -> ACK", o.String())
	o = OptMessageType{MessageType: MessageTypeNak}
	require.Equal(t, "DHCP Message Type -> NAK", o.String())
	o = OptMessageType{MessageType: MessageTypeRelease}
	require.Equal(t, "DHCP Message Type -> RELEASE", o.String())
	o = OptMessageType{MessageType: MessageTypeInform}
	require.Equal(t, "DHCP Message Type -> INFORM", o.String())
	o = OptMessageType{MessageType: 99}
	require.Equal(t, "DHCP Message Type -> UNKNOWN", o.String())
}
