package dhcpv4

import (
	"fmt"
	"net"

	"github.com/u-root/u-root/pkg/uio"
)

// This option implements the server identifier option
// https://tools.ietf.org/html/rfc2132

// OptBroadcastAddress represents an option encapsulating the server identifier.
type OptBroadcastAddress struct {
	BroadcastAddress net.IP
}

// ParseOptBroadcastAddress returns a new OptBroadcastAddress from a byte
// stream, or error if any.
func ParseOptBroadcastAddress(data []byte) (*OptBroadcastAddress, error) {
	buf := uio.NewBigEndianBuffer(data)
	return &OptBroadcastAddress{BroadcastAddress: net.IP(buf.CopyN(net.IPv4len))}, buf.FinError()
}

// Code returns the option code.
func (o *OptBroadcastAddress) Code() OptionCode {
	return OptionBroadcastAddress
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptBroadcastAddress) ToBytes() []byte {
	ret := []byte{byte(o.Code()), byte(o.Length())}
	return append(ret, o.BroadcastAddress.To4()...)
}

// String returns a human-readable string.
func (o *OptBroadcastAddress) String() string {
	return fmt.Sprintf("Broadcast Address -> %v", o.BroadcastAddress.String())
}

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptBroadcastAddress) Length() int {
	return len(o.BroadcastAddress.To4())
}
