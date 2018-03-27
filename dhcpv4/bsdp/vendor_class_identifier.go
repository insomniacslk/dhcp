// +build !darwin

package bsdp

// MakeVendorClassIdentifier calls the sysctl syscall on macOS to get the
// platform model.
func MakeVendorClassIdentifier() (string, error) {
	return DefaultMacOSVendorClassIdentifier, nil
}
