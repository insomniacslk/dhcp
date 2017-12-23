package dhcpv6

// This module defines the OptRelayMsg structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

type OptRelayMsg struct {
	relayMessage DHCPv6
}

func (op *OptRelayMsg) Code() OptionCode {
	return OPTION_RELAY_MSG
}

func (op *OptRelayMsg) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_RELAY_MSG))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	buf = append(buf, op.relayMessage.ToBytes()...)
	return buf
}

func (op *OptRelayMsg) RelayMessage() DHCPv6 {
	return op.relayMessage
}

func (op *OptRelayMsg) SetRelayMessage(relayMessage DHCPv6) {
	op.relayMessage = relayMessage
}

func (op *OptRelayMsg) Length() int {
	return op.relayMessage.Length()
}

func (op *OptRelayMsg) String() string {
	return fmt.Sprintf("OptRelayMsg{relaymsg=%v}", op.relayMessage)
}

// build an OptRelayMsg structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptRelayMsg(data []byte) (*OptRelayMsg, error) {
	var err error
	opt := OptRelayMsg{}
	opt.relayMessage, err = FromBytes(data)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
