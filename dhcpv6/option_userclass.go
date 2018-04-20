package dhcpv6

// This module defines the OptUserClass structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// OptUserClass represent a DHCPv6 User Class option
type OptUserClass struct {
	UserClasses [][]byte
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
	u16 := make([]byte, 2)
	for _, uc := range op.UserClasses {
		binary.BigEndian.PutUint16(u16, uint16(len(uc)))
		buf = append(buf, u16...)
		buf = append(buf, uc...)
	}
	return buf
}

// Length returns the option length
func (op *OptUserClass) Length() int {
	ret := 0
	for _, uc := range op.UserClasses {
		ret += 2 + len(uc)
	}
	return ret
}

func (op *OptUserClass) String() string {
	ucStrings := make([]string, 0)
	for _, uc := range op.UserClasses {
		ucStrings = append(ucStrings, string(uc))
	}
	return fmt.Sprintf("OptUserClass{userclass=[%s]}", strings.Join(ucStrings, ", "))
}

// ParseOptUserClass builds an OptUserClass structure from a sequence of
// bytes. The input data does not include option code and length bytes.
func ParseOptUserClass(data []byte) (*OptUserClass, error) {
	opt := OptUserClass{}
	for {
		if len(data) == 0 {
			break
		}
		if len(data) < 2 {
			return nil, errors.New("ParseOptUserClass: short data: missing length field")
		}
		ucLen := int(binary.BigEndian.Uint16(data[2:]))
		if len(data) < ucLen+2 {
			return nil, fmt.Errorf("ParseOptUserClass: short data: less than %d bytes", ucLen+2)
		}
		opt.UserClasses = append(opt.UserClasses, data[2:ucLen])
		data = data[2+ucLen:]
	}
	return &opt, nil
}
