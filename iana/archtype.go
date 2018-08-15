package iana

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


// ArchTypeToString returns a mnemonic name for a given architecture type
func ArchTypeToString(a ArchType) string {
	if at := ArchTypeToStringMap[a]; at != "" {
		return at
	}
	return "Unknown"
}
