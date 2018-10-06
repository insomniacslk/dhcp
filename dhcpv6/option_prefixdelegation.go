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
	options []Option
}

// Code returns the option code
func (op *OptIAForPrefixDelegation) Code() OptionCode {
	return OptionIAPD
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptIAForPrefixDelegation) ToBytes() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionIAPD))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	copy(buf[4:8], op.iaId[:])
	binary.BigEndian.PutUint32(buf[8:12], op.t1)
	binary.BigEndian.PutUint32(buf[12:16], op.t2)
	for _, opt := range op.options {
		buf = append(buf, opt.ToBytes()...)
	}
	return buf
}

// IAID returns the identity association identifier for this option
func (op *OptIAForPrefixDelegation) IAID() []byte {
	return op.iaId[:]
}

// SetIAID sets the identity association identifier for this option
func (op *OptIAForPrefixDelegation) SetIAID(iaId [4]byte) {
	op.iaId = iaId
}

// T1 returns the T1 timer for this option
func (op *OptIAForPrefixDelegation) T1() uint32 {
	return op.t1
}

// SetT1 sets the T1 timer for this option
func (op *OptIAForPrefixDelegation) SetT1(t1 uint32) {
	op.t1 = t1
}

// T2 returns the T2 timer for this option
func (op *OptIAForPrefixDelegation) T2() uint32 {
	return op.t2
}

// SetT2 sets the T2 timer for this option
func (op *OptIAForPrefixDelegation) SetT2(t2 uint32) {
	op.t2 = t2
}

// Options serializes the options and returns them as a sequence of bytes
func (op *OptIAForPrefixDelegation) Options() []byte {
	log.Printf("Warning: OptIAForPrefixDelegation.Options() is deprecated and will be changed to a public field")
	buf := op.ToBytes()
	return buf[16:]
}

// SetOptions sets the options as a sequence of bytes
func (op *OptIAForPrefixDelegation) SetOptions(options []byte) error {
	var err error
	op.options, err = OptionsFromBytes(options)
	if err != nil {
		return err
	}
	return nil
}

// Length returns the option length
func (op *OptIAForPrefixDelegation) Length() int {
	l := 12
	for _, opt := range op.options {
		l += 4 + opt.Length()
	}
	return l
}

// String returns a string representation of the OptIAForPrefixDelegation data
func (op *OptIAForPrefixDelegation) String() string {
	return fmt.Sprintf("OptIAForPrefixDelegation{IAID=%v, t1=%v, t2=%v, options=%v}",
		op.iaId, op.t1, op.t2, op.options)
}

// GetOneOption will get an option of the give type from the Options field, if
// it is present. It will return `nil` otherwise
func (op *OptIAForPrefixDelegation) GetOneOption(code OptionCode) Option {
	return getOption(op.options, code)
}

// DelOption will remove all the options that match a Option code.
func (op *OptIAForPrefixDelegation) DelOption(code OptionCode) {
	op.options = delOption(op.options, code)
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
	opt.options, err = OptionsFromBytes(data[12:])
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
