package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptVersionInterfaceMethods(t *testing.T) {
	o := OptVersion{Version1_1}
	require.Equal(t, OptionVersion, o.Code(), "Code")
	require.Equal(t, 2, o.Length(), "Length")
	require.Equal(t, []byte{2, 2, 1, 1}, o.ToBytes(), "ToBytes")
}

func TestParseOptVersion(t *testing.T) {
	data := []byte{2, 2, 1, 1}
	o, err := ParseOptVersion(data)
	require.NoError(t, err)
	require.Equal(t, &OptVersion{Version1_1}, o)

	// Short byte stream
	data = []byte{2, 2}
	_, err = ParseOptVersion(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 2, 1, 0}
	_, err = ParseOptVersion(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{2, 4, 1, 0}
	_, err = ParseOptVersion(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptVersionString(t *testing.T) {
	// known
	o := OptVersion{Version1_1}
	require.Equal(t, "BSDP Version -> 1.1", o.String())
}
