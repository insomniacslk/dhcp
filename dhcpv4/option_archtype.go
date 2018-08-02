package dhcpv4

// This option implements the Client System Architecture Type option
// https://tools.ietf.org/html/rfc4578

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
	ArchTypes []ArchType
}

func (o *OptClientArchType) Code() OptionCode {
	return OptionClientSystemArchitectureType
}

func (o *OptClientArchType) ToBytes() []byte {
	ret := []byte{byte(o.Code()), byte(o.Length())}
	for _, at := range o.ArchTypes {
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf[0:2], uint16(at))
		ret = append(ret, buf...)
	}
	return ret
}

func (o *OptClientArchType) Length() int {
	return 2*len(o.ArchTypes)
}

func (o *OptClientArchType) String() string {
	var archTypes string
	for idx, at := range o.ArchTypes {
		name, ok := ArchTypeToStringMap[at]
		if !ok {
			name = "Unknown"
		}
		archTypes += name
		if idx < len(o.ArchTypes)-1 {
			archTypes += ", "
		}
	}
	return fmt.Sprintf("Client System Architecture Type -> %v", archTypes)
}

// ParseOptClientArchType builds an OptClientArchType structure from
// a sequence of bytes The input data does not include option code and
// length bytes.
func ParseOptClientArchType(data []byte) (*OptClientArchType, error) {
	if len(data) < 2 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionClientSystemArchitectureType {
		return nil, fmt.Errorf("expected code %v, got %v", OptionClientSystemArchitectureType, code)
	}
	length := int(data[1])
	if length == 0 || length%2 != 0 {
		return nil, fmt.Errorf("Invalid length: expected multiple of 2 larger than 2, got %v", length)
	}
	if len(data) < 2+length {
		return nil, ErrShortByteStream
	}
	archTypes := make([]ArchType, 0, length%2)
	for idx := 0; idx < length; idx += 2 {
		b := data[2+idx : 2+idx+2]
		archTypes = append(archTypes, ArchType(binary.BigEndian.Uint16(b)))
	}
	return &OptClientArchType{ArchTypes: archTypes}, nil
}
