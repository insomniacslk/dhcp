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
	Intelx86PC      ArchType = 0
	NECPC98         ArchType = 1
	EFIItanium      ArchType = 2
	DECAlpha        ArchType = 3
	ARCx86          ArchType = 4
	IntelLeanClient ArchType = 5
	EFIIA32         ArchType = 6
	EFIBC           ArchType = 7
	EFIXscale       ArchType = 8
	EFIx8664        ArchType = 9
)

// ArchTypeToStringMap maps an ArchType to a mnemonic name
var ArchTypeToStringMap = map[ArchType]string{
	Intelx86PC:      "Intel x86PC",
	NECPC98:         "NEC/PC98",
	EFIItanium:      "EFI Itanium",
	DECAlpha:        "DEC Alpha",
	ARCx86:          "Arc x86",
	IntelLeanClient: "Intel Lean Client",
	EFIIA32:         "EFI IA32",
	EFIBC:           "EFI BC",
	EFIXscale:       "EFI Xscale",
	EFIx8664:        "EFI x86-64",
}

// OptClientArchType represents an option CLIENT_ARCH_TYPE
type OptClientArchType struct {
	ArchType ArchType
}

func (op *OptClientArchType) Code() OptionCode {
	return OPTION_CLIENT_ARCH_TYPE
}

func (op *OptClientArchType) ToBytes() []byte {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_CLIENT_ARCH_TYPE))
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
