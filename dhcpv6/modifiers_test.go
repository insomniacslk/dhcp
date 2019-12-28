package dhcpv6

import (
	"log"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestWithClientID(t *testing.T) {
	duid := Duid{
		Type:          DUID_LL,
		HwType:        iana.HWTypeEthernet,
		LinkLayerAddr: net.HardwareAddr([]byte{0xfa, 0xce, 0xb0, 0x00, 0x00, 0x0c}),
	}
	m, err := NewMessage(WithClientID(duid))
	require.NoError(t, err)
	cid := m.Options.ClientID()
	require.Equal(t, cid, &duid)
}

func TestWithServerID(t *testing.T) {
	duid := Duid{
		Type:          DUID_LL,
		HwType:        iana.HWTypeEthernet,
		LinkLayerAddr: net.HardwareAddr([]byte{0xfa, 0xce, 0xb0, 0x00, 0x00, 0x0c}),
	}
	m, err := NewMessage(WithServerID(duid))
	require.NoError(t, err)
	sid := m.Options.ServerID()
	require.Equal(t, sid, &duid)
}

func TestWithRequestedOptions(t *testing.T) {
	// Check if ORO is created when no ORO present
	m, err := NewMessage(WithRequestedOptions(OptionClientID))
	require.NoError(t, err)
	opt := m.GetOneOption(OptionORO)
	require.NotNil(t, opt)
	oro := opt.(*OptRequestedOption)
	require.ElementsMatch(t, oro.RequestedOptions(), []OptionCode{OptionClientID})
	// Check if already set options are preserved
	WithRequestedOptions(OptionServerID)(m)
	opt = m.GetOneOption(OptionORO)
	require.NotNil(t, opt)
	oro = opt.(*OptRequestedOption)
	require.ElementsMatch(t, oro.RequestedOptions(), []OptionCode{OptionClientID, OptionServerID})
}

func TestWithIANA(t *testing.T) {
	var d Message
	WithIANA(OptIAAddress{
		IPv6Addr:          net.ParseIP("::1"),
		PreferredLifetime: 3600,
		ValidLifetime:     5200,
	})(&d)
	require.Equal(t, 1, len(d.Options.Options))
	require.Equal(t, OptionIANA, d.Options.Options[0].Code())
}

func TestWithDNS(t *testing.T) {
	var d Message
	WithDNS([]net.IP{
		net.ParseIP("fe80::1"),
		net.ParseIP("fe80::2"),
	}...)(&d)
	require.Equal(t, 1, len(d.Options.Options))
	dns := d.Options.Options[0].(*OptDNSRecursiveNameServer)
	log.Printf("DNS %+v", dns)
	require.Equal(t, OptionDNSRecursiveNameServer, dns.Code())
	require.Equal(t, 2, len(dns.NameServers))
	require.Equal(t, net.ParseIP("fe80::1"), dns.NameServers[0])
	require.Equal(t, net.ParseIP("fe80::2"), dns.NameServers[1])
	require.NotEqual(t, net.ParseIP("fe80::1"), dns.NameServers[1])
}

func TestWithDomainSearchList(t *testing.T) {
	var d Message
	WithDomainSearchList([]string{
		"slackware.it",
		"dhcp.slackware.it",
	}...)(&d)
	require.Equal(t, 1, len(d.Options.Options))
	osl := d.Options.Options[0].(*OptDomainSearchList)
	require.Equal(t, OptionDomainSearchList, osl.Code())
	require.NotNil(t, osl.DomainSearchList)
	labels := osl.DomainSearchList.Labels
	require.Equal(t, 2, len(labels))
	require.Equal(t, "slackware.it", labels[0])
	require.Equal(t, "dhcp.slackware.it", labels[1])
}

func TestWithFQDN(t *testing.T) {
	var d Message
	WithFQDN(4, "cnos.localhost")(&d)
	require.Equal(t, 1, len(d.Options))
	ofqdn := d.Options[0].(*OptFQDN)
	require.Equal(t, OptionFQDN, ofqdn.Code())
	require.Equal(t, uint8(4), ofqdn.Flags)
	require.Equal(t, "cnos.localhost", ofqdn.DomainName)
}
