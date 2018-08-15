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
	return fmt.Sprintf("%s/i386/%s", AppleVendorID, hwModel), nil
}
