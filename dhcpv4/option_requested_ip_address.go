package dhcpv4

import (
	"fmt"
	"net"

	"github.com/u-root/u-root/pkg/uio"
)

// This option implements the requested IP address option
// https://tools.ietf.org/html/rfc2132

// OptRequestedIPAddress represents an option encapsulating the server
// identifier.
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

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptRequestedIPAddress) Length() int {
	return len(o.RequestedAddr.To4())
}
