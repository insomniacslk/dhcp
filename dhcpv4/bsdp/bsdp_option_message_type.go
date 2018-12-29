package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// MessageType represents the different BSDP message types.
//
// Implements the BSDP option message type. Can be one of LIST, SELECT, or
// FAILED.
type MessageType byte

// BSDP Message types - e.g. LIST, SELECT, FAILED
const (
	MessageTypeList   MessageType = 1
	MessageTypeSelect MessageType = 2
	MessageTypeFailed MessageType = 3
)

func (m MessageType) String() string {
	if s, ok := messageTypeToString[m]; ok {
		return s
	}
	return "Unknown"
}

// messageTypeToString maps each BSDP message type to a human-readable string.
var messageTypeToString = map[MessageType]string{
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
	buf := uio.NewBigEndianBuffer(data)
	return &OptMessageType{Type: MessageType(buf.Read8())}, buf.FinError()
}

// Code returns the option code.
func (o *OptMessageType) Code() dhcpv4.OptionCode {
	return OptionMessageType
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptMessageType) ToBytes() []byte {
	return []byte{byte(o.Type)}
}

// String returns a human-readable string for this option.
func (o *OptMessageType) String() string {
	return fmt.Sprintf("BSDP Message Type -> %s", o.Type.String())
}
