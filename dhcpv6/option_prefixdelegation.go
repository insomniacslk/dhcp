package dhcpv6

// This module defines the OptIAForPrefixDelegation structure.
// https://www.ietf.org/rfc/rfc3633.txt

import (
	"encoding/binary"
	"fmt"
)

type OptIAForPrefixDelegation struct {
	iaId    [4]byte
	t1      uint32
	t2      uint32
	Options []Option
}

func (op *OptIAForPrefixDelegation) Code() OptionCode {
	return OptionIAPD
}

func (op *OptIAForPrefixDelegation) ToBytes() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionIAPD))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	copy(buf[4:8], op.iaId[:])
	binary.BigEndian.PutUint32(buf[8:12], op.t1)
	binary.BigEndian.PutUint32(buf[12:16], op.t2)
	for _, opt := range op.Options {
		buf = append(buf, opt.ToBytes()...)
	}
	return buf
}

func (op *OptIAForPrefixDelegation) IAID() []byte {
	return op.iaId[:]
}

func (op *OptIAForPrefixDelegation) SetIAID(iaId [4]byte) {
	op.iaId = iaId
}

func (op *OptIAForPrefixDelegation) T1() uint32 {
	return op.t1
}

func (op *OptIAForPrefixDelegation) SetT1(t1 uint32) {
	op.t1 = t1
}

func (op *OptIAForPrefixDelegation) T2() uint32 {
	return op.t2
}

func (op *OptIAForPrefixDelegation) SetT2(t2 uint32) {
	op.t2 = t2
}

func (op *OptIAForPrefixDelegation) Length() int {
	l := 12
	for _, opt := range op.Options {
		l += 4 + opt.Length()
	}
	return l
}

func (op *OptIAForPrefixDelegation) String() string {
	return fmt.Sprintf("OptIAForPrefixDelegation{IAID=%v, t1=%v, t2=%v, options=%v}",
		op.iaId, op.t1, op.t2, op.Options)
}

// GetOneOption will get an option of the give type from the Options field, if
// it is present. It will return `nil` otherwise
func (op *OptIAForPrefixDelegation) GetOneOption(code OptionCode) Option {
	return getOption(op.Options, code)
}

// DelOption will remove all the options that match a Option code.
func (op *OptIAForPrefixDelegation) DelOption(code OptionCode) {
	op.Options = delOption(op.Options, code)
}

// build an OptIAForPrefixDelegation structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptIAForPrefixDelegation(data []byte) (*OptIAForPrefixDelegation, error) {
	var err error
	opt := OptIAForPrefixDelegation{}
	if len(data) < 12 {
		return nil, fmt.Errorf("Invalid IA for Prefix Delegation data length. Expected at least 12 bytes, got %v", len(data))
	}
	copy(opt.iaId[:], data[:4])
	opt.t1 = binary.BigEndian.Uint32(data[4:8])
	opt.t2 = binary.BigEndian.Uint32(data[8:12])
	opt.Options, err = OptionsFromBytes(data[12:])
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
