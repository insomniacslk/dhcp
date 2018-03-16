package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptClassIdentifierInterfaceMethods(t *testing.T) {
	o := OptClassIdentifier{Identifier: "foo"}
	require.Equal(t, OptionClassIdentifier, o.Code(), "Code")
	require.Equal(t, 3, o.Length(), "Length")
	require.Equal(t, []byte{byte(OptionClassIdentifier), 3, 'f', 'o', 'o'}, o.ToBytes(), "ToBytes")
}

func TestParseOptClassIdentifier(t *testing.T) {
	data := []byte{byte(OptionClassIdentifier), 4, 't', 'e', 's', 't'} // DISCOVER
	o, err := ParseOptClassIdentifier(data)
	require.NoError(t, err)
	require.Equal(t, &OptClassIdentifier{Identifier: "test"}, o)

	// Short byte stream
	data = []byte{byte(OptionClassIdentifier)}
	_, err = ParseOptClassIdentifier(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 2, 1, 1}
	_, err = ParseOptClassIdentifier(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{byte(OptionClassIdentifier), 6, 1, 1, 1}
	_, err = ParseOptClassIdentifier(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptClassIdentifierString(t *testing.T) {
	o := OptClassIdentifier{Identifier: "testy test"}
	require.Equal(t, "Class Identifier -> testy test", o.String())
}
