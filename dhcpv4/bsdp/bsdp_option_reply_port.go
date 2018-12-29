package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// OptReplyPort represents a BSDP protocol version.
//
// Implements the BSDP option reply port. This is used when BSDP responses
// should be sent to a reply port other than the DHCP default. The macOS GUI
// "Startup Disk Select" sends this option since it's operating in an
// unprivileged context.
type OptReplyPort struct {
	Port uint16
}

// ParseOptReplyPort constructs an OptReplyPort struct from a sequence of
// bytes and returns it, or an error.
func ParseOptReplyPort(data []byte) (*OptReplyPort, error) {
	buf := uio.NewBigEndianBuffer(data)
	return &OptReplyPort{buf.Read16()}, buf.FinError()
}

// Code returns the option code.
func (o *OptReplyPort) Code() dhcpv4.OptionCode {
	return OptionReplyPort
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptReplyPort) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(o.Port)
	return buf.Data()
}

// String returns a human-readable string for this option.
func (o *OptReplyPort) String() string {
	return fmt.Sprintf("BSDP Reply Port -> %v", o.Port)
}

// Length returns the length of the data portion of this option.
func (o *OptReplyPort) Length() int {
	return 2
}
