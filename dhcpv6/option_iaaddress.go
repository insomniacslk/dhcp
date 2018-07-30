package dhcpv6

// This module defines the OptIAAddress structure.
// https://www.ietf.org/rfc/rfc3633.txt

import (
	"encoding/binary"
	"fmt"
	"net"
)

// OptIAAddress represents an OptionIAAddr
type OptIAAddress struct {
	IPv6Addr          net.IP
	PreferredLifetime uint32
	ValidLifetime     uint32
	Options           []Option
}

// Code returns the option's code
func (op *OptIAAddress) Code() OptionCode {
	return OptionIAAddr
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptIAAddress) ToBytes() []byte {
	buf := make([]byte, 28)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionIAAddr))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	copy(buf[4:20], op.IPv6Addr[:])
	binary.BigEndian.PutUint32(buf[20:24], op.PreferredLifetime)
	binary.BigEndian.PutUint32(buf[24:28], op.ValidLifetime)
	for _, opt := range op.Options {
		buf = append(buf, opt.ToBytes()...)
	}
	return buf
}

// Length returns the option length
func (op *OptIAAddress) Length() int {
	opLen := 24
	for _, opt := range op.Options {
		opLen += 4 + opt.Length()
	}
	return opLen
}

func (op *OptIAAddress) String() string {
	return fmt.Sprintf("OptIAAddress{ipv6addr=%v, preferredlifetime=%v, validlifetime=%v, options=%v}",
		net.IP(op.IPv6Addr[:]), op.PreferredLifetime, op.ValidLifetime, op.Options)
}

// ParseOptIAAddress builds an OptIAAddress structure from a sequence
// of bytes. The input data does not include option code and length
// bytes.
func ParseOptIAAddress(data []byte) (*OptIAAddress, error) {
	var err error
	opt := OptIAAddress{}
	if len(data) < 24 {
		return nil, fmt.Errorf("Invalid IA Address data length. Expected at least 24 bytes, got %v", len(data))
	}
	opt.IPv6Addr = net.IP(data[:16])
	opt.PreferredLifetime = binary.BigEndian.Uint32(data[16:20])
	opt.ValidLifetime = binary.BigEndian.Uint32(data[20:24])
	opt.Options, err = OptionsFromBytes(data[24:])
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
