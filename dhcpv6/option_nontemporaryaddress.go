package dhcpv6

// This module defines the OptIANA structure.
// https://www.ietf.org/rfc/rfc3633.txt

import (
	"encoding/binary"
	"fmt"
)

type OptIANA struct {
	iaId    [4]byte
	t1      uint32
	t2      uint32
	options []byte
}

func (op *OptIANA) Code() OptionCode {
	return OPTION_IA_NA
}

func (op *OptIANA) ToBytes() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_IA_NA))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	copy(buf[4:8], op.iaId[:])
	binary.BigEndian.PutUint32(buf[8:12], op.t1)
	binary.BigEndian.PutUint32(buf[12:16], op.t2)
	buf = append(buf, op.options...)
	return buf
}

func (op *OptIANA) IAID() []byte {
	return op.iaId[:]
}

func (op *OptIANA) SetIAID(iaId [4]byte) {
	op.iaId = iaId
}

func (op *OptIANA) T1() uint32 {
	return op.t1
}

func (op *OptIANA) SetT1(t1 uint32) {
	op.t1 = t1
}

func (op *OptIANA) T2() uint32 {
	return op.t2
}

func (op *OptIANA) SetT2(t2 uint32) {
	op.t2 = t2
}

func (op *OptIANA) Options() []byte {
	return op.options
}

func (op *OptIANA) SetOptions(options []byte) {
	op.options = options
}

func (op *OptIANA) Length() int {
	return 12 + len(op.options)
}

func (op *OptIANA) String() string {
	return fmt.Sprintf("OptIANA{IAID=%v, t1=%v, t2=%v, options=%v}",
		op.iaId, op.t1, op.t2, op.options)
}

// build an OptIANA structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptIANA(data []byte) (*OptIANA, error) {
	opt := OptIANA{}
	if len(data) < 12 {
		return nil, fmt.Errorf("Invalid IA for Non-temporary Addresses data length. Expected at least 12 bytes, got %v", len(data))
	}
	copy(opt.iaId[:], data[:4])
	opt.t1 = binary.BigEndian.Uint32(data[4:8])
	opt.t2 = binary.BigEndian.Uint32(data[8:12])
	opt.options = append(data[12:])
	return &opt, nil
}
