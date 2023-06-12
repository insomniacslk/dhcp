package ztpv6

import (
	"errors"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/iana"
)

type MlnxSubOption uint16

const (
	MlnxSubOptionModel   MlnxSubOption = 1
	MlnxSubOptionPartNum MlnxSubOption = 2
	MlnxSubOptionSerial  MlnxSubOption = 3
	MlnxSubOptionMac     MlnxSubOption = 4
	MlnxSubOptionProfile MlnxSubOption = 5
	MlnxSubOptionRelease MlnxSubOption = 6
)

func getMellanoxVendorData(vendorOptsOption *dhcpv6.OptVendorOpts) (*VendorData, error) {
	vd := VendorData{}
	vd.VendorName = iana.EnterpriseIDMellanoxTechnologiesLTD.String()
	for _, opt := range vendorOptsOption.VendorOpts {
		switch MlnxSubOption(opt.Code()) {
		case MlnxSubOptionSerial:
			vd.Serial = string(opt.ToBytes())
		case MlnxSubOptionModel:
			vd.Model = string(opt.ToBytes())
		}
	}
	if (vd.Serial == "") || (vd.Model == "") {
		return nil, errors.New("couldn't parse Mellanox sub-option for serial or model")
	}

	return &vd, nil
}
