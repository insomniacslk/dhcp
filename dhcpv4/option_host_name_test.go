package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptHostNameInterfaceMethods(t *testing.T) {
	o := OptHostName{HostName: "foo"}
	require.Equal(t, OptionHostName, o.Code(), "Code")
	require.Equal(t, 3, o.Length(), "Length")
	require.Equal(t, []byte{byte(OptionHostName), 3, 'f', 'o', 'o'}, o.ToBytes(), "ToBytes")
}

func TestParseOptHostName(t *testing.T) {
	data := []byte{byte(OptionHostName), 4, 't', 'e', 's', 't'}
	o, err := ParseOptHostName(data)
	require.NoError(t, err)
	require.Equal(t, &OptHostName{HostName: "test"}, o)

	// Short byte stream
	data = []byte{byte(OptionHostName)}
	_, err = ParseOptHostName(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 2, 1, 1}
	_, err = ParseOptHostName(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{byte(OptionHostName), 6, 1, 1, 1}
	_, err = ParseOptHostName(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptHostNameString(t *testing.T) {
	o := OptHostName{HostName: "testy test"}
	require.Equal(t, "Host Name -> testy test", o.String())
}
