package dhcpv6

// This module defines the OptUserClass structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

// OptUserClass represent a DHCPv6 User Class option
type OptUserClass struct {
	userClass []byte
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
	binary.BigEndian.PutUint16(buf[4:6], uint16(len(op.userClass)))
	buf = append(buf, op.userClass...)
	return buf
}

// UserClass returns the user class as a sequence of bytes
func (op *OptUserClass) UserClass() []byte {
	return op.userClass
}

// SetUserClass sets the user class from a sequence of bytes
func (op *OptUserClass) SetUserClass(userClass []byte) {
	op.userClass = userClass
}

// Length returns the option length
func (op *OptUserClass) Length() int {
	return 2 + len(op.userClass)
}

func (op *OptUserClass) String() string {
	return fmt.Sprintf("OptUserClass{userclass=%s}", string(op.userClass))
}

// ParseOptUserClass builds an OptUserClass structure from a sequence of
// bytes. The input data does not include option code and length bytes.
func ParseOptUserClass(data []byte) (*OptUserClass, error) {
	opt := OptUserClass{}
	opt.userClass = append(opt.userClass, data...)
	return &opt, nil
}
