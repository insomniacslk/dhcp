package dhcpv6

/*
 This module defines the OptVendorOpts structure.
 https://tools.ietf.org/html/rfc3315#section-22.17

       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |      OPTION_VENDOR_OPTS       |           option-len          |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                       enterprise-number                       |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      .                                                               .
      .                          option-data                          .
      .                                                               .
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/

import (
	"encoding/binary"
	"fmt"
)

// OptVendorOpts implements the OptionVendorOpts option
type OptVendorOpts struct {
	enterpriseNumber uint32
	vendorOpts       []byte
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
	binary.BigEndian.PutUint32(buf[4:8], uint32(op.enterpriseNumber))
	buf = append(buf, op.vendorOpts...)
	return buf
}

// EnterpriseNumber returns data in the EnterpriseNumber field
func (op *OptVendorOpts) EnterpriseNumber() uint32 {
	return op.enterpriseNumber
}

// SetEnterpriseNumber sets the EnterpriseNumber
func (op *OptVendorOpts) SetEnterpriseNumber(enterpriseNumber uint32) {
	op.enterpriseNumber = enterpriseNumber
}

// VendorOpts returns data in the VendorOpts field
func (op *OptVendorOpts) VendorOpts() []byte {
	return op.vendorOpts
}

// SetVendorOpts sets the VendorOpts data
func (op *OptVendorOpts) SetVendorOpts(vendorOpts []byte) {
	op.vendorOpts = append([]byte(nil), vendorOpts...)
}

// Length returns the option length in bytes
func (op *OptVendorOpts) Length() int {
	return 4 + len(op.vendorOpts)
}

// String returns the option data in a formatted string
func (op *OptVendorOpts) String() string {
	return fmt.Sprintf("OptVendorOpts{enterprisenum=%v, vendorOpts=%s}",
		op.enterpriseNumber, op.vendorOpts,
	)
}

// ParseOptVendorOpts builds an OptVendorOpts structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptVendorOpts(data []byte) (*OptVendorOpts, error) {
	opt := OptVendorOpts{}
	if len(data) < 4 {
		return nil, fmt.Errorf("Invalid vendor opts data length. Expected at least 4 bytes, got %v", len(data))
	}
	opt.enterpriseNumber = binary.BigEndian.Uint32(data[:4])
	opt.vendorOpts = append([]byte(nil), data[4:]...)
	return &opt, nil
}
