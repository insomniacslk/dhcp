package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptDomainNameInterfaceMethods(t *testing.T) {
	o := OptDomainName{DomainName: "foo"}
	require.Equal(t, OptionDomainName, o.Code(), "Code")
	require.Equal(t, 3, o.Length(), "Length")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.ToBytes(), "ToBytes")
}

func TestParseOptDomainName(t *testing.T) {
	data := []byte{'t', 'e', 's', 't'}
	o, err := ParseOptDomainName(data)
	require.NoError(t, err)
	require.Equal(t, &OptDomainName{DomainName: "test"}, o)
}

func TestOptDomainNameString(t *testing.T) {
	o := OptDomainName{DomainName: "testy test"}
	require.Equal(t, "Domain Name -> testy test", o.String())
}
