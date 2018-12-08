package netboot

import (
	"log"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func getAdv(modifiers ...dhcpv6.Modifier) *dhcpv6.DHCPv6Message {
	hwaddr, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	if err != nil {
		log.Panic(err)
	}

	duid := dhcpv6.Duid{
		Type:          dhcpv6.DUID_LLT,
		HwType:        iana.HwTypeEthernet,
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
	_, err := GetNetConfFromPacketv6(adv)
	require.NoError(t, err)
}
