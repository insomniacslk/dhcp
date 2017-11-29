package options

// This module defines the OptIAAddress structure.
// https://www.ietf.org/rfc/rfc3633.txt

import (
	"encoding/binary"
	"fmt"
	"net"
)

type OptIAAddress struct {
	ipv6Addr          [16]byte
	preferredLifetime uint32
	validLifetime     uint32
	options           []byte
}

func (op *OptIAAddress) Code() OptionCode {
	return OPTION_IAADDR
}

func (op *OptIAAddress) ToBytes() []byte {
	buf := make([]byte, 28)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_IAADDR))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	copy(buf[4:20], op.ipv6Addr[:])
	binary.BigEndian.PutUint32(buf[20:24], op.preferredLifetime)
	binary.BigEndian.PutUint32(buf[24:28], op.validLifetime)
	buf = append(buf, op.options...)
	return buf
}

func (op *OptIAAddress) IPv6Addr() []byte {
	return op.ipv6Addr[:]
}

func (op *OptIAAddress) SetIPv6Addr(addr [16]byte) {
	op.ipv6Addr = addr
}

func (op *OptIAAddress) PreferredLifetime() uint32 {
	return op.preferredLifetime
}

func (op *OptIAAddress) SetPreferredLifetime(pl uint32) {
	op.preferredLifetime = pl
}

func (op *OptIAAddress) ValidLifetime() uint32 {
	return op.validLifetime
}

func (op *OptIAAddress) SetValidLifetime(vl uint32) {
	op.validLifetime = vl
}
func (op *OptIAAddress) Options() []byte {
	return op.options
}

func (op *OptIAAddress) SetOptions(options []byte) {
	op.options = options
}

func (op *OptIAAddress) Length() int {
	return 24 + len(op.options)
}

func (op *OptIAAddress) String() string {
	return fmt.Sprintf("OptIAAddress{ipv6addr=%v, preferredlifetime=%v, validlifetime=%v, options=%v}",
		net.IP(op.ipv6Addr[:]), op.preferredLifetime, op.validLifetime, op.options)
}

// build an OptIAAddress structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptIAAddress(data []byte) (*OptIAAddress, error) {
	opt := OptIAAddress{}
	if len(data) < 24 {
		return nil, fmt.Errorf("Invalid IA Address data length. Expected at least 24 bytes, got %v", len(data))
	}
	copy(opt.ipv6Addr[:], data[:16])
	opt.preferredLifetime = binary.BigEndian.Uint32(data[16:20])
	opt.validLifetime = binary.BigEndian.Uint32(data[20:24])
	copy(opt.options, data[24:])
	return &opt, nil
}
