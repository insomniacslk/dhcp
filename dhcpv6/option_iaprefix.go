package dhcpv6

// This module defines the OptIAPrefix structure.
// https://www.ietf.org/rfc/rfc3633.txt

import (
	"encoding/binary"
	"fmt"
	"net"
)

type OptIAPrefix struct {
	preferredLifetime uint32
	validLifetime     uint32
	prefixLength      byte
	ipv6Prefix        [16]byte
	options           []byte
}

func (op *OptIAPrefix) Code() OptionCode {
	return OPTION_IAPREFIX
}

func (op *OptIAPrefix) ToBytes() []byte {
	buf := make([]byte, 25)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_IAPREFIX))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint32(buf[4:8], op.preferredLifetime)
	binary.BigEndian.PutUint32(buf[8:12], op.validLifetime)
	buf = append(buf, op.prefixLength)
	buf = append(buf, op.ipv6Prefix[:]...)
	buf = append(buf, op.options...)
	return buf
}

func (op *OptIAPrefix) PreferredLifetime() uint32 {
	return op.preferredLifetime
}

func (op *OptIAPrefix) SetPreferredLifetime(pl uint32) {
	op.preferredLifetime = pl
}

func (op *OptIAPrefix) ValidLifetime() uint32 {
	return op.validLifetime
}

func (op *OptIAPrefix) SetValidLifetime(vl uint32) {
	op.validLifetime = vl
}

func (op *OptIAPrefix) PrefixLength() byte {
	return op.prefixLength
}

func (op *OptIAPrefix) SetPrefixLength(pl byte) {
	op.prefixLength = pl
}

func (op *OptIAPrefix) IPv6Prefix() []byte {
	return op.ipv6Prefix[:]
}

func (op *OptIAPrefix) SetIPv6Prefix(p [16]byte) {
	op.ipv6Prefix = p
}

func (op *OptIAPrefix) Options() []byte {
	return op.options
}

func (op *OptIAPrefix) SetOptions(options []byte) {
	op.options = options
}

func (op *OptIAPrefix) Length() int {
	return 25 + len(op.options)
}

func (op *OptIAPrefix) String() string {
	return fmt.Sprintf("OptIAPrefix{preferredlifetime=%v, validlifetime=%v, prefixlength=%v, ipv6prefix=%v, options=%v}",
		op.preferredLifetime, op.validLifetime, op.prefixLength, net.IP(op.ipv6Prefix[:]), op.options)
}

// build an OptIAPrefix structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptIAPrefix(data []byte) (*OptIAPrefix, error) {
	opt := OptIAPrefix{}
	if len(data) < 12 {
		return nil, fmt.Errorf("Invalid IA for Prefix Delegation data length. Expected at least 12 bytes, got %v", len(data))
	}
	opt.preferredLifetime = binary.BigEndian.Uint32(data[:4])
	opt.validLifetime = binary.BigEndian.Uint32(data[4:8])
	opt.prefixLength = data[9]
	copy(opt.ipv6Prefix[:], data[9:17])
	copy(opt.options, data[17:])
	return &opt, nil
}
