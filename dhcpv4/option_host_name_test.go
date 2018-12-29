package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptHostNameInterfaceMethods(t *testing.T) {
	o := OptHostName{HostName: "foo"}
	require.Equal(t, OptionHostName, o.Code(), "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.ToBytes(), "ToBytes")
}

func TestParseOptHostName(t *testing.T) {
	data := []byte{'t', 'e', 's', 't'}
	o, err := ParseOptHostName(data)
	require.NoError(t, err)
	require.Equal(t, &OptHostName{HostName: "test"}, o)
}

func TestOptHostNameString(t *testing.T) {
	o := OptHostName{HostName: "testy test"}
	require.Equal(t, "Host Name -> testy test", o.String())
}
