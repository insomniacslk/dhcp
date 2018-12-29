package iana

// Arch encodes an architecture type per RFC 4578, Section 2.1.
type Arch uint16

// See RFC 4578.
const (
	INTEL_X86PC       Arch = 0
	NEC_PC98          Arch = 1
	EFI_ITANIUM       Arch = 2
	DEC_ALPHA         Arch = 3
	ARC_X86           Arch = 4
	INTEL_LEAN_CLIENT Arch = 5
	EFI_IA32          Arch = 6
	EFI_BC            Arch = 7
	EFI_XSCALE        Arch = 8
	EFI_X86_64        Arch = 9
)

// archTypeToStringMap maps an Arch to a mnemonic name
var archTypeToStringMap = map[Arch]string{
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

// String returns a mnemonic name for a given architecture type.
func (a Arch) String() string {
	if at := archTypeToStringMap[a]; at != "" {
		return at
	}
	return "unknown"
}
