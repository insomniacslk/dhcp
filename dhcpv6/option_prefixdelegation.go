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
	options []byte
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
	buf = append(buf, op.options...)
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

func (op *OptIAForPrefixDelegation) Options() []byte {
	return op.options
}

func (op *OptIAForPrefixDelegation) SetOptions(options []byte) {
	op.options = options
}

func (op *OptIAForPrefixDelegation) Length() int {
	return 12 + len(op.options)
}

func (op *OptIAForPrefixDelegation) String() string {
	return fmt.Sprintf("OptIAForPrefixDelegation{IAID=%v, t1=%v, t2=%v, options=%v}",
		op.iaId, op.t1, op.t2, op.options)
}

// build an OptIAForPrefixDelegation structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptIAForPrefixDelegation(data []byte) (*OptIAForPrefixDelegation, error) {
	opt := OptIAForPrefixDelegation{}
	if len(data) < 12 {
		return nil, fmt.Errorf("Invalid IA for Prefix Delegation data length. Expected at least 12 bytes, got %v", len(data))
	}
	copy(opt.iaId[:], data[:4])
	opt.t1 = binary.BigEndian.Uint32(data[4:8])
	opt.t2 = binary.BigEndian.Uint32(data[8:12])
	opt.options = append(opt.options, data[12:]...)
	return &opt, nil
}
