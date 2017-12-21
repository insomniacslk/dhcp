package dhcpv6

// This module defines the OptBootFileURL structure.
// https://www.ietf.org/rfc/rfc5970.txt

import (
	"encoding/binary"
	"fmt"
)

type OptBootFileURL struct {
	bootFileUrl []byte
}

func (op *OptBootFileURL) Code() OptionCode {
	return OPT_BOOTFILE_URL
}

func (op *OptBootFileURL) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPT_BOOTFILE_URL))
	binary.BigEndian.PutUint16(buf[2:4], 2)
	buf = append(buf, op.bootFileUrl...)
	return buf
}

func (op *OptBootFileURL) BootFileURL() []byte {
	return op.bootFileUrl
}

func (op *OptBootFileURL) SetBootFileURL(bootFileUrl []byte) {
	op.bootFileUrl = bootFileUrl
}

func (op *OptBootFileURL) Length() int {
	return len(op.bootFileUrl)
}

func (op *OptBootFileURL) String() string {
	return fmt.Sprintf("OptBootFileURL{BootFileUrl=%v}", op.bootFileUrl)
}

// build an OptBootFileURL structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptBootFileURL(data []byte) (*OptBootFileURL, error) {
	opt := OptBootFileURL{}
	opt.bootFileUrl = append([]byte(nil), data...)
	return &opt, nil
}
