package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptReplyPortInterfaceMethods(t *testing.T) {
	o := OptReplyPort{1234}
	require.Equal(t, OptionReplyPort, o.Code(), "Code")
	require.Equal(t, 2, o.Length(), "Length")
	require.Equal(t, []byte{byte(OptionReplyPort), 2, 4, 210}, o.ToBytes(), "ToBytes")
}

func TestParseOptReplyPort(t *testing.T) {
	data := []byte{byte(OptionReplyPort), 2, 0, 1}
	o, err := ParseOptReplyPort(data)
	require.NoError(t, err)
	require.Equal(t, &OptReplyPort{1}, o)

	// Short byte stream
	data = []byte{byte(OptionReplyPort), 2}
	_, err = ParseOptReplyPort(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 2, 1, 0}
	_, err = ParseOptReplyPort(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{byte(OptionReplyPort), 4, 1, 0}
	_, err = ParseOptReplyPort(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptReplyPortString(t *testing.T) {
	// known
	o := OptReplyPort{1234}
	require.Equal(t, "BSDP Reply Port -> 1234", o.String())
}
