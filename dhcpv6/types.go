package dhcpv6

// from http://www.networksorcery.com/enp/protocol/dhcpv6.htm

type MessageType uint8

const (
	SOLICIT             MessageType = 1
	ADVERTISE           MessageType = 2
	REQUEST             MessageType = 3
	CONFIRM             MessageType = 4
	RENEW               MessageType = 5
	REBIND              MessageType = 6
	REPLY               MessageType = 7
	RELEASE             MessageType = 8
	DECLINE             MessageType = 9
	RECONFIGURE         MessageType = 10
	INFORMATION_REQUEST MessageType = 11
	RELAY_FORW          MessageType = 12
	RELAY_REPL          MessageType = 13
	LEASEQUERY          MessageType = 14
	LEASEQUERY_REPLY    MessageType = 15
	LEASEQUERY_DONE     MessageType = 16
	LEASEQUERY_DATA     MessageType = 17
)

func MessageTypeToString(t MessageType) string {
	if m := MessageTypeToStringMap[t]; m != "" {
		return m
	}
	return "Unknown"
}

var MessageTypeToStringMap = map[MessageType]string{
	SOLICIT:             "SOLICIT",
	ADVERTISE:           "ADVERTISE",
	REQUEST:             "REQUEST",
	CONFIRM:             "CONFIRM",
	RENEW:               "RENEW",
	REBIND:              "REBIND",
	REPLY:               "REPLY",
	RELEASE:             "RELEASE",
	DECLINE:             "DECLINE",
	RECONFIGURE:         "RECONFIGURE",
	INFORMATION_REQUEST: "INFORMATION-REQUEST",
	RELAY_FORW:          "RELAY-FORW",
	RELAY_REPL:          "RELAY-REPL",
	LEASEQUERY:          "LEASEQUERY",
	LEASEQUERY_REPLY:    "LEASEQUERY-REPLY",
	LEASEQUERY_DONE:     "LEASEQUERY-DONE",
	LEASEQUERY_DATA:     "LEASEQUERY-DATA",
}
