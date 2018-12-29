package dhcpv4

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// This option implements the message type option
// https://tools.ietf.org/html/rfc2132

// OptMessageType represents the DHCP message type option.
type OptMessageType struct {
	MessageType MessageType
}

// ParseOptMessageType constructs an OptMessageType struct from a sequence of
// bytes and returns it, or an error.
func ParseOptMessageType(data []byte) (*OptMessageType, error) {
	buf := uio.NewBigEndianBuffer(data)
	return &OptMessageType{MessageType: MessageType(buf.Read8())}, buf.FinError()
}

// Code returns the option code.
func (o *OptMessageType) Code() OptionCode {
	return OptionDHCPMessageType
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptMessageType) ToBytes() []byte {
	return []byte{byte(o.Code()), byte(o.Length()), byte(o.MessageType)}
}

// String returns a human-readable string for this option.
func (o *OptMessageType) String() string {
	return fmt.Sprintf("DHCP Message Type -> %s", o.MessageType.String())
}

// Length returns the length of the data portion (excluding option code and byte
// for length, if any).
func (o *OptMessageType) Length() int {
	return 1
}
