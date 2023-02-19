package ztpv6

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseVendorDataWithVendorOpts(t *testing.T) {
	tt := []struct {
		name string
		vc   string
		want *VendorData
		fail bool
	}{
		{name: "empty", fail: true},
		{name: "unknownVendor", vc: "VendorX;BFR10K;XX12345", fail: true, want: nil},
		{name: "truncatedArista", vc: "Arista;1234", fail: true, want: nil},
		{name: "truncatedCisco", vc: "Cisco;1234", fail: true, want: nil},
		{name: "truncatedZPE", vc: "ZPESystems:1234", fail: true, want: nil},
		{
			name: "arista",
			vc:   "Arista;DCS-7050S-64;01.23;JPE12345678",
			want: &VendorData{VendorName: "Arista", Model: "DCS-7050S-64", Serial: "JPE12345678"},
		}, {
			name: "cisco",
			vc:   "Cisco;SYS-8801;01.23;FOC12345678",
			want: &VendorData{VendorName: "Cisco", Model: "SYS-8801", Serial: "FOC12345678"},
		}, {
			name: "zpe",
			vc:   "ZPESystems:NSC:001234567",
			want: &VendorData{VendorName: "ZPESystems", Model: "NSC", Serial: "001234567"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			packet, err := dhcpv6.NewMessage()
			if err != nil {
				t.Fatalf("failed to creat dhcpv6 packet object: %v", err)
			}

			opts := []dhcpv6.Option{&dhcpv6.OptionGeneric{OptionData: []byte(tc.vc), OptionCode: 1}}
			packet.AddOption(&dhcpv6.OptVendorOpts{
				VendorOpts: opts, EnterpriseNumber: 0000})

			vd, err := ParseVendorData(packet)
			if err != nil && !tc.fail {
				t.Errorf("unexpected failure: %v", err)
			}

			if vd != nil {
				require.Equal(t, *tc.want, *vd, "comparing vendor option data")
			} else {
				require.Equal(t, tc.want, vd, "comparing vendor option data")
			}
		})
	}
}

func TestParseVendorDataWithVendorClass(t *testing.T) {
	tt := []struct {
		name     string
		vc       string
		clientId dhcpv6.DUID
		want     *VendorData
		fail     bool
	}{
		{name: "empty", fail: true},
		{name: "unknownVendor", vc: "VendorX;BFR10K;XX12345", fail: true, want: nil},
		{name: "truncatedArista", vc: "Arista;1234", fail: true, want: nil},
		{name: "truncatedZPE", vc: "ZPESystems:1234", fail: true, want: nil},
		{
			name: "arista",
			vc:   "Arista;DCS-7050S-64;01.23;JPE12345678",
			want: &VendorData{VendorName: "Arista", Model: "DCS-7050S-64", Serial: "JPE12345678"},
		}, {
			name: "zpe",
			vc:   "ZPESystems:NSC:001234567",
			want: &VendorData{VendorName: "ZPESystems", Model: "NSC", Serial: "001234567"},
		},
		{
			name: "Ciena",
			vc:   "1271-23422Z11-123",
			clientId: &dhcpv6.DUIDEN{
				EnterpriseIdentifier: []byte("001234567"),
			},
			want: &VendorData{VendorName: iana.EnterpriseIDCienaCorporation.String(), Model: "23422Z11-123", Serial: "001234567"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			packet, err := dhcpv6.NewMessage()
			if err != nil {
				t.Fatalf("failed to creat dhcpv6 packet object: %v", err)
			}

			packet.AddOption(&dhcpv6.OptVendorClass{
				EnterpriseNumber: 0000, Data: [][]byte{[]byte(tc.vc)}})
			if tc.clientId != nil {
				packet.AddOption(dhcpv6.OptClientID(tc.clientId))
			}
			vd, err := ParseVendorData(packet)
			if err != nil && !tc.fail {
				t.Errorf("unexpected failure: %v", err)
			}

			if vd != nil {
				require.Equal(t, *tc.want, *vd, "comparing vendor class data")
			} else {
				require.Equal(t, tc.want, vd, "comparing vendor class data")
			}
		})
	}
}
