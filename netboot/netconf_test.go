package netboot

import (
	"log"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func getAdv(modifiers ...dhcpv6.Modifier) *dhcpv6.Message {
	hwaddr, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	if err != nil {
		log.Panic(err)
	}

	duid := dhcpv6.Duid{
		Type:          dhcpv6.DUID_LLT,
		HwType:        iana.HWTypeEthernet,
		Time:          dhcpv6.GetTime(),
		LinkLayerAddr: hwaddr,
	}
	sol, err := dhcpv6.NewSolicitWithCID(duid, modifiers...)
	if err != nil {
		log.Panic(err)
	}
	d, err := dhcpv6.NewAdvertiseFromSolicit(sol, modifiers...)
	if err != nil {
		log.Panic(err)
	}
	adv := d.(*dhcpv6.DHCPv6Message)
	return adv
}

func TestGetNetConfFromPacketv6Invalid(t *testing.T) {
	adv := getAdv()
	_, err := GetNetConfFromPacketv6(adv)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv6NoAddrsNoDNS(t *testing.T) {
	adv := getAdv(dhcpv6.WithIANA())
	_, err := GetNetConfFromPacketv6(adv)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv6NoDNS(t *testing.T) {
	addrs := []dhcpv6.OptIAAddress{
		dhcpv6.OptIAAddress{
			IPv6Addr:          net.ParseIP("::1"),
			PreferredLifetime: 3600,
			ValidLifetime:     5200,
		},
	}
	adv := getAdv(dhcpv6.WithIANA(addrs...))
	_, err := GetNetConfFromPacketv6(adv)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv6NoSearchList(t *testing.T) {
	addrs := []dhcpv6.OptIAAddress{
		dhcpv6.OptIAAddress{
			IPv6Addr:          net.ParseIP("::1"),
			PreferredLifetime: 3600,
			ValidLifetime:     5200,
		},
	}
	adv := getAdv(
		dhcpv6.WithIANA(addrs...),
		dhcpv6.WithDNS(net.ParseIP("fe80::1")),
	)
	_, err := GetNetConfFromPacketv6(adv)
	require.NoError(t, err)
}

func TestGetNetConfFromPacketv6(t *testing.T) {
	addrs := []dhcpv6.OptIAAddress{
		dhcpv6.OptIAAddress{
			IPv6Addr:          net.ParseIP("::1"),
			PreferredLifetime: 3600,
			ValidLifetime:     5200,
		},
	}
	adv := getAdv(
		dhcpv6.WithIANA(addrs...),
		dhcpv6.WithDNS(net.ParseIP("fe80::1")),
		dhcpv6.WithDomainSearchList("slackware.it"),
	)
	netconf, err := GetNetConfFromPacketv6(adv)
	require.NoError(t, err)
	// check addresses
	require.Equal(t, 1, len(netconf.Addresses))
	require.Equal(t, net.ParseIP("::1"), netconf.Addresses[0].IPNet.IP)
	require.Equal(t, 3600, netconf.Addresses[0].PreferredLifetime)
	require.Equal(t, 5200, netconf.Addresses[0].ValidLifetime)
	// check DNSes
	require.Equal(t, 1, len(netconf.DNSServers))
	require.Equal(t, net.ParseIP("fe80::1"), netconf.DNSServers[0])
	// check DNS search list
	require.Equal(t, 1, len(netconf.DNSSearchList))
	require.Equal(t, "slackware.it", netconf.DNSSearchList[0])
	// check routers
	require.Equal(t, 0, len(netconf.Routers))
}

func TestGetNetConfFromPacketv4AddrZero(t *testing.T) {
	d, _ := dhcpv4.New(dhcpv4.WithYourIP(net.IPv4zero))
	_, err := GetNetConfFromPacketv4(d)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv4NoMask(t *testing.T) {
	d, _ := dhcpv4.New(dhcpv4.WithYourIP(net.ParseIP("10.0.0.1")))
	_, err := GetNetConfFromPacketv4(d)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv4NullMask(t *testing.T) {
	d, _ := dhcpv4.New(
		dhcpv4.WithNetmask(net.IPv4Mask(0, 0, 0, 0)),
		dhcpv4.WithYourIP(net.ParseIP("10.0.0.1")),
	)
	_, err := GetNetConfFromPacketv4(d)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv4NoLeaseTime(t *testing.T) {
	d, _ := dhcpv4.New(
		dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
		dhcpv4.WithYourIP(net.ParseIP("10.0.0.1")),
	)
	_, err := GetNetConfFromPacketv4(d)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv4NoDNS(t *testing.T) {
	d, _ := dhcpv4.New(
		dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
		dhcpv4.WithLeaseTime(uint32(0)),
		dhcpv4.WithYourIP(net.ParseIP("10.0.0.1")),
	)
	_, err := GetNetConfFromPacketv4(d)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv4EmptyDNSList(t *testing.T) {
	d, _ := dhcpv4.New(
		dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
		dhcpv4.WithLeaseTime(uint32(0)),
		dhcpv4.WithDNS(),
		dhcpv4.WithYourIP(net.ParseIP("10.0.0.1")),
	)
	_, err := GetNetConfFromPacketv4(d)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv4NoSearchList(t *testing.T) {
	d, _ := dhcpv4.New(
		dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
		dhcpv4.WithLeaseTime(uint32(0)),
		dhcpv4.WithDNS(net.ParseIP("10.10.0.1"), net.ParseIP("10.10.0.2")),
		dhcpv4.WithYourIP(net.ParseIP("10.0.0.1")),
	)
	_, err := GetNetConfFromPacketv4(d)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv4EmptySearchList(t *testing.T) {
	d, _ := dhcpv4.New(
		dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
		dhcpv4.WithLeaseTime(uint32(0)),
		dhcpv4.WithDNS(net.ParseIP("10.10.0.1"), net.ParseIP("10.10.0.2")),
		dhcpv4.WithDomainSearchList(),
		dhcpv4.WithYourIP(net.ParseIP("10.0.0.1")),
	)
	_, err := GetNetConfFromPacketv4(d)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv4NoRouter(t *testing.T) {
	d, _ := dhcpv4.New(
		dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
		dhcpv4.WithLeaseTime(uint32(0)),
		dhcpv4.WithDNS(net.ParseIP("10.10.0.1"), net.ParseIP("10.10.0.2")),
		dhcpv4.WithDomainSearchList("slackware.it", "dhcp.slackware.it"),
		dhcpv4.WithYourIP(net.ParseIP("10.0.0.1")),
	)
	_, err := GetNetConfFromPacketv4(d)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv4EmptyRouter(t *testing.T) {
	d, _ := dhcpv4.New(
		dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
		dhcpv4.WithLeaseTime(uint32(0)),
		dhcpv4.WithDNS(net.ParseIP("10.10.0.1"), net.ParseIP("10.10.0.2")),
		dhcpv4.WithDomainSearchList("slackware.it", "dhcp.slackware.it"),
		dhcpv4.WithRouter(),
		dhcpv4.WithYourIP(net.ParseIP("10.0.0.1")),
	)
	_, err := GetNetConfFromPacketv4(d)
	require.Error(t, err)
}

func TestGetNetConfFromPacketv4(t *testing.T) {
	d, _ := dhcpv4.New(
		dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
		dhcpv4.WithLeaseTime(uint32(5200)),
		dhcpv4.WithDNS(net.ParseIP("10.10.0.1"), net.ParseIP("10.10.0.2")),
		dhcpv4.WithDomainSearchList("slackware.it", "dhcp.slackware.it"),
		dhcpv4.WithRouter(net.ParseIP("10.0.0.254")),
		dhcpv4.WithYourIP(net.ParseIP("10.0.0.1")),
	)

	netconf, err := GetNetConfFromPacketv4(d)
	require.NoError(t, err)
	// check addresses
	require.Equal(t, 1, len(netconf.Addresses))
	require.Equal(t, net.ParseIP("10.0.0.1"), netconf.Addresses[0].IPNet.IP)
	require.Equal(t, 0, netconf.Addresses[0].PreferredLifetime)
	require.Equal(t, 5200, netconf.Addresses[0].ValidLifetime)
	// check DNSes
	require.Equal(t, 2, len(netconf.DNSServers))
	require.Equal(t, net.ParseIP("10.10.0.1").To4(), netconf.DNSServers[0])
	require.Equal(t, net.ParseIP("10.10.0.2").To4(), netconf.DNSServers[1])
	// check DNS search list
	require.Equal(t, 2, len(netconf.DNSSearchList))
	require.Equal(t, "slackware.it", netconf.DNSSearchList[0])
	require.Equal(t, "dhcp.slackware.it", netconf.DNSSearchList[1])
	// check routers
	require.Equal(t, 1, len(netconf.Routers))
	require.Equal(t, net.ParseIP("10.0.0.254").To4(), netconf.Routers[0])
}
