package ztpv6

import (
	"errors"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv6"
)

var (
	errVendorOptionMalformed = errors.New("malformed vendor option")
)

// VendorData contains fields extracted from Option 17 data
type VendorData struct {
	VendorName, Model, Serial string
}

// ParseVendorData will try to parse dhcp6 Vendor Specific Information options data
// looking for more specific vendor data (like model, serial number, etc).
// If the options are missing we will just return nil
func ParseVendorData(packet dhcpv6.DHCPv6) (*VendorData, error) {
	opt := packet.GetOneOption(dhcpv6.OptionVendorOpts)
	if opt == nil {
		return nil, errors.New("vendor options not found")
	}

	vd := VendorData{}
	vo := opt.(*dhcpv6.OptVendorOpts).VendorOpts

	for _, opt := range vo {
		optData := string(opt.(*dhcpv6.OptionGeneric).OptionData)
		switch {
		// Arista;DCS-0000;00.00;ZZZ00000000
		case strings.HasPrefix(optData, "Arista;"):
			p := strings.Split(optData, ";")
			if len(p) < 4 {
				return nil, errVendorOptionMalformed
			}

			vd.VendorName = p[0]
			vd.Model = p[1]
			vd.Serial = p[3]
			return &vd, nil

		// ZPESystems:NSC:000000000
		case strings.HasPrefix(optData, "ZPESystems:"):
			p := strings.Split(optData, ":")
			if len(p) < 3 {
				return nil, errVendorOptionMalformed
			}

			vd.VendorName = p[0]
			vd.Model = p[1]
			vd.Serial = p[2]
			return &vd, nil
		}
	}
	return nil, errors.New("failed to parse vendor option data")
}
