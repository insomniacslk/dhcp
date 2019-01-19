package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionIDModifier(t *testing.T) {
	d, err := New()
	require.NoError(t, err)
	WithTransactionID(TransactionID{0xdd, 0xcc, 0xbb, 0xaa})(d)
	require.Equal(t, TransactionID{0xdd, 0xcc, 0xbb, 0xaa}, d.TransactionID)
}

func TestBroadcastModifier(t *testing.T) {
	d, err := New()
	require.NoError(t, err)

	// set and test broadcast
	WithBroadcast(true)(d)
	require.Equal(t, true, d.IsBroadcast())

	// set and test unicast
	WithBroadcast(false)(d)
	require.Equal(t, true, d.IsUnicast())
}

func TestHwAddrModifier(t *testing.T) {
	hwaddr := net.HardwareAddr{0xa, 0xb, 0xc, 0xd, 0xe, 0xf}
	d, err := New(WithHwAddr(hwaddr))
	require.NoError(t, err)
	require.Equal(t, hwaddr, d.ClientHWAddr)
}

func TestWithOptionModifier(t *testing.T) {
	d, err := New(WithOption(&OptDomainName{DomainName: "slackware.it"}))
	require.NoError(t, err)

	opt := d.GetOneOption(OptionDomainName)
	require.NotNil(t, opt)
	dnOpt := opt.(*OptDomainName)
	require.Equal(t, "slackware.it", dnOpt.DomainName)
}

func TestUserClassModifier(t *testing.T) {
	d, err := New(WithUserClass([]byte("linuxboot"), false))
	require.NoError(t, err)

	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, "User Class Information -> linuxboot", d.Options[0].String())
	require.Equal(t, expected, d.Options[0].ToBytes())
}

func TestUserClassModifierRFC(t *testing.T) {
	d, err := New(WithUserClass([]byte("linuxboot"), true))
	require.NoError(t, err)

	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, "User Class Information -> linuxboot", d.Options[0].String())
	require.Equal(t, expected, d.Options[0].ToBytes())
}

func TestWithNetboot(t *testing.T) {
	d, err := New(WithNetboot)
	require.NoError(t, err)

	require.Equal(t, "Parameter Request List -> TFTP Server Name, Bootfile Name", d.Options[0].String())
}

func TestWithNetbootExistingTFTP(t *testing.T) {
	d, err := New()
	require.NoError(t, err)
	d.UpdateOption(&OptParameterRequestList{
		RequestedOpts: []OptionCode{OptionTFTPServerName},
	})
	WithNetboot(d)
	require.Equal(t, "Parameter Request List -> TFTP Server Name, Bootfile Name", d.Options[0].String())
}

func TestWithNetbootExistingBootfileName(t *testing.T) {
	d, _ := New()
	d.UpdateOption(&OptParameterRequestList{
		RequestedOpts: []OptionCode{OptionBootfileName},
	})
	WithNetboot(d)
	require.Equal(t, "Parameter Request List -> Bootfile Name, TFTP Server Name", d.Options[0].String())
}

func TestWithNetbootExistingBoth(t *testing.T) {
	d, _ := New()
	d.UpdateOption(&OptParameterRequestList{
		RequestedOpts: []OptionCode{OptionBootfileName, OptionTFTPServerName},
	})
	WithNetboot(d)
	require.Equal(t, "Parameter Request List -> Bootfile Name, TFTP Server Name", d.Options[0].String())
}

func TestWithRequestedOptions(t *testing.T) {
	// Check if OptionParameterRequestList is created when not present
	d, err := New(WithRequestedOptions(OptionFQDN))
	require.NoError(t, err)
	require.NotNil(t, d)
	o := d.GetOneOption(OptionParameterRequestList)
	require.NotNil(t, o)
	opts := o.(*OptParameterRequestList)
	require.ElementsMatch(t, opts.RequestedOpts, []OptionCode{OptionFQDN})

	// Check if already set options are preserved
	WithRequestedOptions(OptionHostName)(d)
	require.NotNil(t, d)
	o = d.GetOneOption(OptionParameterRequestList)
	require.NotNil(t, o)
	opts = o.(*OptParameterRequestList)
	require.ElementsMatch(t, opts.RequestedOpts, []OptionCode{OptionFQDN, OptionHostName})
}

func TestWithRelay(t *testing.T) {
	ip := net.IP{10, 0, 0, 1}
	d, err := New(WithRelay(ip))
	require.NoError(t, err)

	require.True(t, d.IsUnicast(), "expected unicast")
	require.Equal(t, ip, d.GatewayIPAddr)
	require.Equal(t, uint8(1), d.HopCount)
}

func TestWithNetmask(t *testing.T) {
	d, err := New(WithNetmask(net.IPv4Mask(255, 255, 255, 0)))
	require.NoError(t, err)

	require.Equal(t, 1, len(d.Options))
	require.Equal(t, OptionSubnetMask, d.Options[0].Code())
	osm := d.Options[0].(*OptSubnetMask)
	require.Equal(t, net.IPv4Mask(255, 255, 255, 0), osm.SubnetMask)
}

func TestWithLeaseTime(t *testing.T) {
	d, err := New(WithLeaseTime(uint32(3600)))
	require.NoError(t, err)

	require.Equal(t, 1, len(d.Options))
	require.Equal(t, OptionIPAddressLeaseTime, d.Options[0].Code())
	olt := d.Options[0].(*OptIPAddressLeaseTime)
	require.Equal(t, uint32(3600), olt.LeaseTime)
}

func TestWithDNS(t *testing.T) {
	d, err := New(WithDNS(net.ParseIP("10.0.0.1"), net.ParseIP("10.0.0.2")))
	require.NoError(t, err)

	require.Equal(t, 1, len(d.Options))
	require.Equal(t, OptionDomainNameServer, d.Options[0].Code())
	olt := d.Options[0].(*OptDomainNameServer)
	require.Equal(t, 2, len(olt.NameServers))
	require.Equal(t, net.ParseIP("10.0.0.1"), olt.NameServers[0])
	require.Equal(t, net.ParseIP("10.0.0.2"), olt.NameServers[1])
	require.NotEqual(t, net.ParseIP("10.0.0.1"), olt.NameServers[1])
}

func TestWithDomainSearchList(t *testing.T) {
	d, err := New(WithDomainSearchList("slackware.it", "dhcp.slackware.it"))
	require.NoError(t, err)

	require.Equal(t, 1, len(d.Options))
	osl := d.Options[0].(*OptDomainSearch)
	require.Equal(t, OptionDNSDomainSearchList, osl.Code())
	require.NotNil(t, osl.DomainSearch)
	require.Equal(t, 2, len(osl.DomainSearch.Labels))
	require.Equal(t, "slackware.it", osl.DomainSearch.Labels[0])
	require.Equal(t, "dhcp.slackware.it", osl.DomainSearch.Labels[1])
}

func TestWithRouter(t *testing.T) {
	rtr := net.ParseIP("10.0.0.254")
	d, err := New(WithRouter(rtr))
	require.NoError(t, err)

	require.Equal(t, 1, len(d.Options))
	ortr := d.Options[0].(*OptRouter)
	require.Equal(t, OptionRouter, ortr.Code())
	require.Equal(t, 1, len(ortr.Routers))
	require.Equal(t, rtr, ortr.Routers[0])
}
