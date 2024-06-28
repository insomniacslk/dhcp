package dhcpv4

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptIPAddressLeaseTime(t *testing.T) {
	o := OptIPAddressLeaseTime(43200 * time.Second)
	require.Equal(t, OptionIPAddressLeaseTime, o.Code, "Code")
	require.Equal(t, []byte{0, 0, 168, 192}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "IP Addresses Lease Time: 12h0m0s", o.String(), "String")
}

func TestGetIPAddressLeaseTime(t *testing.T) {
	m, _ := New(WithGeneric(OptionIPAddressLeaseTime, []byte{0, 0, 168, 192}))
	leaseTime := m.IPAddressLeaseTime(0)
	require.Equal(t, 43200*time.Second, leaseTime)

	// Too short.
	m, _ = New(WithGeneric(OptionIPAddressLeaseTime, []byte{168, 192}))
	leaseTime = m.IPAddressLeaseTime(0)
	require.Equal(t, time.Duration(0), leaseTime)

	// Too long.
	m, _ = New(WithGeneric(OptionIPAddressLeaseTime, []byte{1, 1, 1, 1, 1}))
	leaseTime = m.IPAddressLeaseTime(0)
	require.Equal(t, time.Duration(0), leaseTime)

	// Empty.
	m, _ = New()
	require.Equal(t, time.Duration(10), m.IPAddressLeaseTime(10))
}

func TestOptRenewTimeValue(t *testing.T) {
	o := OptRenewTimeValue(43200 * time.Second)
	require.Equal(t, OptionRenewTimeValue, o.Code, "Code")
	require.Equal(t, []byte{0, 0, 168, 192}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Renew Time Value: 12h0m0s", o.String(), "String")
}

func TestGetRenewTimeValue(t *testing.T) {
	m, _ := New(WithGeneric(OptionRenewTimeValue, []byte{0, 0, 168, 192}))
	leaseTime := m.IPAddressRenewalTime(0)
	require.Equal(t, 43200*time.Second, leaseTime)

	// Too short.
	m, _ = New(WithGeneric(OptionRenewTimeValue, []byte{168, 192}))
	leaseTime = m.IPAddressRenewalTime(0)
	require.Equal(t, time.Duration(0), leaseTime)

	// Too long.
	m, _ = New(WithGeneric(OptionRenewTimeValue, []byte{1, 1, 1, 1, 1}))
	leaseTime = m.IPAddressRenewalTime(0)
	require.Equal(t, time.Duration(0), leaseTime)

	// Empty.
	m, _ = New()
	require.Equal(t, time.Duration(10), m.IPAddressRenewalTime(10))
}

func TestOptRebindingTimeValue(t *testing.T) {
	o := OptRebindingTimeValue(43200 * time.Second)
	require.Equal(t, OptionRebindingTimeValue, o.Code, "Code")
	require.Equal(t, []byte{0, 0, 168, 192}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Rebinding Time Value: 12h0m0s", o.String(), "String")
}

func TestGetRebindingTimeValue(t *testing.T) {
	m, _ := New(WithGeneric(OptionRebindingTimeValue, []byte{0, 0, 168, 192}))
	leaseTime := m.IPAddressRebindingTime(0)
	require.Equal(t, 43200*time.Second, leaseTime)

	// Too short.
	m, _ = New(WithGeneric(OptionRebindingTimeValue, []byte{168, 192}))
	leaseTime = m.IPAddressRebindingTime(0)
	require.Equal(t, time.Duration(0), leaseTime)

	// Too long.
	m, _ = New(WithGeneric(OptionRebindingTimeValue, []byte{1, 1, 1, 1, 1}))
	leaseTime = m.IPAddressRebindingTime(0)
	require.Equal(t, time.Duration(0), leaseTime)

	// Empty.
	m, _ = New()
	require.Equal(t, time.Duration(10), m.IPAddressRebindingTime(10))
}

func TestOptIPv6OnlyPreferred(t *testing.T) {
	o := OptIPv6OnlyPreferred(43200 * time.Second)
	require.Equal(t, OptionIPv6OnlyPreferred, o.Code, "Code")
	require.Equal(t, []byte{0, 0, 168, 192}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "IPv6-Only Preferred: 12h0m0s", o.String(), "String")
}

func TestOptIPv6OnlyPreferredZero(t *testing.T) {
	o := OptIPv6OnlyPreferred(0)
	require.Equal(t, OptionIPv6OnlyPreferred, o.Code, "Code")
	require.Equal(t, []byte{0, 0, 0, 0}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "IPv6-Only Preferred: 0s", o.String(), "String")
}

func TestGetIPv6OnlyPreferred(t *testing.T) {
	m, _ := New(WithGeneric(OptionIPv6OnlyPreferred, []byte{0, 0, 168, 192}))
	v6onlyWait, ok := m.IPv6OnlyPreferred()
	require.True(t, ok)
	require.Equal(t, 43200*time.Second, v6onlyWait)

	// Too short.
	m, _ = New(WithGeneric(OptionIPv6OnlyPreferred, []byte{168, 192}))
	_, ok = m.IPv6OnlyPreferred()
	require.False(t, ok)

	// Too long.
	m, _ = New(WithGeneric(OptionIPv6OnlyPreferred, []byte{1, 1, 1, 1, 1}))
	_, ok = m.IPv6OnlyPreferred()
	require.False(t, ok)

	// Missing.
	m, _ = New()
	_, ok = m.IPv6OnlyPreferred()
	require.False(t, ok)
}
