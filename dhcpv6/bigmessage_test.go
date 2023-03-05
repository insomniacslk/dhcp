package dhcpv6

import (
	"net"
	"testing"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/iana"
)

func makeBigMessage() (*Message, *RelayMessage) {
	m4, _ := dhcpv4.NewDiscovery(net.HardwareAddr{0x1, 0x2, 0xde, 0xad, 0xbe, 0xef})

	m, _ := NewSolicit(net.HardwareAddr{0x1, 0x2, 0xde, 0xad, 0xbe, 0xef}, WithRapidCommit)

	oneiana := m.Options.OneIANA()
	iaaddr := &OptIAAddress{IPv6Addr: net.ParseIP("fe80::1")}
	iaaddr.Options.Add(&OptStatusCode{StatusCode: iana.StatusSuccess, StatusMessage: "yes"})
	oneiana.Options.Add(iaaddr)

	oneiata := &OptIATA{}
	oneiata.Options.Add(iaaddr)

	fourrd := &Opt4RD{}
	fourrd.Add(&Opt4RDMapRule{
		Prefix4: net.IPNet{
			IP:   net.IP{123, 123, 0, 0},
			Mask: net.CIDRMask(16, 32),
		},
		Prefix6: net.IPNet{
			IP:   net.ParseIP("fc80::"),
			Mask: net.CIDRMask(64, 128),
		},
	})
	fourrd.Add(&Opt4RDNonMapRule{
		HubAndSpoke: true,
	})

	iapd := &OptIAPD{
		IaId: [4]byte{0x1, 0x2, 0x3, 0x4},
	}
	iaprefix := &OptIAPrefix{
		Prefix: &net.IPNet{
			IP:   net.ParseIP("fc80::"),
			Mask: net.CIDRMask(64, 128),
		},
	}
	iaprefix.Options.Add(&OptStatusCode{StatusCode: iana.StatusSuccess, StatusMessage: "yeah whatever"})
	iapd.Options.Add(iaprefix)

	vendorOpts := &OptVendorOpts{
		EnterpriseNumber: 123,
	}
	vendorOpts.VendorOpts.Add(&OptionGeneric{OptionCode: 400, OptionData: []byte("foobar")})

	adv, _ := NewReplyFromMessage(m,
		WithOption(OptClientArchType(iana.INTEL_X86PC, iana.EFI_X86_64)),
		WithOption(OptBootFileURL("http://foobar")),
		WithOption(OptBootFileParam("loglevel=10", "uroot.nohwrng")),
		WithOption(OptClientLinkLayerAddress(iana.HWTypeEthernet, net.HardwareAddr{0x1, 0x2, 0xbe, 0xef, 0xde, 0xad})),
		WithOption(fourrd),
		WithOption(&OptDHCPv4Msg{m4}),
		WithOption(&OptDHCP4oDHCP6Server{[]net.IP{net.ParseIP("fe81::1")}}),
		WithOption(OptDNS(net.ParseIP("fe82::1"))),
		WithOption(iapd),
		WithOption(OptInformationRefreshTime(1*time.Second)),
		WithOption(OptInterfaceID([]byte{0x1, 0x2})),
		WithOption(&OptNetworkInterfaceID{
			Typ:   NII_PXE_GEN_I,
			Major: 1,
		}),
		WithOption(OptRelayPort(1026)),
		WithOption(&OptRemoteID{EnterpriseNumber: 300, RemoteID: []byte{0xde, 0xad, 0xbe, 0xed}}),
		WithOption(OptRequestedOption(OptionBootfileURL, OptionBootfileParam)),
		WithOption(OptServerID(&DUIDLL{HWType: iana.HWTypeEthernet, LinkLayerAddr: net.HardwareAddr{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}})),
		WithOption(&OptUserClass{[][]byte{[]byte("foo"), []byte("bar")}}),
		WithOption(oneiana),
		WithOption(oneiata),
		WithOption(&OptVendorClass{EnterpriseNumber: 300, Data: [][]byte{[]byte("foo"), []byte("bar")}}),
		WithOption(vendorOpts),
	)

	relayfw := &RelayMessage{
		MessageType: MessageTypeRelayForward,
	}
	relayfw.Options.Add(OptRelayMessage(adv))
	relayfw.Options.Add(&OptRemoteID{
		EnterpriseNumber: 0x123,
		RemoteID:         []byte{0x1, 0x2},
	})
	return adv, relayfw
}

func TestPrint(t *testing.T) {
	m, r := makeBigMessage()
	t.Log(m.String())
	t.Log(m.Summary())

	t.Log(r.String())
	t.Log(r.Summary())
}

func BenchmarkToBytes(b *testing.B) {
	_, r := makeBigMessage()
	for i := 0; i < b.N; i++ {
		_ = r.ToBytes()
	}
}

func BenchmarkFromBytes(b *testing.B) {
	_, r := makeBigMessage()
	buf := r.ToBytes()
	for i := 0; i < b.N; i++ {
		_, _ = FromBytes(buf)
	}
}
