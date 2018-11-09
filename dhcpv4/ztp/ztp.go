package ztpv4

import (
	"errors"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// VendorData is optional data a particular vendor may or may not include
// in the Vendor Class options.
type VendorData struct {
	VendorName string
	Model      string
	Serial     string
}

var errVendorOptionMalformed = errors.New("malformed vendor option")

// ParseVendorData will try to parse dhcp4 options looking for more
// specific vendor data (like model, serial number, etc).
func ParseVendorData(packet *dhcpv4.DHCPv4) (*VendorData, error) {
	return parseV4VendorClass(packet);
}

// parseV4Opt60 will attempt to look at the Vendor Class option (Option 60) on
// DHCPv4.  The option is formatted as a string with the content being specific
// for the vendor, usually using a delimitator to separate the values.
// See: https://tools.ietf.org/html/rfc1533#section-9.11
func parseV4VendorClass(packet *dhcpv4.DHCPv4) (*VendorData, error) {
	opt := packet.GetOneOption(dhcpv4.OptionClassIdentifier)
	if opt == nil {
		return nil, nil
	}
	vc := opt.(*dhcpv4.OptClassIdentifier).Identifier
	vd := &VendorData{}

	switch {
	// Arista;DCS-7050S-64;01.23;JPE12221671
	case strings.HasPrefix(vc, "Arista;"):
		p := strings.Split(vc, ";")
		if len(p) < 4 {
			return nil, errVendorOptionMalformed
		}

		vd.VendorName = p[0]
		vd.Model = p[1]
		vd.Serial = p[3]
		return vd, nil

	// ZPESystems:NSC:002251623
	case strings.HasPrefix(vc, "ZPESystems:"):
		p := strings.Split(vc, ":")
		if len(p) < 3 {
			return nil, errVendorOptionMalformed
		}

		vd.VendorName = p[0]
		vd.Model = p[1]
		vd.Serial = p[2]
		return vd, nil

	// Juniper option 60 parsing is a bit more nuanced.  The following are all
	// "valid" identifying stings for Juniper:
	//    Juniper-ptx1000-DD576      <vendor>-<model>-<serial
	//    Juniper-qfx10008           <vendor>-<model> (serial in hostname option)
	//    Juniper-qfx10002-361-DN817 <vendor>-<model>-<serial> (model has a dash in it!)
	case strings.HasPrefix(vc, "Juniper-"):
		// strip of the prefix
		vc := vc[len("Juniper-"):]
		vd.VendorName = "Juniper"

		sepIdx := strings.LastIndex(vc, "-")
		if sepIdx == -1 {
			// No separator was found. Attempt serial number from the hostname
			if opt := packet.GetOneOption(dhcpv4.OptionHostName); opt != nil {
				vd.Serial = opt.(*dhcpv4.OptHostName).HostName
			} else {
				return nil, errVendorOptionMalformed
			}
		} else {
			if len(vc) == sepIdx+1 {
				return nil, errVendorOptionMalformed
			}
			vd.Serial = vc[sepIdx+1:]
			vc = vc[:sepIdx]
		}
		vd.Model = vc

		return vd, nil
	}

	// We didn't match anything.
	return nil, nil
}
