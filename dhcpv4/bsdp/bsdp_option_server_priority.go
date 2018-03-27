package bsdp

import (
	"encoding/binary"
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// This option implements the server identifier option
// https://tools.ietf.org/html/rfc2132

// OptServerPriority represents an option encapsulating the server priority.
type OptServerPriority struct {
	Priority int
}

// ParseOptServerPriority returns a new OptServerPriority from a byte stream, or
// error if any.
func ParseOptServerPriority(data []byte) (*OptServerPriority, error) {
	if len(data) < 4 {
		return nil, dhcpv4.ErrShortByteStream
	}
	code := dhcpv4.OptionCode(data[0])
	if code != OptionServerPriority {
		return nil, fmt.Errorf("expected code %v, got %v", OptionServerPriority, code)
	}
	length := int(data[1])
	if length != 2 {
		return nil, fmt.Errorf("unexpected length: expected 2, got %v", length)
	}
	return &OptServerPriority{Priority: int(binary.BigEndian.Uint16(data[2:4]))}, nil
}

// Code returns the option code.
func (o *OptServerPriority) Code() dhcpv4.OptionCode {
	return OptionServerPriority
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptServerPriority) ToBytes() []byte {
	serialized := make([]byte, 2)
	binary.BigEndian.PutUint16(serialized, uint16(o.Priority))
	return append([]byte{byte(o.Code()), byte(o.Length())}, serialized...)
}

// String returns a human-readable string.
func (o *OptServerPriority) String() string {
	return fmt.Sprintf("BSDP Server Priority -> %v", o.Priority)
}

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptServerPriority) Length() int {
	return 2
}
