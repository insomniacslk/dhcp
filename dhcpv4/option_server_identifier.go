package dhcpv4

import (
	"fmt"
	"net"
)

// This option implements the server identifier option
// https://tools.ietf.org/html/rfc2132

// OptServerIdentifier represents an option encapsulating the server identifier.
type OptServerIdentifier struct {
	ServerID net.IP
}

// ParseOptServerIdentifier returns a new OptServerIdentifier from a byte
// stream, or error if any.
func ParseOptServerIdentifier(data []byte) (*OptServerIdentifier, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf("unexepcted length: expected 4, got %v", len(data))
	}
	return &OptServerIdentifier{ServerID: net.IP(data)}, nil
}

// Code returns the option code.
func (o *OptServerIdentifier) Code() OptionCode {
	return OptionServerIdentifier
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptServerIdentifier) ToBytes() []byte {
	ret := []byte{byte(o.Code()), byte(o.Length())}
	return append(ret, o.ServerID.To4()...)
}

// String returns a human-readable string.
func (o *OptServerIdentifier) String() string {
	return fmt.Sprintf("Server Identifier -> %v", o.ServerID.String())
}

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptServerIdentifier) Length() int {
	return len(o.ServerID.To4())
}
