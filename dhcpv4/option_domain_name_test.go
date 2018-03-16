package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptDomainNameInterfaceMethods(t *testing.T) {
	o := OptDomainName{DomainName: "foo"}
	require.Equal(t, OptionDomainName, o.Code(), "Code")
	require.Equal(t, 3, o.Length(), "Length")
	require.Equal(t, []byte{byte(OptionDomainName), 3, 'f', 'o', 'o'}, o.ToBytes(), "ToBytes")
}

func TestParseOptDomainName(t *testing.T) {
	data := []byte{byte(OptionDomainName), 4, 't', 'e', 's', 't'} // DISCOVER
	o, err := ParseOptDomainName(data)
	require.NoError(t, err)
	require.Equal(t, &OptDomainName{DomainName: "test"}, o)

	// Short byte stream
	data = []byte{byte(OptionDomainName)}
	_, err = ParseOptDomainName(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 2, 1, 1}
	_, err = ParseOptDomainName(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{byte(OptionDomainName), 6, 1, 1, 1}
	_, err = ParseOptDomainName(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptDomainNameString(t *testing.T) {
	o := OptDomainName{DomainName: "testy test"}
	require.Equal(t, "Domain Name -> testy test", o.String())
}
