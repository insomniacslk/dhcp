package dhcpv6

// This module defines the OptUserClass structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

// OptUserClass represent a DHCPv6 User Class option
type OptUserClass struct {
	UserClass []byte
}

// Code returns the option code
func (op *OptUserClass) Code() OptionCode {
	return OPTION_USER_CLASS
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptUserClass) ToBytes() []byte {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_USER_CLASS))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	// user-class-data has an internal data length field too..
	binary.BigEndian.PutUint16(buf[4:6], uint16(len(op.UserClass)))
	buf = append(buf, op.UserClass...)
	return buf
}

// Length returns the option length
func (op *OptUserClass) Length() int {
	return 2 + len(op.UserClass)
}

func (op *OptUserClass) String() string {
	return fmt.Sprintf("OptUserClass{userclass=%s}", string(op.UserClass))
}

// ParseOptUserClass builds an OptUserClass structure from a sequence of
// bytes. The input data does not include option code and length bytes.
func ParseOptUserClass(data []byte) (*OptUserClass, error) {
	opt := OptUserClass{}
	dataLen := int(binary.BigEndian.Uint16(data[:2]))
	if dataLen != len(data)-2 {
		return nil, fmt.Errorf("ParseOptUserClass: declared data length does not match actual length: %d != %d", dataLen, len(data)-2)
	}
	opt.UserClass = append(opt.UserClass, data[2:]...)
	return &opt, nil
}
