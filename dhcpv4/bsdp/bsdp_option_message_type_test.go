package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptMessageTypeInterfaceMethods(t *testing.T) {
	o := OptMessageType{MessageTypeList}
	require.Equal(t, OptionMessageType, o.Code(), "Code")
	require.Equal(t, 1, o.Length(), "Length")
	require.Equal(t, []byte{1, 1, 1}, o.ToBytes(), "ToBytes")
}

func TestParseOptMessageType(t *testing.T) {
	data := []byte{1, 1, 1} // DISCOVER
	o, err := ParseOptMessageType(data)
	require.NoError(t, err)
	require.Equal(t, &OptMessageType{MessageTypeList}, o)

	// Short byte stream
	data = []byte{1, 1}
	_, err = ParseOptMessageType(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 1, 1}
	_, err = ParseOptMessageType(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{1, 5, 1}
	_, err = ParseOptMessageType(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptMessageTypeString(t *testing.T) {
	// known
	o := OptMessageType{MessageTypeList}
	require.Equal(t, "BSDP Message Type -> LIST", o.String())

	// unknown
	o = OptMessageType{99}
	require.Equal(t, "BSDP Message Type -> Unknown", o.String())
}
