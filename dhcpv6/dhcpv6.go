package dhcpv6

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/insomniacslk/dhcp/iana"
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
	if messageType == MessageTypeRelayForward || messageType == MessageTypeRelayReply {
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
		// TODO fail if no OptRelayMessage is present
		if err := d.options.FromBytes(data[34:]); err != nil {
			return nil, err
		}
		return &d, nil
	} else {
		d := DHCPv6Message{
			messageType: messageType,
		}
		copy(d.transactionID[:], data[1:4])
		if err := d.options.FromBytes(data[4:]); err != nil {
			return nil, err
		}
		return &d, nil
	}
}

// NewMessage creates a new DHCPv6 message with default options
func NewMessage(modifiers ...Modifier) (DHCPv6, error) {
	tid, err := GenerateTransactionID()
	if err != nil {
		return nil, err
	}
	msg := DHCPv6Message{
		messageType:   MessageTypeSolicit,
		transactionID: tid,
	}
	// apply modifiers
	d := DHCPv6(&msg)
	for _, mod := range modifiers {
		d = mod(d)
	}
	return d, nil
}

// DecapsulateRelay extracts the content of a relay message. It does not recurse
// if there are nested relay messages. Returns the original packet if is not not
// a relay message
func DecapsulateRelay(l DHCPv6) (DHCPv6, error) {
	if !l.IsRelay() {
		return l, nil
	}
	opt := l.GetOneOption(OptionRelayMsg)
	if opt == nil {
		return nil, fmt.Errorf("No OptRelayMsg found")
	}
	relayOpt := opt.(*OptRelayMsg)
	if relayOpt.RelayMessage() == nil {
		return nil, fmt.Errorf("Relay message cannot be nil")
	}
	return relayOpt.RelayMessage(), nil
}

// DecapsulateRelayIndex extracts the content of a relay message. It takes an
// integer as index (e.g. if 0 return the outermost relay, 1 returns the
// second, etc, and -1 returns the last). Returns the original packet if
// it is not not a relay message.
func DecapsulateRelayIndex(l DHCPv6, index int) (DHCPv6, error) {
	if !l.IsRelay() {
		return l, nil
	}
	if index < -1 {
		return nil, fmt.Errorf("Invalid index: %d", index)
	} else if index == -1 {
		for {
			d, err := DecapsulateRelay(l)
			if err != nil {
				return nil, err
			}
			if !d.IsRelay() {
				return l, nil
			}
			l = d
		}
	}
	for i := 0; i <= index; i++ {
		d, err := DecapsulateRelay(l)
		if err != nil {
			return nil, err
		}
		l = d
	}
	return l, nil
}

// EncapsulateRelay creates a DHCPv6Relay message containing the passed DHCPv6
// message as payload. The passed message type must be  either RELAY_FORW or
// RELAY_REPL
func EncapsulateRelay(d DHCPv6, mType MessageType, linkAddr, peerAddr net.IP) (DHCPv6, error) {
	if mType != MessageTypeRelayForward && mType != MessageTypeRelayReply {
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

// IsUsingUEFI function takes a DHCPv6 message and returns true if
// the machine trying to netboot is using UEFI of false if it is not.
func IsUsingUEFI(msg DHCPv6) bool {
	// RFC 4578 says:
	// As of the writing of this document, the following pre-boot
	//    architecture types have been requested.
	//             Type   Architecture Name
	//             ----   -----------------
	//               0    Intel x86PC
	//               1    NEC/PC98
	//               2    EFI Itanium
	//               3    DEC Alpha
	//               4    Arc x86
	//               5    Intel Lean Client
	//               6    EFI IA32
	//               7    EFI BC
	//               8    EFI Xscale
	//               9    EFI x86-64
	if opt := msg.GetOneOption(OptionClientArchType); opt != nil {
		optat := opt.(*OptClientArchType)
		for _, at := range optat.ArchTypes {
			// TODO investigate if other types are appropriate
			if at == iana.EFI_BC || at == iana.EFI_X86_64 {
				return true
			}
		}
	}
	if opt := msg.GetOneOption(OptionUserClass); opt != nil {
		optuc := opt.(*OptUserClass)
		for _, uc := range optuc.UserClasses {
			if strings.Contains(string(uc), "EFI") {
				return true
			}
		}
	}
	return false
}

// GetTransactionID returns a transactionID of a message or its inner message
// in case of relay
func GetTransactionID(packet DHCPv6) (TransactionID, error) {
	if message, ok := packet.(*DHCPv6Message); ok {
		return message.TransactionID(), nil
	}
	if relay, ok := packet.(*DHCPv6Relay); ok {
		message, err := relay.GetInnerMessage()
		if err != nil {
			return TransactionID{0, 0, 0}, err
		}
		return GetTransactionID(message)
	}
	return TransactionID{0, 0, 0}, errors.New("Invalid DHCPv6 packet")
}
