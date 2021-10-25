package ztpv4

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseClassIdentifier(t *testing.T) {
	tt := []struct {
		name         string
		vc, hostname string
		ci           []byte // Client Identifier
		want         *VendorData
		fail         bool
	}{
		{name: "empty", fail: true},
		{name: "unknownVendor", vc: "VendorX;BFR10K;XX12345", fail: true},
		{name: "truncatedVendor", vc: "Arista;1234", fail: true},
		{
			name: "arista",
			vc:   "Arista;DCS-7050S-64;01.23;JPE12345678",
			want: &VendorData{VendorName: "Arista", Model: "DCS-7050S-64", Serial: "JPE12345678"},
		},
		{
			name: "juniper",
			vc:   "Juniper-ptx1000-DD123",
			want: &VendorData{VendorName: "Juniper", Model: "ptx1000", Serial: "DD123"},
		},
		{
			name: "juniperModelDash",
			vc:   "Juniper-qfx10002-36q-DN817",
			want: &VendorData{VendorName: "Juniper", Model: "qfx10002-36q", Serial: "DN817"},
		},
		{
			name:     "juniperHostnameSerial",
			vc:       "Juniper-qfx10008",
			hostname: "DE123",
			want:     &VendorData{VendorName: "Juniper", Model: "qfx10008", Serial: "DE123"},
		},
		{name: "juniperNoSerial", vc: "Juniper-qfx10008", fail: true},
		{
			name: "zpe",
			vc:   "ZPESystems:NSC:001234567",
			want: &VendorData{VendorName: "ZPESystems", Model: "NSC", Serial: "001234567"},
		},
		{
			name: "cisco",
			vc:   "FPR4100",
			ci:   []byte("JMX2525X0BW"),
			want: &VendorData{VendorName: "Cisco Systems", Model: "FPR4100", Serial: "JMX2525X0BW"},
		},
		{name: "ciscoNoSerial", vc: "FPR4100", fail: true},
		{
			name: "ciena",
			vc:   "1271-00011E00-032",
			ci:   []byte("JUSTASN"),
			want: &VendorData{VendorName: "Ciena Corporation", Model: "00011E00-032", Serial: "JUSTASN"},
		},
		{name: "cienaInvalidVendorClass", vc: "127100011E00032", fail: true},
		{name: "cienaNoSerial", vc: "1271-00011E00-032", fail: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			packet, err := dhcpv4.New()
			if err != nil {
				t.Fatalf("failed to creat dhcpv4 packet object: %v", err)
			}

			if tc.vc != "" {
				packet.UpdateOption(dhcpv4.OptClassIdentifier(tc.vc))
			}
			if tc.hostname != "" {
				packet.UpdateOption(dhcpv4.OptHostName(tc.hostname))
			}
			if tc.ci != nil {
				packet.UpdateOption(dhcpv4.OptClientIdentifier(tc.ci))
			}

			vd, err := ParseVendorData(packet)
			if tc.fail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, vd)
			}
		})
	}
}

func TestParseVIVC(t *testing.T) {
	tt := []struct {
		name  string
		vivc  string
		entID iana.EntID
		want  *VendorData
		fail  bool
	}{
		{
			name:  "cisco",
			entID: iana.EntIDCiscoSystems,
			vivc:  "SN:0;PID:R-IOSXRV9000-CC",
			want:  &VendorData{VendorName: "Cisco Systems", Model: "R-IOSXRV9000-CC", Serial: "0"},
		},
		{
			name:  "ciscoMultipleColonDelimiters",
			entID: iana.EntIDCiscoSystems,
			vivc:  "SN:0:123;PID:R-IOSXRV9000-CC:456",
			fail:  true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			packet, err := dhcpv4.New()
			if err != nil {
				t.Fatalf("failed to creat dhcpv4 packet object: %v", err)
			}

			if tc.vivc != "" {
				vivc := dhcpv4.VIVCIdentifier{EntID: uint32(tc.entID), Data: []byte(tc.vivc)}
				packet.UpdateOption(dhcpv4.OptVIVC(vivc))
			}

			vd, err := ParseVendorData(packet)
			if tc.fail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, vd)
			}
		})
	}
}
