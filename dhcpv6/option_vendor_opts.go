package dhcpv6

/*
 This module defines the OptVendorOpts structure.
 https://tools.ietf.org/html/rfc3315#section-22.17

				Option 17
       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |      OPTION_VENDOR_OPTS       |           option-len          |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                       enterprise-number                       |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      .                                                               .
      .                   option-data (sub-options)         					.
      .                                                               .
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

				Sub-Option
	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|          opt-code             |             option-len        |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	.                                                               .
	.                          option-data                          .
	.                                                               .
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// OptVendorOpts represents a DHCPv6 Status Code option
type OptVendorOpts struct {
	EnterpriseNumber uint32
	VendorOpts       []Option
}

// Code returns the option code
func (op *OptVendorOpts) Code() OptionCode {
	return OptionVendorOpts
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptVendorOpts) ToBytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionVendorOpts))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint32(buf[4:8], uint32(op.EnterpriseNumber))
	for _, opt := range op.VendorOpts {
		buf = append(buf, opt.ToBytes()...)
	}
	return buf
}

// Length returns the option length
func (op *OptVendorOpts) Length() int {
	l := 4 // 4 bytes for Enterprise Number
	for _, opt := range op.VendorOpts {
		l += 4 + opt.Length() // 4 bytes for Code and Length from Vendor
	}
	return l
}

// String returns a string representation of the VendorOpts data
func (op *OptVendorOpts) String() string {
	return fmt.Sprintf("OptVendorOpts{enterprisenum=%v, vendorOpts=%v}",
		op.EnterpriseNumber, op.VendorOpts,
	)
}

// ParseOptVendorOpts builds an OptVendorOpts structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptVendorOpts(data []byte) (*OptVendorOpts, error) {
	opt := OptVendorOpts{}
	if len(data) < 4 {
		return nil, fmt.Errorf("Invalid vendor opts data length. Expected at least 4 bytes, got %v", len(data))
	}
	opt.EnterpriseNumber = binary.BigEndian.Uint32(data[:4])

	var err error
	opt.VendorOpts, err = VendorOptionsFromBytes(data[4:])
	if err != nil {
		return nil, err
	}
	return &opt, nil
}

// VendorOptionsFromBytes builds a slice of GenericOptions from a slice of bytes
func VendorOptionsFromBytes(data []byte) ([]Option, error) {
	// Parse a sequence of bytes until the end and build a list of options from
	// it. Returns an error if any invalid option or length is found.
	options := make([]Option, 0, 10)
	if len(data) == 0 {
		// no options, no party
		return options, nil
	}
	if len(data) < 4 {
		// cannot be shorter than option code (2 bytes) + length (2 bytes)
		return nil, fmt.Errorf("Invalid options: shorter than 4 bytes")
	}
	idx := 0
	for {
		if idx == len(data) {
			break
		}
		if idx > len(data) {
			// this should never happen
			return nil, fmt.Errorf("Error: reading past the end of options")
		}
		opt, err := VendParseOption(data[idx:])
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
		idx += opt.Length() + 4 // 4 bytes for type + length
	}
	return options, nil
}

// VendParseOption builds a GenericOption from a slice of bytes
// We cannot use the exisitng ParseOption function in options.go because the
// sub-options include codes specific to each vendor. There are overlaps in these
// codes with RFC standard codes.
func VendParseOption(dataStart []byte) (Option, error) {
	// Parse a sequence of bytes as a single DHCPv6 option.
	// Returns the option structure, or an error if any.
	opt := &OptionGeneric{}

	if len(dataStart) < 4 {
		return opt, fmt.Errorf("Invalid DHCPv6 option: less than 4 bytes")
	}
	code := OptionCode(binary.BigEndian.Uint16(dataStart[:2]))
	length := int(binary.BigEndian.Uint16(dataStart[2:4]))
	if len(dataStart) < length+4 {
		return opt, fmt.Errorf("Invalid option length for option %v. Declared %v, actual %v",
			code, length, len(dataStart)-4,
		)
	}

	data := dataStart[4 : 4+length]
	if len(data) < 2 {
		return opt, errors.New("ParseOptVendorOpts: short data: missing length field")
	}
	opt.OptionData = data
	opt.OptionCode = code

	if len(opt.OptionData) < 1 {
		return opt, errors.New("ParseOptVendorOpts: at least one vendor options data is required")
	}

	if length != opt.Length() {
		return opt, fmt.Errorf("Error: declared length is different from actual length for option %d: %d != %d",
			code, opt.Length(), length)
	}

	return opt, nil
}
