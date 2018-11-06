package ztpv4

import (
	"bytes"
	"errors"
	"strings"

	"github.com/golang/glog"
	"github.com/insomniacslk/dhcp/dhcpv4"
)

// VendorData is optional data a particular vendor may or may not include
// in the Vendor Class options.  All values are optional and will be zero
// values if not found.
type VendorData struct {
	VendorName string
	Model      string
	Serial     string
}

var errVendorOptionMalformed = errors.New("malformed vendor option")

// VendorDataV4 will try to parse dhcp4 options data looking for more specific
// vendor data (like model, serial number, etc).  If the options are missing
func VendorDataV4(packet *dhcpv4.DHCPv4) VendorData {
	vd := VendorData{}

	if err := parseV4VendorClass(&vd, packet); err != nil {
		glog.Errorf("failed to parse vendor data from vendor class: %v", err)
	}

	if err := parseV4VIVC(&vd, packet); err != nil {
		glog.Errorf("failed to parse vendor data from vendor-idenitifying vendor class: %v", err)
	}

	return vd
}

// parseV4Opt60 will attempt to look at the Vendor Class option (Option 60) on
// DHCPv4.  The option is formatted as a string with the content being specific
// for the vendor, usually using a deliminator to separate the values.
// See: https://tools.ietf.org/html/rfc1533#section-9.11
func parseV4VendorClass(vd *VendorData, packet *dhcpv4.DHCPv4) error {
	opt := packet.GetOneOption(dhcpv4.OptionClassIdentifier)
	if opt == nil {
		return nil
	}
	vc := opt.(*dhcpv4.OptClassIdentifier).Identifier

	switch {
	// Arista;DCS-7050S-64;01.23;JPE12221671
	case strings.HasPrefix(vc, "Arista;"):
		p := strings.Split(vc, ";")
		if len(p) < 4 {
			return errVendorOptionMalformed
		}

		vd.VendorName = p[0]
		vd.Model = p[1]
		vd.Serial = p[3]
		return nil

	// ZPESystems:NSC:002251623
	case strings.HasPrefix(vc, "ZPESystems:"):
		p := strings.Split(vc, ":")
		if len(p) < 3 {
			return errVendorOptionMalformed
		}

		vd.VendorName = p[0]
		vd.Model = p[1]
		vd.Serial = p[2]
		return nil

	// Juniper option 60 parsing is a bit more nuanced.  The following are all
	// "valid" indetifing stings for Juniper:
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
			}
		} else {
			vd.Serial = vc[sepIdx+1:]
			vc = vc[:sepIdx]
		}
		vd.Model = vc

		return nil
	}

	// We didn't match anything, just return an empty vendor data.
	return nil
}

const entIDCiscoSystems = 0x9

// parseV4Opt124 will attempt to read the Vendor-Identifying Vendor Class
// (Option 124) on a DHCPv4 packet.  The data is represented in a length/value
// format with an indentifying enterprise number.
//
// See: https://tools.ietf.org/html/rfc3925
func parseV4VIVC(vd *VendorData, packet *dhcpv4.DHCPv4) error {
	opt := packet.GetOneOption(dhcpv4.OptionVendorIdentifyingVendorClass)
	if opt == nil {
		return nil
	}
	ids := opt.(*dhcpv4.OptVIVC).Identifiers

	for _, id := range ids {
		if id.EntID == entIDCiscoSystems {
			vd.VendorName = "Cisco Systems"

			//SN:0;PID:R-IOSXRV9000-CC
			for _, f := range bytes.Split(id.Data, []byte(";")) {
				p := bytes.SplitN(f, []byte(":"), 2)
				if len(p) != 2 {
					return errVendorOptionMalformed
				}

				switch string(p[0]) {
				case "SN":
					vd.Serial = string(p[1])
				case "PID":
					vd.Model = string(p[1])
				}
			}
		}
	}
	return nil
}
