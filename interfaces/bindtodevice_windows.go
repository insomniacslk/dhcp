package interfaces

// BindToInterface on Windows is a no-op.
// Windows implementation uses ipv4.PacketConn with SetControlMessage for interface filtering
// instead of binding to a specific interface at the socket level.
func BindToInterface(fd int, ifname string) error {
	// No-op on Windows. Interface filtering is handled by ipv4.PacketConn.SetControlMessage.
	return nil
}
