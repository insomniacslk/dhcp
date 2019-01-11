package dhcpv4

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptMaximumDHCPMessageSize implements the maximum DHCP message size option
// described by RFC 2132, Section 9.10.
type OptMaximumDHCPMessageSize struct {
	Size uint16
}

// ParseOptMaximumDHCPMessageSize constructs an OptMaximumDHCPMessageSize
// struct from a sequence of bytes and returns it, or an error.
func ParseOptMaximumDHCPMessageSize(data []byte) (*OptMaximumDHCPMessageSize, error) {
	buf := uio.NewBigEndianBuffer(data)
	return &OptMaximumDHCPMessageSize{Size: buf.Read16()}, buf.FinError()
}

// Code returns the option code.
func (o *OptMaximumDHCPMessageSize) Code() OptionCode {
	return OptionMaximumDHCPMessageSize
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptMaximumDHCPMessageSize) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(o.Size)
	return buf.Data()
}

// String returns a human-readable string for this option.
func (o *OptMaximumDHCPMessageSize) String() string {
	return fmt.Sprintf("Maximum DHCP Message Size -> %v", o.Size)
}
