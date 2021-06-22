package main

import (
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/iana"
)

func main() {
	m4, _ := dhcpv4.NewDiscovery(net.HardwareAddr{0x1, 0x2, 0xde, 0xad, 0xbe, 0xef})

	m, _ := dhcpv6.NewSolicit(net.HardwareAddr{0x1, 0x2, 0xde, 0xad, 0xbe, 0xef}, dhcpv6.WithRapidCommit)
	fmt.Println(m.Summary())

	oneiana := m.Options.OneIANA()
	iaaddr := &dhcpv6.OptIAAddress{IPv6Addr: net.ParseIP("fe80::1")}
	iaaddr.Options.Add(&dhcpv6.OptStatusCode{StatusCode: iana.StatusSuccess, StatusMessage: "yes"})
	oneiana.Options.Add(iaaddr)

	oneiata := &dhcpv6.OptIATA{}
	oneiata.Options.Add(iaaddr)

	fourrd := &dhcpv6.Opt4RD{}
	fourrd.Options.Add(&dhcpv6.Opt4RDMapRule{
		Prefix4: net.IPNet{
			IP:   net.IP{123, 123, 0, 0},
			Mask: net.CIDRMask(16, 32),
		},
		Prefix6: net.IPNet{
			IP:   net.ParseIP("fc80::"),
			Mask: net.CIDRMask(64, 128),
		},
	})
	fourrd.Options.Add(&dhcpv6.Opt4RDNonMapRule{
		HubAndSpoke: true,
	})

	iapd := &dhcpv6.OptIAPD{
		IaId: [4]byte{0x1, 0x2, 0x3, 0x4},
	}
	iaprefix := &dhcpv6.OptIAPrefix{
		Prefix: &net.IPNet{
			IP:   net.ParseIP("fc80::"),
			Mask: net.CIDRMask(64, 128),
		},
	}
	iaprefix.Options.Add(&dhcpv6.OptStatusCode{StatusCode: iana.StatusSuccess, StatusMessage: "yeah whatever"})
	iapd.Options.Add(iaprefix)

	vendorOpts := &dhcpv6.OptVendorOpts{
		EnterpriseNumber: 123,
	}
	vendorOpts.VendorOpts.Add(&dhcpv6.OptionGeneric{OptionCode: 400, OptionData: []byte("foobar")})

	adv, _ := dhcpv6.NewReplyFromMessage(m,
		dhcpv6.WithOption(dhcpv6.OptBootFileURL("http://foobar")),
		dhcpv6.WithOption(dhcpv6.OptBootFileParam("loglevel=10", "uroot.nohwrng")),
		dhcpv6.WithOption(dhcpv6.OptClientLinkLayerAddress(iana.HWTypeEthernet, net.HardwareAddr{0x1, 0x2, 0xbe, 0xef, 0xde, 0xad})),
		dhcpv6.WithOption(fourrd),
		dhcpv6.WithOption(&dhcpv6.OptDHCPv4Msg{m4}),
		dhcpv6.WithOption(&dhcpv6.OptDHCP4oDHCP6Server{[]net.IP{net.ParseIP("fe81::1")}}),
		dhcpv6.WithOption(dhcpv6.OptDNS(net.ParseIP("fe82::1"))),
		dhcpv6.WithOption(iapd),
		dhcpv6.WithOption(dhcpv6.OptInformationRefreshTime(1*time.Second)),
		dhcpv6.WithOption(dhcpv6.OptInterfaceID([]byte{0x1, 0x2})),
		dhcpv6.WithOption(&dhcpv6.OptNetworkInterfaceID{
			Typ:   dhcpv6.NII_PXE_GEN_I,
			Major: 1,
		}),
		dhcpv6.WithOption(vendorOpts),
		dhcpv6.WithOption(oneiana),
		dhcpv6.WithOption(oneiata),
	)
	adv.Options.Add(&dhcpv6.OptIATA{})

	relayfw := dhcpv6.RelayMessage{
		MessageType: dhcpv6.MessageTypeRelayForward,
	}
	relayfw.Options.Add(dhcpv6.OptRelayMessage(adv))
	relayfw.Options.Add(&dhcpv6.OptRemoteID{
		EnterpriseNumber: 0x123,
		RemoteID:         []byte{0x1, 0x2},
	})
	fmt.Println(relayfw.Summary())
}
