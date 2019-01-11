package dhcpv4

import (
	"fmt"
	"net"

	"github.com/u-root/u-root/pkg/uio"
)

// OptRequestedIPAddress implements the requested IP address option described
// by RFC 2132, Section 9.1.
type OptRequestedIPAddress struct {
	RequestedAddr net.IP
}

// ParseOptRequestedIPAddress returns a new OptServerIdentifier from a byte
// stream, or error if any.
func ParseOptRequestedIPAddress(data []byte) (*OptRequestedIPAddress, error) {
	buf := uio.NewBigEndianBuffer(data)
	return &OptRequestedIPAddress{RequestedAddr: net.IP(buf.CopyN(net.IPv4len))}, buf.FinError()
}

// Code returns the option code.
func (o *OptRequestedIPAddress) Code() OptionCode {
	return OptionRequestedIPAddress
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptRequestedIPAddress) ToBytes() []byte {
	return o.RequestedAddr.To4()
}

// String returns a human-readable string.
func (o *OptRequestedIPAddress) String() string {
	return fmt.Sprintf("Requested IP Address -> %v", o.RequestedAddr.String())
}
