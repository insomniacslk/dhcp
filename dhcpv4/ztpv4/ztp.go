package ztpv4

import (
	"bytes"
	"errors"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/iana"
)

// VendorData is optional data a particular vendor may or may not include
// in the Vendor Class options.
type VendorData struct {
	VendorName, Model, Serial string
}

var errVendorOptionMalformed = errors.New("malformed vendor option")

func parseClassIdentifier(packet *dhcpv4.DHCPv4) (*VendorData, error) {
	vd := &VendorData{}

	switch vc := packet.ClassIdentifier(); {
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
		p := strings.Split(vc, "-")
		if len(p) < 3 {
			vd.Model = p[1]
			vd.Serial = packet.HostName()
			if len(vd.Serial) == 0 {
				return nil, errors.New("host name option is missing")
			}
		} else {
			vd.Model = strings.Join(p[1:len(p)-1], "-")
			vd.Serial = p[len(p)-1]
		}

		vd.VendorName = p[0]
		return vd, nil

	// Cisco Firepower FPR4100/9300 models use Opt 60 for model info
	// and Opt 61 contains the serial number
	case vc == "FPR4100" || vc == "FPR9300":
		vd.VendorName = iana.EntIDCiscoSystems.String()
		vd.Model = vc
		vd.Serial = dhcpv4.GetString(dhcpv4.OptionClientIdentifier, packet.Options)
		return vd, nil

	}

	return nil, nil
}

func parseVIVC(packet *dhcpv4.DHCPv4) (*VendorData, error) {
	vd := &VendorData{}

	for _, id := range packet.VIVC() {
		if id.EntID == uint32(iana.EntIDCiscoSystems) {
			vd.VendorName = iana.EntIDCiscoSystems.String()
			//SN:0;PID:R-IOSXRV9000-CC
			for _, f := range bytes.Split(id.Data, []byte(";")) {
				p := bytes.Split(f, []byte(":"))
				if len(p) != 2 {
					return nil, errVendorOptionMalformed
				}

				switch string(p[0]) {
				case "SN":
					vd.Serial = string(p[1])
				case "PID":
					vd.Model = string(p[1])
				}
			}
			return vd, nil
		}
	}

	return nil, nil
}

// ParseVendorData will try to parse dhcp4 options looking for more
// specific vendor data (like model, serial number, etc).
func ParseVendorData(packet *dhcpv4.DHCPv4) (*VendorData, error) {
	vd, err := parseClassIdentifier(packet)
	if err != nil {
		return nil, err
	}

	// If VendorData is set, return early
	if vd != nil {
		return vd, nil
	}

	vd, err = parseVIVC(packet)
	if err != nil {
		return nil, err
	}

	if vd != nil {
		return vd, nil
	}

	return nil, errors.New("no known ZTP vendor found")
}
