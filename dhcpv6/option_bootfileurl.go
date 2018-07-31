package dhcpv6

// This module defines the OptBootFileURL structure.
// https://www.ietf.org/rfc/rfc5970.txt

import (
	"encoding/binary"
	"fmt"
)

// OptBootFileURL implements the OptionBootfileURL option
type OptBootFileURL struct {
	BootFileURL []byte
}

// Code returns the option code
func (op *OptBootFileURL) Code() OptionCode {
	return OptionBootfileURL
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptBootFileURL) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionBootfileURL))
	binary.BigEndian.PutUint16(buf[2:4], uint16(len(op.BootFileURL)))
	buf = append(buf, op.BootFileURL...)
	return buf
}

// Length returns the option length in bytes
func (op *OptBootFileURL) Length() int {
	return len(op.BootFileURL)
}

func (op *OptBootFileURL) String() string {
	return fmt.Sprintf("OptBootFileURL{BootFileUrl=%s}", op.BootFileURL)
}

// ParseOptBootFileURL builds an OptBootFileURL structure from a sequence
// of bytes. The input data does not include option code and length bytes.
func ParseOptBootFileURL(data []byte) (*OptBootFileURL, error) {
	opt := OptBootFileURL{}
	opt.BootFileURL = append([]byte(nil), data...)
	return &opt, nil
}
