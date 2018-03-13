package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptParameterRequestListInterfaceMethods(t *testing.T) {
	requestedOpts := []OptionCode{OptionBootfileName, OptionNameServer}
	o := NewOptParameterRequestList(requestedOpts...)
	require.Equal(t, OptionParameterRequestList, o.Code(), "Code")
	require.Equal(t, requestedOpts, o.RequestList(), "RequestList")

	expectedBytes := []byte{55, 2, 67, 5}
	require.Equal(t, expectedBytes, o.ToBytes(), "ToBytes")

	expectedString := "Parameter Request List -> [67 5]"
	require.Equal(t, expectedString, o.String(), "String")
}

func TestParseOptParameterRequestList(t *testing.T) {
	var (
		o   *OptParameterRequestList
		err error
	)
	o, err = ParseOptParameterRequestList([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptParameterRequestList([]byte{55, 2})
	require.Error(t, err, "short byte stream")

	o, err = ParseOptParameterRequestList([]byte{53, 2, 1, 1})
	require.Error(t, err, "wrong option code")

	o, err = ParseOptParameterRequestList([]byte{55, 2, 67, 5})
	require.NoError(t, err)
	expectedOpts := []OptionCode{OptionBootfileName, OptionNameServer}
	require.Equal(t, expectedOpts, o.RequestList())
}
