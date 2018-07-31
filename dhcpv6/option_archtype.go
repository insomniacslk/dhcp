package dhcpv6

// This module defines the OptClientArchType structure.
// https://www.ietf.org/rfc/rfc5970.txt

import (
	"encoding/binary"
	"fmt"
)

//ArchType encodes an architecture type in an uint16
type ArchType uint16

// see rfc4578
const (
	INTEL_X86PC       ArchType = 0
	NEC_PC98          ArchType = 1
	EFI_ITANIUM       ArchType = 2
	DEC_ALPHA         ArchType = 3
	ARC_X86           ArchType = 4
	INTEL_LEAN_CLIENT ArchType = 5
	EFI_IA32          ArchType = 6
	EFI_BC            ArchType = 7
	EFI_XSCALE        ArchType = 8
	EFI_X86_64        ArchType = 9
)

// ArchTypeToStringMap maps an ArchType to a mnemonic name
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

// OptClientArchType represents an option CLIENT_ARCH_TYPE
type OptClientArchType struct {
	ArchType ArchType
}

func (op *OptClientArchType) Code() OptionCode {
	return OptionClientArchType
}

func (op *OptClientArchType) ToBytes() []byte {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionClientArchType))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint16(buf[4:6], uint16(op.ArchType))
	return buf
}

func (op *OptClientArchType) Length() int {
	return 2
}

func (op *OptClientArchType) String() string {
	name, ok := ArchTypeToStringMap[op.ArchType]
	if !ok {
		name = "Unknown"
	}
	return fmt.Sprintf("OptClientArchType{archtype=%v}", name)
}

// ParseOptClientArchType builds an OptClientArchType structure from
// a sequence of bytes The input data does not include option code and
// length bytes.
func ParseOptClientArchType(data []byte) (*OptClientArchType, error) {
	opt := OptClientArchType{}
	if len(data) != 2 {
		return nil, fmt.Errorf("Invalid arch type data length. Expected 2 bytes, got %v", len(data))
	}
	opt.ArchType = ArchType(binary.BigEndian.Uint16(data))
	return &opt, nil
}
