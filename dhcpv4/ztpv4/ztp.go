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

func getVIVC(vd *VendorData, packet *dhcpv4.DHCPv4) error {
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

func parseArista(vd *VendorData, vc string) error {
	// Arista;DCS-7050S-64;01.23;JPE12221671
	p := strings.Split(vc, ";")
	if len(p) < 4 {
		return errVendorOptionMalformed
	}
	vd.VendorName = p[0]
	vd.Model = p[1]
	vd.Serial = p[3]
	return nil
}

func parseZPE(vd *VendorData, vc string) error {
	// ZPESystems:NSC:002251623
	p := strings.Split(vc, ":")
	if len(p) < 3 {
		return errVendorOptionMalformed
	}
	vd.VendorName = p[0]
	vd.Model = p[1]
	vd.Serial = p[2]
	return nil
}

func parseJuniper(vd *VendorData, vc string, packet *dhcpv4.DHCPv4) error {
	// Juniper option 60 parsing is a bit more nuanced.  The following are all
	// "valid" identifying stings for Juniper:
	//    Juniper-ptx1000-DD576      <vendor>-<model>-<serial
	//    Juniper-qfx10008           <vendor>-<model> (serial in hostname option)
	//    Juniper-qfx10002-361-DN817 <vendor>-<model>-<serial> (model has a dash in it!)
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

// ParseVendorData will try to parse dhcp4 options looking for more
// specific vendor data (like model, serial number, etc).
func ParseVendorData(packet *dhcpv4.DHCPv4) (*VendorData, error) {
	vd := &VendorData{}

	if err := getVIVC(vd, packet); err != nil {
		return nil, err
	}

	// Check if VendorData got set
	if (VendorData{} != *vd) {
		return vd, nil
	}

	vc := packet.ClassIdentifier()
	if len(vc) == 0 {
		return nil, errors.New("vendor options not found")
	}
	vd := &VendorData{}

	vivc := packet.VIVC()
	for _, id := range vivc {
		if id.EntID == entIDCiscoSystems {
			vd.VendorName = "Cisco Systems"
			//SN:0;PID:R-IOSXRV9000-CC
			for _, f := range bytes.Split(id.Data, []byte(";")) {
				p := bytes.SplitN(f, []byte(":"), 2)
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

	switch {
	case strings.HasPrefix(vc, "Arista;"):
		if err := parseArista(vd, vc); err != nil {
			return nil, err
		}
		return vd, nil

	case strings.HasPrefix(vc, "ZPESystems:"):
		if err := parseZPE(vd, vc); err != nil {
			return nil, err
		}
		return vd, nil

	case strings.HasPrefix(vc, "Juniper-"):
		if err := parseJuniper(vd, vc, packet); err != nil {
			return nil, err
		}
		return vd, nil
	}

	// We didn't match anything.
	return nil, errors.New("no known ZTP vendor found")
}

