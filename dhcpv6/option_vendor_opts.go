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
      .                   option-data (sub-options)                   .
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
	opt.VendorOpts, err = OptionsFromBytesWithParser(data[4:], vendParseOption)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}

// vendParseOption builds a GenericOption from a slice of bytes
// We cannot use the existing ParseOption function in options.go because the
// sub-options include codes specific to each vendor. There are overlaps in these
// codes with RFC standard codes.
func vendParseOption(dataStart []byte) (Option, error) {
	// Parse a sequence of bytes as a single DHCPv6 option.
	// Returns the option structure, or an error if any.

	if len(dataStart) < 4 {
		return nil, fmt.Errorf("Invalid DHCPv6 vendor option: less than 4 bytes")
	}
	code := OptionCode(binary.BigEndian.Uint16(dataStart[:2]))
	length := int(binary.BigEndian.Uint16(dataStart[2:4]))
	if len(dataStart) < length+4 {
		return nil, fmt.Errorf("Invalid option length for vendor option %v. Declared %v, actual %v",
			code, length, len(dataStart)-4,
		)
	}

	optData := dataStart[4 : 4+length]
	if len(optData) < 1 {
		return nil, errors.New("vendParseOption: at least one vendor options data is required")
	}

	return &OptionGeneric{OptionCode: code, OptionData: optData}, nil
}
