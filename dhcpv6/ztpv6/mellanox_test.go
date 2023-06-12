package ztpv6

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseMellanoxVendorData(t *testing.T) {
	tt := []struct {
		name       string
		vendorOpts []dhcpv6.Option
		want       *VendorData
		fail       bool
	}{
		{name: "empty", fail: true},
		{
			name: "ok",
			fail: false,
			vendorOpts: []dhcpv6.Option{
				&dhcpv6.OptionGeneric{OptionData: []byte("SomeModel"), OptionCode: dhcpv6.OptionCode(MlnxSubOptionModel)},
				&dhcpv6.OptionGeneric{OptionData: []byte("SomeModel-1234"), OptionCode: dhcpv6.OptionCode(MlnxSubOptionPartNum)},
				&dhcpv6.OptionGeneric{OptionData: []byte("ABC1234"), OptionCode: dhcpv6.OptionCode(MlnxSubOptionSerial)},
				&dhcpv6.OptionGeneric{OptionData: []byte("1.2.3"), OptionCode: dhcpv6.OptionCode(MlnxSubOptionRelease)},
			},
			want: &VendorData{
				VendorName: iana.EnterpriseIDMellanoxTechnologiesLTD.String(),
				Model:      "SomeModel",
				Serial:     "ABC1234",
			},
		},
		{
			name: "no model",
			fail: true,
			vendorOpts: []dhcpv6.Option{
				&dhcpv6.OptionGeneric{OptionData: []byte("ABC1234"), OptionCode: dhcpv6.OptionCode(MlnxSubOptionSerial)},
			},
			want: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			packet, err := dhcpv6.NewMessage()
			if err != nil {
				t.Fatalf("failed to creat dhcpv6 packet object: %v", err)
			}

			packet.AddOption(&dhcpv6.OptVendorOpts{
				VendorOpts: tc.vendorOpts, EnterpriseNumber: uint32(iana.EnterpriseIDMellanoxTechnologiesLTD)})

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
