package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptServerPriorityInterfaceMethods(t *testing.T) {
	o := OptServerPriority{Priority: 100}
	require.Equal(t, OptionServerPriority, o.Code(), "Code")
	require.Equal(t, []byte{0, 100}, o.ToBytes(), "ToBytes")
	require.Equal(t, "BSDP Server Priority -> 100", o.String(), "String")
}

func TestParseOptServerPriority(t *testing.T) {
	var (
		o   *OptServerPriority
		err error
	)
	o, err = ParseOptServerPriority([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptServerPriority([]byte{1})
	require.Error(t, err, "short byte stream")

	o, err = ParseOptServerPriority([]byte{0, 100})
	require.NoError(t, err)
	require.Equal(t, uint16(100), o.Priority)
}
