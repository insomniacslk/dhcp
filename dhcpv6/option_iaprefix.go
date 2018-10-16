package dhcpv6

// This module defines the OptIAPrefix structure.
// https://www.ietf.org/rfc/rfc3633.txt

import (
	"encoding/binary"
	"fmt"
	"net"
)

type OptIAPrefix struct {
	PreferredLifetime uint32
	ValidLifetime     uint32
	prefixLength      byte
	ipv6Prefix        net.IP
	Options           []Option
}

func (op *OptIAPrefix) Code() OptionCode {
	return OptionIAPrefix
}

func (op *OptIAPrefix) ToBytes() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionIAPrefix))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint32(buf[4:8], op.PreferredLifetime)
	binary.BigEndian.PutUint32(buf[8:12], op.ValidLifetime)
	buf = append(buf, op.prefixLength)
	buf = append(buf, op.ipv6Prefix...)
	for _, opt := range op.Options {
		buf = append(buf, opt.ToBytes()...)
	}
	return buf
}

func (op *OptIAPrefix) PrefixLength() byte {
	return op.prefixLength
}

func (op *OptIAPrefix) SetPrefixLength(pl byte) {
	op.prefixLength = pl
}

// IPv6Prefix returns the ipv6Prefix
func (op *OptIAPrefix) IPv6Prefix() net.IP {
	return op.ipv6Prefix
}

// SetIPv6Prefix sets the ipv6Prefix
func (op *OptIAPrefix) SetIPv6Prefix(p net.IP) {
	op.ipv6Prefix = p
}

// Length returns the option length
func (op *OptIAPrefix) Length() int {
	opLen := 25
	for _, opt := range op.Options {
		opLen += 4 + opt.Length()
	}
	return opLen
}

func (op *OptIAPrefix) String() string {
	return fmt.Sprintf("OptIAPrefix{preferredlifetime=%v, validlifetime=%v, prefixlength=%v, ipv6prefix=%v, options=%v}",
		op.PreferredLifetime, op.ValidLifetime, op.PrefixLength(), op.IPv6Prefix(), op.Options)
}

// GetOneOption will get an option of the give type from the Options field, if
// it is present. It will return `nil` otherwise
func (op *OptIAPrefix) GetOneOption(code OptionCode) Option {
	return getOption(op.Options, code)
}

// DelOption will remove all the options that match a Option code.
func (op *OptIAPrefix) DelOption(code OptionCode) {
	op.Options = delOption(op.Options, code)
}

// build an OptIAPrefix structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptIAPrefix(data []byte) (*OptIAPrefix, error) {
	var err error
	opt := OptIAPrefix{}
	if len(data) < 25 {
		return nil, fmt.Errorf("Invalid IA for Prefix Delegation data length. Expected at least 25 bytes, got %v", len(data))
	}
	opt.PreferredLifetime = binary.BigEndian.Uint32(data[:4])
	opt.ValidLifetime = binary.BigEndian.Uint32(data[4:8])
	opt.prefixLength = data[8]
	opt.ipv6Prefix = net.IP(data[9:25])
	opt.Options, err = OptionsFromBytes(data[25:])
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
