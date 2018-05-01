package dhcpv6

import (
	"fmt"
	"net"
)

type DHCPv6 interface {
	Type() MessageType
	ToBytes() []byte
	Options() []Option
	String() string
	Summary() string
	Length() int
	IsRelay() bool
	GetOption(code OptionCode) []Option
	GetOneOption(code OptionCode) Option
	SetOptions(options []Option)
	AddOption(Option)
	UpdateOption(Option)
}

// Modifier defines the signature for functions that can modify DHCPv6
// structures. This is used to simplify packet manipulation
type Modifier func(d DHCPv6) DHCPv6

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

// NewMessage creates a new DHCPv6 message with default options
func NewMessage() (DHCPv6, error) {
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

func getOptions(options []Option, code OptionCode, onlyFirst bool) []Option {
	var ret []Option
	for _, opt := range options {
		if opt.Code() == code {
			ret = append(ret, opt)
			if onlyFirst {
				break
			}
		}
	}
	return ret
}

func getOption(options []Option, code OptionCode) Option {
	opts := getOptions(options, code, true)
	if opts == nil {
		return nil
	}
	return opts[0]
}

// DecapsulateRelay extracts the content of a relay message. It does not recurse
// if there are nested relay messages. Returns the original packet if is not not
// a relay message
func DecapsulateRelay(l DHCPv6) (DHCPv6, error) {
	if !l.IsRelay() {
		return l, nil
	}
	opt := l.GetOneOption(OPTION_RELAY_MSG)
	if opt == nil {
		return nil, fmt.Errorf("No OptRelayMsg found")
	}
	relayOpt := opt.(*OptRelayMsg)
	if relayOpt.RelayMessage() == nil {
		return nil, fmt.Errorf("Relay message cannot be nil")
	}
	return relayOpt.RelayMessage(), nil
}

// EncapsulateRelay creates a DHCPv6Relay message containing the passed DHCPv6
// message as payload. The passed message type must be  either RELAY_FORW or
// RELAY_REPL
func EncapsulateRelay(d DHCPv6, mType MessageType, linkAddr, peerAddr net.IP) (DHCPv6, error) {
	if mType != RELAY_FORW && mType != RELAY_REPL {
		return nil, fmt.Errorf("Message type must be either RELAY_FORW or RELAY_REPL")
	}
	outer := DHCPv6Relay{
		messageType: mType,
		linkAddr:    linkAddr,
		peerAddr:    peerAddr,
	}
	if d.IsRelay() {
		relay := d.(*DHCPv6Relay)
		outer.hopCount = relay.hopCount + 1
	} else {
		outer.hopCount = 0
	}
	orm := OptRelayMsg{relayMessage: d}
	outer.AddOption(&orm)
	return &outer, nil
}
