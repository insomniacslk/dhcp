package options

// This module defines the OptRelayMsg structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"fmt"
)

type OptRelayMsg struct {
	relayMessage []byte // FIXME this has to become []DHCPv6
}

func (op *OptRelayMsg) Code() OptionCode {
	return OPTION_RELAY_MSG
}

func (op *OptRelayMsg) ToBytes() []byte {
	return []byte(op.relayMessage)
}

func (op *OptRelayMsg) RelayMessage() []byte {
	return op.relayMessage
}

func (op *OptRelayMsg) SetRelayMessage(relayMessage []byte) {
	op.relayMessage = relayMessage
}

func (op *OptRelayMsg) Length() int {
	return len(op.relayMessage)
}

func (op *OptRelayMsg) String() string {
	return fmt.Sprintf("OptRelayMsg{relaymsg=%v}", op.relayMessage)
}

// build an OptRelayMsg structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptRelayMsg(data []byte) (*OptRelayMsg, error) {
	opt := OptRelayMsg{}
	opt.relayMessage = []byte(data)
	return &opt, nil
}
