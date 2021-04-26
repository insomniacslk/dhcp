package ztpv4

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseV4VendorClass(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		hostname string
		entID    iana.EntID
		want     *VendorData
		fail     bool
	}{
		{name: "empty", fail: true},
		{name: "unknownVendor", input: "VendorX;BFR10K;XX12345", fail: true},
		{name: "truncatedVendor", input: "Arista;1234", fail: true},
		{
			name:  "arista",
			input: "Arista;DCS-7050S-64;01.23;JPE12345678",
			want:  &VendorData{VendorName: "Arista", Model: "DCS-7050S-64", Serial: "JPE12345678"},
		},
		{
			name:  "juniper",
			input: "Juniper-ptx1000-DD123",
			want:  &VendorData{VendorName: "Juniper", Model: "ptx1000", Serial: "DD123"},
		},
		{
			name:  "juniperModelDash",
			input: "Juniper-qfx10002-36q-DN817",
			want:  &VendorData{VendorName: "Juniper", Model: "qfx10002-36q", Serial: "DN817"},
		},
		{
			name:     "juniperHostnameSerial",
			input:    "Juniper-qfx10008",
			hostname: "DE123",
			want:     &VendorData{VendorName: "Juniper", Model: "qfx10008", Serial: "DE123"},
		},
		{name: "juniperNoSerial", input: "Juniper-qfx10008", fail: true},
		{
			name:  "zpe",
			input: "ZPESystems:NSC:001234567",
			want:  &VendorData{VendorName: "ZPESystems", Model: "NSC", Serial: "001234567"},
		},
		{
			name:  "cisco",
			entID: iana.EntIDCiscoSystems,
			input: "SN:0;PID:R-IOSXRV9000-CC",
			want:  &VendorData{VendorName: "Cisco Systems", Model: "R-IOSXRV9000-CC", Serial: "0"},
		},
		{
			name:  "ciscoMultipleColonDelimiters",
			entID: iana.EntIDCiscoSystems,
			input: "SN:0:123;PID:R-IOSXRV9000-CC:456",
			fail:  true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			packet, err := dhcpv4.New()
			if err != nil {
				t.Fatalf("failed to creat dhcpv4 packet object: %v", err)
			}

			if tc.input != "" {
				packet.UpdateOption(dhcpv4.OptClassIdentifier(tc.input))
			}
			if tc.hostname != "" {
				packet.UpdateOption(dhcpv4.OptHostName(tc.hostname))
			}
			if tc.entID == iana.EntIDCiscoSystems {
				vivc := dhcpv4.VIVCIdentifier{EntID: uint32(tc.entID), Data: []byte(tc.input)}
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
