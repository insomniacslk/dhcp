package dhcpv6

import (
	"strings"
)

// IsNetboot function takes a DHCPv6 message and returns true if the machine
// is trying to netboot. It checks if "boot file" is one of the requested
// options, which is useful for SOLICIT/REQUEST packet types, it also checks
// if the "boot file" option is included in the packet, which is useful for
// ADVERTISE/REPLY packet.
func IsNetboot(msg DHCPv6) bool {
	for _, optoro := range msg.GetOption(OPTION_ORO) {
		for _, o := range optoro.(*OptRequestedOption).RequestedOptions() {
			if o == OPT_BOOTFILE_URL {
				return true
			}
		}
	}
	if optbf := msg.GetOneOption(OPT_BOOTFILE_URL); optbf != nil {
		return true
	}
	return false
}

// IsUsingUEFI function takes a DHCPv6 message and returns true if
// the machine trying to netboot is using UEFI of false if it is not.
func IsUsingUEFI(msg DHCPv6) bool {
	// RFC 4578 says:
	// As of the writing of this document, the following pre-boot
	//    architecture types have been requested.
	//             Type   Architecture Name
	//             ----   -----------------
	//               0    Intel x86PC
	//               1    NEC/PC98
	//               2    EFI Itanium
	//               3    DEC Alpha
	//               4    Arc x86
	//               5    Intel Lean Client
	//               6    EFI IA32
	//               7    EFI BC
	//               8    EFI Xscale
	//               9    EFI x86-64
	if opt := msg.GetOneOption(OPTION_CLIENT_ARCH_TYPE); opt != nil {
		optat := opt.(*OptClientArchType)
		if optat.ArchType == EFI_BC || optat.ArchType == EFI_X86_64 {
			return true
		}
	}
	// our iPXE roms have been built to include the architecture type in the
	// user_class field. e.g. FBipxeUEFI or FBipxeLegacy
	if opt := msg.GetOneOption(OPTION_USER_CLASS); opt != nil {
		optuc := opt.(*OptUserClass)
		for _, uc := range optuc.UserClasses {
			if strings.Contains(string(uc), "EFI") {
				return true
			}
		}
	}
	return false
}
