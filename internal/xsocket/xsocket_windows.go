//go:build windows

package xsocket

import "errors"

// CloexecSocket is not supported on Windows.
// Windows does not have the close-on-exec flag concept in the same way Unix does.
func CloexecSocket(domain, typ, proto int) (int, error) {
	return -1, errors.New("xsocket: CloexecSocket not supported on Windows, use net.Dial or net.ListenPacket instead")
}
