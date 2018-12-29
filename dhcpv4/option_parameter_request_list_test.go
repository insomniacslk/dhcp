package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptParameterRequestListInterfaceMethods(t *testing.T) {
	requestedOpts := []OptionCode{OptionBootfileName, OptionNameServer}
	o := &OptParameterRequestList{RequestedOpts: requestedOpts}
	require.Equal(t, OptionParameterRequestList, o.Code(), "Code")

	expectedBytes := []byte{55, 2, 67, 5}
	require.Equal(t, expectedBytes, o.ToBytes(), "ToBytes")

	expectedString := "Parameter Request List -> [Bootfile Name, Name Server]"
	require.Equal(t, expectedString, o.String(), "String")
}

func TestParseOptParameterRequestList(t *testing.T) {
	var (
		o   *OptParameterRequestList
		err error
	)
	o, err = ParseOptParameterRequestList([]byte{67, 5})
	require.NoError(t, err)
	expectedOpts := []OptionCode{OptionBootfileName, OptionNameServer}
	require.Equal(t, expectedOpts, o.RequestedOpts)
}
