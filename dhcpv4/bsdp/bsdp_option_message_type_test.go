package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptMessageTypeInterfaceMethods(t *testing.T) {
	o := OptMessageType{MessageTypeList}
	require.Equal(t, OptionMessageType, o.Code(), "Code")
	require.Equal(t, []byte{1}, o.ToBytes(), "ToBytes")
}

func TestParseOptMessageType(t *testing.T) {
	data := []byte{1} // DISCOVER
	o, err := ParseOptMessageType(data)
	require.NoError(t, err)
	require.Equal(t, &OptMessageType{MessageTypeList}, o)
}

func TestOptMessageTypeString(t *testing.T) {
	// known
	o := OptMessageType{MessageTypeList}
	require.Equal(t, "BSDP Message Type -> LIST", o.String())

	// unknown
	o = OptMessageType{99}
	require.Equal(t, "BSDP Message Type -> Unknown", o.String())
}
