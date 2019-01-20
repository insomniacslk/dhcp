package dhcpv6

import (
	"fmt"
)

// MessageType represents the kind of DHCPv6 message.
type MessageType uint8

// The DHCPv6 message types defined per RFC 3315, Section 5.3.
const (
	// MessageTypeNone is used internally and is not part of the RFC.
	MessageTypeNone               MessageType = 0
	MessageTypeSolicit            MessageType = 1
	MessageTypeAdvertise          MessageType = 2
	MessageTypeRequest            MessageType = 3
	MessageTypeConfirm            MessageType = 4
	MessageTypeRenew              MessageType = 5
	MessageTypeRebind             MessageType = 6
	MessageTypeReply              MessageType = 7
	MessageTypeRelease            MessageType = 8
	MessageTypeDecline            MessageType = 9
	MessageTypeReconfigure        MessageType = 10
	MessageTypeInformationRequest MessageType = 11
	MessageTypeRelayForward       MessageType = 12
	MessageTypeRelayReply         MessageType = 13
	MessageTypeLeaseQuery         MessageType = 14
	MessageTypeLeaseQueryReply    MessageType = 15
	MessageTypeLeaseQueryDone     MessageType = 16
	MessageTypeLeaseQueryData     MessageType = 17
)

// String prints the message type name.
func (m MessageType) String() string {
	if s, ok := messageTypeToStringMap[m]; ok {
		return s
	}
	return fmt.Sprintf("unknown (%d)", m)
}

// messageTypeToStringMap contains the mapping of MessageTypes to
// human-readable strings.
var messageTypeToStringMap = map[MessageType]string{
	MessageTypeSolicit:            "SOLICIT",
	MessageTypeAdvertise:          "ADVERTISE",
	MessageTypeRequest:            "REQUEST",
	MessageTypeConfirm:            "CONFIRM",
	MessageTypeRenew:              "RENEW",
	MessageTypeRebind:             "REBIND",
	MessageTypeReply:              "REPLY",
	MessageTypeRelease:            "RELEASE",
	MessageTypeDecline:            "DECLINE",
	MessageTypeReconfigure:        "RECONFIGURE",
	MessageTypeInformationRequest: "INFORMATION-REQUEST",
	MessageTypeRelayForward:       "RELAY-FORW",
	MessageTypeRelayReply:         "RELAY-REPL",
	MessageTypeLeaseQuery:         "LEASEQUERY",
	MessageTypeLeaseQueryReply:    "LEASEQUERY-REPLY",
	MessageTypeLeaseQueryDone:     "LEASEQUERY-DONE",
	MessageTypeLeaseQueryData:     "LEASEQUERY-DATA",
}
