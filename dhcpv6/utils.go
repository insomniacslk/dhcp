package dhcpv6

import (
	"errors"
	"strings"
)

// IsNetboot function takes a DHCPv6 message and returns true if the machine
// is trying to netboot. It checks if "boot file" is one of the requested
// options, which is useful for SOLICIT/REQUEST packet types, it also checks
// if the "boot file" option is included in the packet, which is useful for
// ADVERTISE/REPLY packet.
func IsNetboot(msg DHCPv6) bool {
	if IsRequested(msg, OptionBootfileURL) {
		return true
	}
	if optbf := msg.GetOneOption(OptionBootfileURL); optbf != nil {
		return true
	}
	return false
}

// IsRequested function takes a DHCPv6 message and an OptionCode, and returns
// true if that option is within the requested options of the DHCPv6 message.
func IsRequested(msg DHCPv6, requested OptionCode) bool {
	for _, optoro := range msg.GetOption(OptionORO) {
		for _, o := range optoro.(*OptRequestedOption).RequestedOptions() {
			if o == requested {
				return true
			}
		}
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
	if opt := msg.GetOneOption(OptionClientArchType); opt != nil {
		optat := opt.(*OptClientArchType)
		// TODO investigate if other types are appropriate
		if optat.ArchType == EFI_BC || optat.ArchType == EFI_X86_64 {
			return true
		}
	}
	if opt := msg.GetOneOption(OptionUserClass); opt != nil {
		optuc := opt.(*OptUserClass)
		for _, uc := range optuc.UserClasses {
			if strings.Contains(string(uc), "EFI") {
				return true
			}
		}
	}
	return false
}

// GetTransactionID returns a transactionID of a message or its inner message
// in case of relay
func GetTransactionID(packet DHCPv6) (uint32, error) {
	if message, ok := packet.(*DHCPv6Message); ok {
		return message.TransactionID(), nil
	}
	if relay, ok := packet.(*DHCPv6Relay); ok {
		message, err := relay.GetInnerMessage()
		if err != nil {
			return 0, err
		}
		return GetTransactionID(message)
	}
	return 0, errors.New("Invalid DHCPv6 packet")
}
