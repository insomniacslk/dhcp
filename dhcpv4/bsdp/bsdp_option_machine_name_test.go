// +build darwin

package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptMachineNameInterfaceMethods(t *testing.T) {
	o := OptMachineName{"somebox"}
	require.Equal(t, OptionMachineName, o.Code(), "Code")
	require.Equal(t, 7, o.Length(), "Length")
	expectedBytes := []byte{130, 7, 's', 'o', 'm', 'e', 'b', 'o', 'x'}
	require.Equal(t, expectedBytes, o.ToBytes(), "ToBytes")
}

func TestParseOptMachineName(t *testing.T) {
	data := []byte{130, 7, 's', 'o', 'm', 'e', 'b', 'o', 'x'}
	o, err := ParseOptMachineName(data)
	require.NoError(t, err)
	require.Equal(t, &OptMachineName{"somebox"}, o)

	// Short byte stream
	data = []byte{130}
	_, err = ParseOptMachineName(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 1, 1}
	_, err = ParseOptMachineName(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{130, 5, 1}
	_, err = ParseOptMachineName(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptMachineNameString(t *testing.T) {
	o := OptMachineName{"somebox"}
	require.Equal(t, "BSDP Machine Name -> somebox", o.String())
}
