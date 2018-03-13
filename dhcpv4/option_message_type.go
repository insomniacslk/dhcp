package dhcpv4

import (
	"errors"
	"fmt"
)

// OptMessageType represents the DHCP message type option.
type OptMessageType struct {
	MessageType MessageType
}

// ParseOptMessageType constructs an OptMessageType struct from a sequence of
// bytes and returns it, or an error.
func ParseOptMessageType(data []byte) (*OptMessageType, error) {
	// Should at least have code, length, and message type.
	if len(data) < 3 {
		return nil, errors.New("too short of byte stream")
	}
	code := OptionCode(data[0])
	if code != OptionDHCPMessageType {
		return nil, fmt.Errorf("expected option %v, got %v instead", OptionDHCPMessageType, code)
	}
	length := int(data[1])
	if length != 1 {
		return nil, fmt.Errorf("expected length 1, got %v instead", length)
	}
	messageType := MessageType(data[2])
	return &OptMessageType{MessageType: messageType}, nil
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
	var s string
	switch o.MessageType {
	case MessageTypeDiscover:
		s = "DISCOVER"
	case MessageTypeOffer:
		s = "OFFER"
	case MessageTypeRequest:
		s = "REQUEST"
	case MessageTypeDecline:
		s = "DECLINE"
	case MessageTypeAck:
		s = "ACK"
	case MessageTypeNak:
		s = "NAK"
	case MessageTypeRelease:
		s = "RELEASE"
	case MessageTypeInform:
		s = "INFORM"
	default:
		s = "UNKNOWN"
	}
	return fmt.Sprintf("DHCP Message Type -> %s", s)
}

// Length returns the length of the data portion (excluding option code and byte
// for length, if any).
func (o *OptMessageType) Length() int {
	return 1
}
