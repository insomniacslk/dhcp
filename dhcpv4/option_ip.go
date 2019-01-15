package dhcpv4

import (
	"fmt"
	"net"

	"github.com/u-root/u-root/pkg/uio"
)

// OptBroadcastAddress implements the broadcast address option described in RFC
// 2132, Section 5.3.
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
	return []byte(o.BroadcastAddress.To4())
}

// String returns a human-readable string.
func (o *OptBroadcastAddress) String() string {
	return fmt.Sprintf("Broadcast Address -> %v", o.BroadcastAddress.String())
}

// OptRequestedIPAddress implements the requested IP address option described
// by RFC 2132, Section 9.1.
type OptRequestedIPAddress struct {
	RequestedAddr net.IP
}

// ParseOptRequestedIPAddress returns a new OptRequestedIPAddress from a byte
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

// OptServerIdentifier implements the server identifier option described by RFC
// 2132, Section 9.7.
type OptServerIdentifier struct {
	ServerID net.IP
}

// ParseOptServerIdentifier returns a new OptServerIdentifier from a byte
// stream, or error if any.
func ParseOptServerIdentifier(data []byte) (*OptServerIdentifier, error) {
	buf := uio.NewBigEndianBuffer(data)
	return &OptServerIdentifier{ServerID: net.IP(buf.CopyN(net.IPv4len))}, buf.FinError()
}

// Code returns the option code.
func (o *OptServerIdentifier) Code() OptionCode {
	return OptionServerIdentifier
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptServerIdentifier) ToBytes() []byte {
	return o.ServerID.To4()
}

// String returns a human-readable string.
func (o *OptServerIdentifier) String() string {
	return fmt.Sprintf("Server Identifier -> %v", o.ServerID.String())
}
