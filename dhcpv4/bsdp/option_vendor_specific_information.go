package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// VendorOptions is like dhcpv4.Options, but stringifies using BSDP-specific
// option codes.
type VendorOptions struct {
	dhcpv4.Options
}

// String prints the contained options using BSDP-specific option code parsing.
func (v VendorOptions) String() string {
	return v.Options.ToString(bsdpHumanizer)
}

// FromBytes parses vendor options from
func (v *VendorOptions) FromBytes(data []byte) error {
	v.Options = make(dhcpv4.Options)
	return v.Options.FromBytes(data)
}

// OptVendorOptions returns the BSDP Vendor Specific Info in o.
func OptVendorOptions(o ...dhcpv4.Option) dhcpv4.Option {
	return dhcpv4.Option{
		Code:  dhcpv4.OptionVendorSpecificInformation,
		Value: VendorOptions{dhcpv4.OptionsFromList(o...)},
	}
}

// GetVendorOptions returns a new BSDP Vendor Specific Info option.
func GetVendorOptions(o dhcpv4.Options) *VendorOptions {
	v := o.Get(dhcpv4.OptionVendorSpecificInformation)
	if v == nil {
		return nil
	}
	var vo VendorOptions
	if err := vo.FromBytes(v); err != nil {
		return nil
	}
	return &vo
}

var bsdpHumanizer = dhcpv4.OptionHumanizer{
	ValueHumanizer: parseOption,
	CodeHumanizer: func(c uint8) dhcpv4.OptionCode {
		return optionCode(c)
	},
}

// parseOption is similar to dhcpv4.parseOption, except that it interprets
// option codes based on the BSDP-specific options.
func parseOption(code dhcpv4.OptionCode, data []byte) fmt.Stringer {
	var d dhcpv4.OptionDecoder
	switch code {
	case OptionMachineName:
		var s dhcpv4.String
		d = &s

	case OptionServerIdentifier:
		d = &dhcpv4.IP{}

	case OptionServerPriority, OptionReplyPort:
		var u dhcpv4.Uint16
		d = &u

	case OptionBootImageList:
		d = &BootImageList{}

	case OptionDefaultBootImageID, OptionSelectedBootImageID:
		d = &BootImageID{}

	case OptionMessageType:
		var m MessageType
		d = &m

	case OptionVersion:
		d = &Version{}
	}
	if d != nil && d.FromBytes(data) == nil {
		return d
	}
	return dhcpv4.OptionGeneric{data}
}
