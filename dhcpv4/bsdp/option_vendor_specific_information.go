package bsdp

import (
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// OptVendorSpecificInformation encapsulates the BSDP-specific options used for
// the protocol.
type OptVendorSpecificInformation struct {
	Options dhcpv4.Options
}

// parseOption is similar to dhcpv4.ParseOption, except that it switches based
// on the BSDP specific options.
func parseOption(code dhcpv4.OptionCode, data []byte) (dhcpv4.Option, error) {
	var (
		opt dhcpv4.Option
		err error
	)
	switch code {
	case OptionBootImageList:
		opt, err = ParseOptBootImageList(data)
	case OptionDefaultBootImageID:
		opt, err = ParseOptDefaultBootImageID(data)
	case OptionMachineName:
		opt, err = ParseOptMachineName(data)
	case OptionMessageType:
		opt, err = ParseOptMessageType(data)
	case OptionReplyPort:
		opt, err = ParseOptReplyPort(data)
	case OptionSelectedBootImageID:
		opt, err = ParseOptSelectedBootImageID(data)
	case OptionServerIdentifier:
		opt, err = ParseOptServerIdentifier(data)
	case OptionServerPriority:
		opt, err = ParseOptServerPriority(data)
	case OptionVersion:
		opt, err = ParseOptVersion(data)
	default:
		opt, err = ParseOptGeneric(code, data)
	}
	if err != nil {
		return nil, err
	}
	return opt, nil
}

// codeGetter is a dhcpv4.OptionCodeGetter for BSDP optionCodes.
func codeGetter(c uint8) dhcpv4.OptionCode {
	return optionCode(c)
}

// ParseOptVendorSpecificInformation constructs an OptVendorSpecificInformation struct from a sequence of
// bytes and returns it, or an error.
func ParseOptVendorSpecificInformation(data []byte) (*OptVendorSpecificInformation, error) {
	options, err := dhcpv4.OptionsFromBytesWithParser(data, codeGetter, parseOption, false /* don't check for OptionEnd tag */)
	if err != nil {
		return nil, err
	}
	return &OptVendorSpecificInformation{options}, nil
}

// Code returns the option code.
func (o *OptVendorSpecificInformation) Code() dhcpv4.OptionCode {
	return dhcpv4.OptionVendorSpecificInformation
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptVendorSpecificInformation) ToBytes() []byte {
	return uio.ToBigEndian(o.Options)
}

// String returns a human-readable string for this option.
func (o *OptVendorSpecificInformation) String() string {
	s := "Vendor Specific Information ->"
	for _, opt := range o.Options {
		optString := opt.String()
		// If this option has sub-structures, offset them accordingly.
		if strings.Contains(optString, "\n") {
			optString = strings.Replace(optString, "\n  ", "\n    ", -1)
		}
		s += "\n  " + optString
	}
	return s
}

// GetOneOption returns the first suboption that matches the OptionCode code.
func (o *OptVendorSpecificInformation) GetOneOption(code dhcpv4.OptionCode) dhcpv4.Option {
	return o.Options.GetOne(code)
}
