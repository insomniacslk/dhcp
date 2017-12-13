package dhcpv6

// This module defines the OptClientArchType structure.
// https://www.ietf.org/rfc/rfc5970.txt

import (
	"encoding/binary"
	"fmt"
)

type ArchType uint16

// see rfc4578
const (
	INTEL_X86PC ArchType = iota
	NEC_PC98
	EFI_ITANIUM
	DEC_ALPHA
	ARC_X86
	INTEL_LEAN_CLIENT
	EFI_IA32
	EFI_BC
	EFI_XSCALE
	EFI_X86_64
)

var ArchTypeToStringMap = map[ArchType]string{
	INTEL_X86PC:       "Intel x86PC",
	NEC_PC98:          "NEC/PC98",
	EFI_ITANIUM:       "EFI Itanium",
	DEC_ALPHA:         "DEC Alpha",
	ARC_X86:           "Arc x86",
	INTEL_LEAN_CLIENT: "Intel Lean Client",
	EFI_IA32:          "EFI IA32",
	EFI_BC:            "EFI BC",
	EFI_XSCALE:        "EFI Xscale",
	EFI_X86_64:        "EFI x86-64",
}

type OptClientArchType struct {
	archType ArchType
}

func (op *OptClientArchType) Code() OptionCode {
	return OPTION_CLIENT_ARCH_TYPE
}

func (op *OptClientArchType) ToBytes() []byte {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_CLIENT_ARCH_TYPE))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint16(buf[4:6], uint16(op.archType))
	return buf
}

func (op *OptClientArchType) ArchType() ArchType {
	return op.archType
}

func (op *OptClientArchType) SetArchType(archType ArchType) {
	op.archType = archType
}

func (op *OptClientArchType) Length() int {
	return 2
}

func (op *OptClientArchType) String() string {
	name, ok := ArchTypeToStringMap[op.archType]
	if !ok {
		name = "Unknown"
	}
	return fmt.Sprintf("OptClientArchType{archtype=%v}", name)
}

// build an OptClientArchType structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptClientArchType(data []byte) (*OptClientArchType, error) {
	opt := OptClientArchType{}
	if len(data) != 2 {
		return nil, fmt.Errorf("Invalid arch type data length. Expected 2 bytes, got %v", len(data))
	}
	opt.archType = ArchType(binary.BigEndian.Uint16(data))
	return &opt, nil
}
