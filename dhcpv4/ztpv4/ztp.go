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

func parseV4VIVC(vd *VendorData, packet *dhcpv4.DHCPv4) error {
	vivc := packet.VIVC()
	for _, id := range vivc {
		if id.EntID == iana.EntIDCiscoSystems {
			vd.VendorName = "Cisco Systems"
			//SN:0;PID:R-IOSXRV9000-CC
			for _, f := range bytes.Split(id.Data, []byte(";")) {
				p := bytes.Split(f, []byte(":"))
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

func parseV4VendorClass(vd *VendorData, packet *dhcpv4.DHCPv4) error {
	vc := packet.ClassIdentifier()
	if len(vc) == 0 {
		return errors.New("vendor options not found")
	}

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
				return errors.New("host name option is missing")
			}
		} else {
			vd.Model = strings.Join(p[1:len(p)-1], "-")
			vd.Serial = p[len(p)-1]
		}

		vd.VendorName = p[0]
		return nil
	}

	return nil
}

// ParseVendorData will try to parse dhcp4 options looking for more
// specific vendor data (like model, serial number, etc).
func ParseVendorData(packet *dhcpv4.DHCPv4) (*VendorData, error) {
	vd := &VendorData{}

	if err := parseV4VIVC(vd, packet); err != nil {
		return nil, err
	}

	if err := parseV4VendorClass(vd, packet); err != nil {
		return nil, err
	}

	// Check if VendorData did not get set
	if (VendorData{} == *vd) {
		// We didn't match anything.
		return nil, errors.New("no known ZTP vendor found")
	}

	return vd, nil
}
