// +build darwin

package bsdp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// OptVendorSpecificInformation encapsulates the BSDP-specific options used for
// the protocol.
type OptVendorSpecificInformation struct {
	Options []dhcpv4.Option
}

// parseOption is similar to dhcpv4.ParseOption, except that it switches based
// on the BSDP specific options.
func parseOption(data []byte) (dhcpv4.Option, error) {
	if len(data) == 0 {
		return nil, dhcpv4.ErrZeroLengthByteStream
	}
	var (
		opt dhcpv4.Option
		err error
	)
	switch dhcpv4.OptionCode(data[0]) {
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
		opt, err = ParseOptGeneric(data)
	}
	if err != nil {
		return nil, err
	}
	return opt, nil
}

// ParseOptVendorSpecificInformation constructs an OptVendorSpecificInformation struct from a sequence of
// bytes and returns it, or an error.
func ParseOptVendorSpecificInformation(data []byte) (*OptVendorSpecificInformation, error) {
	// Should at least have code + length
	if len(data) < 2 {
		return nil, dhcpv4.ErrShortByteStream
	}
	code := dhcpv4.OptionCode(data[0])
	if code != dhcpv4.OptionVendorSpecificInformation {
		return nil, fmt.Errorf("expected option %v, got %v instead", dhcpv4.OptionVendorSpecificInformation, code)
	}
	length := int(data[1])
	if len(data) < length+2 {
		return nil, fmt.Errorf("expected length 2, got %d instead", length)
	}

	options := make([]dhcpv4.Option, 0, 10)
	idx := 2
	for {
		if idx == len(data) {
			break
		}
		// This should never happen.
		if idx > len(data) {
			return nil, errors.New("read past the end of options")
		}
		opt, err := parseOption(data[idx:])
		if err != nil {
			return nil, err
		}
		options = append(options, opt)

		// Account for code + length bytes
		idx += 2 + opt.Length()
	}

	return &OptVendorSpecificInformation{options}, nil
}

// Code returns the option code.
func (o *OptVendorSpecificInformation) Code() dhcpv4.OptionCode {
	return dhcpv4.OptionVendorSpecificInformation
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptVendorSpecificInformation) ToBytes() []byte {
	bs := []byte{byte(o.Code()), byte(o.Length())}

	// Append data section
	for _, opt := range o.Options {
		bs = append(bs, opt.ToBytes()...)
	}
	return bs
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

// Length returns the length of the data portion of this option. Take into
// account code + data length bytes for each sub option.
func (o *OptVendorSpecificInformation) Length() int {
	var length int
	for _, opt := range o.Options {
		length += 2 + opt.Length()
	}
	return length
}

// GetOptions returns all suboptions that match the given OptionCode code.
func (o *OptVendorSpecificInformation) GetOptions(code dhcpv4.OptionCode) []dhcpv4.Option {
	var opts []dhcpv4.Option
	for _, opt := range o.Options {
		if opt.Code() == code {
			opts = append(opts, opt)
		}
	}
	return opts
}

// GetOption returns the first suboption that matches the OptionCode code.
func (o *OptVendorSpecificInformation) GetOption(code dhcpv4.OptionCode) dhcpv4.Option {
	opts := o.GetOptions(code)
	if len(opts) == 0 {
		return nil
	}
	return opts[0]
}
