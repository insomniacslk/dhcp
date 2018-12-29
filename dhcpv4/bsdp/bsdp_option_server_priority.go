package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// OptServerPriority represents an option encapsulating the server priority.
type OptServerPriority struct {
	Priority uint16
}

// ParseOptServerPriority returns a new OptServerPriority from a byte stream, or
// error if any.
func ParseOptServerPriority(data []byte) (*OptServerPriority, error) {
	buf := uio.NewBigEndianBuffer(data)
	return &OptServerPriority{Priority: buf.Read16()}, buf.FinError()
}

// Code returns the option code.
func (o *OptServerPriority) Code() dhcpv4.OptionCode {
	return OptionServerPriority
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptServerPriority) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(o.Priority)
	return buf.Data()
}

// String returns a human-readable string.
func (o *OptServerPriority) String() string {
	return fmt.Sprintf("BSDP Server Priority -> %v", o.Priority)
}
