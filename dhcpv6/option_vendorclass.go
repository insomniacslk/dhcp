package dhcpv6

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// OptVendorClass represents a DHCPv6 Vendor Class option
type OptVendorClass struct {
	EnterpriseNumber uint32
	Data             [][]byte
}

// Code returns the option code
func (op *OptVendorClass) Code() OptionCode {
	return OptionVendorClass
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptVendorClass) ToBytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionVendorClass))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint32(buf[4:8], uint32(op.EnterpriseNumber))
	u16 := make([]byte, 2)
	for _, data := range op.Data {
		binary.BigEndian.PutUint16(u16, uint16(len(data)))
		buf = append(buf, u16...)
		buf = append(buf, data...)
	}
	return buf
}

// Length returns the option length
func (op *OptVendorClass) Length() int {
	ret := 0
	for _, data := range op.Data {
		ret += 2 + len(data)
	}
	return 4 + ret
}

// String returns a string representation of the VendorClass data
func (op *OptVendorClass) String() string {
	vcStrings := make([]string, 0)
	for _, data := range op.Data {
		vcStrings = append(vcStrings, string(data))
	}
	return fmt.Sprintf("OptVendorClass{enterprisenum=%d, data=[%s]}", op.EnterpriseNumber, strings.Join(vcStrings, ", "))
}

// ParseOptVendorClass builds an OptVendorClass structure from a sequence of
// bytes. The input data does not include option code and length bytes.
func ParseOptVendorClass(data []byte) (*OptVendorClass, error) {
	opt := OptVendorClass{}
	if len(data) < 4 {
		return nil, fmt.Errorf("Invalid vendor opts data length. Expected at least 4 bytes, got %v", len(data))
	}
	opt.EnterpriseNumber = binary.BigEndian.Uint32(data[:4])
	data = data[4:]
	for {
		if len(data) == 0 {
			break
		}
		if len(data) < 2 {
			return nil, errors.New("ParseOptVendorClass: short data: missing length field")
		}
		vcLen := int(binary.BigEndian.Uint16(data[:2]))
		if len(data) < vcLen+2 {
			return nil, fmt.Errorf("ParseOptVendorClass: short data: less than %d bytes", vcLen+2)
		}
		opt.Data = append(opt.Data, data[2:vcLen+2])
		data = data[2+vcLen:]
	}
	return &opt, nil
}
