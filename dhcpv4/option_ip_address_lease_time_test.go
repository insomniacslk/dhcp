package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptIPAddressLeaseTimeInterfaceMethods(t *testing.T) {
	o := OptIPAddressLeaseTime{LeaseTime: 43200}
	require.Equal(t, OptionIPAddressLeaseTime, o.Code(), "Code")
	require.Equal(t, 4, o.Length(), "Length")
	require.Equal(t, []byte{51, 4, 0, 0, 168, 192}, o.ToBytes(), "ToBytes")
}

func TestParseOptIPAddressLeaseTime(t *testing.T) {
	data := []byte{51, 4, 0, 0, 168, 192}
	o, err := ParseOptIPAddressLeaseTime(data)
	require.NoError(t, err)
	require.Equal(t, &OptIPAddressLeaseTime{LeaseTime: 43200}, o)

	// Short byte stream
	data = []byte{51, 4}
	_, err = ParseOptIPAddressLeaseTime(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 4, 0, 0, 168, 192}
	_, err = ParseOptIPAddressLeaseTime(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{51, 5, 1, 1, 1, 1, 1}
	_, err = ParseOptMaximumDHCPMessageSize(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptIPAddressLeaseTimeString(t *testing.T) {
	o := OptIPAddressLeaseTime{LeaseTime: 43200}
	require.Equal(t, "IP Addresses Lease Time -> 43200", o.String())
}
