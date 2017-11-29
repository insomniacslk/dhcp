package dhcpv6

// from http://www.networksorcery.com/enp/protocol/dhcpv6.htm

type MessageType uint8

const (
	_ MessageType = iota // skip 0
	SOLICIT
	ADVERTISE
	REQUEST
	CONFIRM
	RENEW
	REBIND
	REPLY
	RELEASE
	DECLINE
	RECONFIGURE
	INFORMATION_REQUEST
	RELAY_FORW
	RELAY_REPL
	LEASEQUERY
	LEASEQUERY_REPLY
	LEASEQUERY_DONE
	LEASEQUERY_DATA
)

var MessageToString = map[MessageType]string{
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
