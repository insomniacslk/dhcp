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

type OptVendorOpts struct {
	enterpriseNumber uint32
	vendorOpts       []byte
}

func (op *OptVendorOpts) Code() OptionCode {
	return OPTION_VENDOR_OPTS
}

func (op *OptVendorOpts) ToBytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_VENDOR_OPTS))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint32(buf[4:8], uint32(op.enterpriseNumber))
	buf = append(buf, op.vendorOpts...)
	return buf
}

func (op *OptVendorOpts) EnterpriseNumber() uint32 {
	return op.enterpriseNumber
}

func (op *OptVendorOpts) SetEnterpriseNumber(enterpriseNumber uint32) {
	op.enterpriseNumber = enterpriseNumber
}

func (op *OptVendorOpts) VendorOpts() []byte {
	return op.vendorOpts
}

func (op *OptVendorOpts) SetVendorOpts(vendorOpts []byte) {
	op.vendorOpts = append([]byte(nil), vendorOpts...)
}

func (op *OptVendorOpts) Length() int {
	return 4 + len(op.vendorOpts)
}

func (op *OptVendorOpts) String() string {
	return fmt.Sprintf("OptVendorOpts{enterprisenum=%v, vendorOpts=%v}",
		op.enterpriseNumber, op.vendorOpts,
	)
}

// build an OptVendorOpts structure from a sequence of bytes.
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
