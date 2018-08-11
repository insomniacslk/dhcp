package dhcpv6

// This module defines the OptIANA structure.
// https://www.ietf.org/rfc/rfc3633.txt

import (
	"encoding/binary"
	"fmt"
)

type OptIANA struct {
	IaId    [4]byte
	T1      uint32
	T2      uint32
	Options []Option
}

func (op *OptIANA) Code() OptionCode {
	return OptionIANA
}

func (op *OptIANA) ToBytes() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionIANA))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	copy(buf[4:8], op.IaId[:])
	binary.BigEndian.PutUint32(buf[8:12], op.T1)
	binary.BigEndian.PutUint32(buf[12:16], op.T2)
	for _, opt := range op.Options {
		buf = append(buf, opt.ToBytes()...)
	}
	return buf
}

func (op *OptIANA) Length() int {
	l := 12
	for _, opt := range op.Options {
		l += 4 + opt.Length()
	}
	return l
}

func (op *OptIANA) String() string {
	return fmt.Sprintf("OptIANA{IAID=%v, t1=%v, t2=%v, options=%v}",
		op.IaId, op.T1, op.T2, op.Options)
}

// GetOneOption will get an option of the give type from the Options field, if
// it is present. It will return `nil` otherwise
func (op *OptIANA) GetOneOption(code OptionCode) Option {
	return getOption(op.Options, code)
}

// DelOption will remove all the options that match a Option code.
func (op *OptIANA) DelOption(code OptionCode) {
	op.Options = delOption(op.Options, code)
}

// build an OptIANA structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptIANA(data []byte) (*OptIANA, error) {
	var err error
	opt := OptIANA{}
	if len(data) < 12 {
		return nil, fmt.Errorf("Invalid IA for Non-temporary Addresses data length. Expected at least 12 bytes, got %v", len(data))
	}
	copy(opt.IaId[:], data[:4])
	opt.T1 = binary.BigEndian.Uint32(data[4:8])
	opt.T2 = binary.BigEndian.Uint32(data[8:12])
	opt.Options, err = OptionsFromBytes(data[12:])
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
