package dhcpv6

import (
	"fmt"
)

type DHCPv6 interface {
	Type() MessageType
	ToBytes() []byte
	String() string
	Summary() string
	Length() int
}

func FromBytes(data []byte) (DHCPv6, error) {
	var (
		isRelay     = false
		headerSize  int
		messageType = MessageType(data[0])
	)
	if messageType == RELAY_FORW || messageType == RELAY_REPL {
		isRelay = true
	}
	if isRelay {
		headerSize = RelayHeaderSize
	} else {
		headerSize = MessageHeaderSize
	}
	if len(data) < headerSize {
		return nil, fmt.Errorf("Invalid header size: shorter than %v bytes", headerSize)
	}
	if isRelay {
		var (
			linkAddr, peerAddr []byte
		)
		d := DHCPv6Relay{
			messageType: messageType,
			hopCount:    uint8(data[1]),
		}
		linkAddr = append(linkAddr, data[2:18]...)
		d.linkAddr = linkAddr
		peerAddr = append(peerAddr, data[18:34]...)
		d.peerAddr = peerAddr
		options, err := OptionsFromBytes(data[34:])
		if err != nil {
			return nil, err
		}
		// TODO fail if no OptRelayMessage is present
		d.options = options
		return &d, nil
	} else {
		tid, err := BytesToTransactionID(data[1:4])
		if err != nil {
			return nil, err
		}
		d := DHCPv6Message{
			messageType:   messageType,
			transactionID: *tid,
		}
		options, err := OptionsFromBytes(data[4:])
		if err != nil {
			return nil, err
		}
		d.options = options
		return &d, nil
	}
}

func NewMessage() (*DHCPv6Message, error) {
	tid, err := GenerateTransactionID()
	if err != nil {
		return nil, err
	}
	d := DHCPv6Message{
		messageType:   SOLICIT,
		transactionID: *tid,
	}
	return &d, nil
}
