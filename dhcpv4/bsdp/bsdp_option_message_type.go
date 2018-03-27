package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Implements the BSDP option message type. Can be one of LIST, SELECT, or
// FAILED.

// MessageType represents the different BSDP message types.
type MessageType byte

// BSDP Message types - e.g. LIST, SELECT, FAILED
const (
	MessageTypeList   MessageType = 1
	MessageTypeSelect MessageType = 2
	MessageTypeFailed MessageType = 3
)

// MessageTypeToString maps each BSDP message type to a human-readable string.
var MessageTypeToString = map[MessageType]string{
	MessageTypeList:   "LIST",
	MessageTypeSelect: "SELECT",
	MessageTypeFailed: "FAILED",
}

// OptMessageType represents a BSDP message type.
type OptMessageType struct {
	Type MessageType
}

// ParseOptMessageType constructs an OptMessageType struct from a sequence of
// bytes and returns it, or an error.
func ParseOptMessageType(data []byte) (*OptMessageType, error) {
	if len(data) < 3 {
		return nil, dhcpv4.ErrShortByteStream
	}
	code := dhcpv4.OptionCode(data[0])
	if code != OptionMessageType {
		return nil, fmt.Errorf("expected option %v, got %v instead", OptionMessageType, code)
	}
	length := int(data[1])
	if length != 1 {
		return nil, fmt.Errorf("expected length 1, got %d instead", length)
	}
	return &OptMessageType{Type: MessageType(data[2])}, nil
}

// Code returns the option code.
func (o *OptMessageType) Code() dhcpv4.OptionCode {
	return OptionMessageType
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptMessageType) ToBytes() []byte {
	return []byte{byte(o.Code()), 1, byte(o.Type)}
}

// String returns a human-readable string for this option.
func (o *OptMessageType) String() string {
	s, ok := MessageTypeToString[o.Type]
	if !ok {
		s = "UNKNOWN"
	}
	return fmt.Sprintf("BSDP Message Type -> %s", s)
}

// Length returns the length of the data portion of this option.
func (o *OptMessageType) Length() int {
	return 1
}
