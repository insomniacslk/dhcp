// +build darwin

package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptServerPriorityInterfaceMethods(t *testing.T) {
	o := OptServerPriority{Priority: 100}
	require.Equal(t, OptionServerPriority, o.Code(), "Code")
	require.Equal(t, []byte{4, 2, 0, 100}, o.ToBytes(), "ToBytes")
	require.Equal(t, 2, o.Length(), "Length")
	require.Equal(t, "BSDP Server Priority -> 100", o.String(), "String")
}

func TestParseOptServerPriority(t *testing.T) {
	var (
		o   *OptServerPriority
		err error
	)
	o, err = ParseOptServerPriority([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptServerPriority([]byte{4, 2, 1})
	require.Error(t, err, "short byte stream")

	o, err = ParseOptServerPriority([]byte{4, 3, 1, 1})
	require.Error(t, err, "wrong priority length")

	o, err = ParseOptServerPriority([]byte{53, 2, 168, 1})
	require.Error(t, err, "wrong option code")

	o, err = ParseOptServerPriority([]byte{4, 2, 0, 100})
	require.NoError(t, err)
	require.Equal(t, 100, o.Priority)
}
