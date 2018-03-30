package bsdp

import (
	"fmt"
	"syscall"
)

// MakeVendorClassIdentifier calls the sysctl syscall on macOS to get the
// platform model.
func MakeVendorClassIdentifier() (string, error) {
	// Fetch hardware model for class ID.
	hwModel, err := syscall.Sysctl("hw.model")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("AAPLBSDPC/i386/%s", hwModel), nil
}
