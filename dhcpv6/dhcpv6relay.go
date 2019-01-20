package dhcpv6

import (
	"errors"
	"fmt"
	"net"

	"github.com/u-root/u-root/pkg/uio"
)

const RelayHeaderSize = 34

type DHCPv6Relay struct {
	messageType MessageType
	hopCount    uint8
	linkAddr    net.IP
	peerAddr    net.IP
	options     Options
}

func (r *DHCPv6Relay) Type() MessageType {
	return r.messageType
}

func (r *DHCPv6Relay) MessageTypeToString() string {
	return r.messageType.String()
}

func (r *DHCPv6Relay) String() string {
	ret := fmt.Sprintf(
		"DHCPv6Relay(messageType=%v hopcount=%v, linkaddr=%v, peeraddr=%v, %d options)",
		r.Type().String(), r.hopCount, r.linkAddr, r.peerAddr, len(r.options),
	)
	return ret
}

func (r *DHCPv6Relay) Summary() string {
	ret := fmt.Sprintf(
		"DHCPv6Relay\n"+
			"  messageType=%v\n"+
			"  hopcount=%v\n"+
			"  linkaddr=%v\n"+
			"  peeraddr=%v\n"+
			"  options=%v\n",
		r.Type().String(),
		r.hopCount,
		r.linkAddr,
		r.peerAddr,
		r.options,
	)
	return ret
}

// ToBytes returns the serialized version of this relay message as defined by
// RFC 3315, Section 6.
func (r *DHCPv6Relay) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(make([]byte, 0, RelayHeaderSize))
	buf.Write8(byte(r.messageType))
	buf.Write8(byte(r.hopCount))
	buf.WriteBytes(r.linkAddr.To16())
	buf.WriteBytes(r.peerAddr.To16())
	buf.WriteBytes(r.options.ToBytes())
	return buf.Data()
}

func (r *DHCPv6Relay) SetMessageType(messageType MessageType) {
	// not enforcing if message type is not a RELAY_FORW or a RELAY_REPL message
	r.messageType = messageType
}

func (r *DHCPv6Relay) HopCount() uint8 {
	return r.hopCount
}

func (r *DHCPv6Relay) SetHopCount(hopCount uint8) {
	r.hopCount = hopCount
}

func (r *DHCPv6Relay) LinkAddr() net.IP {
	return r.linkAddr
}

func (r *DHCPv6Relay) SetLinkAddr(linkAddr net.IP) {
	r.linkAddr = linkAddr
}

func (r *DHCPv6Relay) PeerAddr() net.IP {
	return r.peerAddr
}

func (r *DHCPv6Relay) SetPeerAddr(peerAddr net.IP) {
	r.peerAddr = peerAddr
}

func (r *DHCPv6Relay) Options() []Option {
	return r.options
}

func (r *DHCPv6Relay) GetOption(code OptionCode) []Option {
	return r.options.Get(code)
}

func (r *DHCPv6Relay) GetOneOption(code OptionCode) Option {
	return r.options.GetOne(code)
}

func (r *DHCPv6Relay) SetOptions(options []Option) {
	r.options = options
}

func (r *DHCPv6Relay) AddOption(option Option) {
	r.options.Add(option)
}

// UpdateOption replaces the first option of the same type as the specified one.
func (r *DHCPv6Relay) UpdateOption(option Option) {
	r.options.Update(option)
}

func (r *DHCPv6Relay) IsRelay() bool {
	return true
}

// Recurse into a relay message and extract and return the inner DHCPv6Message.
// Return nil if none found (e.g. not a relay message).
func (d *DHCPv6Relay) GetInnerMessage() (DHCPv6, error) {
	var (
		p   DHCPv6
		err error
	)
	p = d
	for {
		if !p.IsRelay() {
			return p, nil
		}
		p, err = DecapsulateRelay(p)
		if err != nil {
			return nil, err
		}
	}
}

// NewRelayReplFromRelayForw creates a MessageTypeRelayReply based on a
// MessageTypeRelayForward and replaces the inner message with the passed
// DHCPv6 message. It copies the OptionInterfaceID and OptionRemoteID if the
// options are present in the Relay packet.
func NewRelayReplFromRelayForw(relayForw, msg DHCPv6) (DHCPv6, error) {
	var (
		err                error
		linkAddr, peerAddr []net.IP
		optiid             []Option
		optrid             []Option
	)
	if relayForw == nil {
		return nil, errors.New("Relay message cannot be nil")
	}
	relay, ok := relayForw.(*DHCPv6Relay)
	if !ok {
		return nil, errors.New("Not a DHCPv6Relay")
	}
	if relay.Type() != MessageTypeRelayForward {
		return nil, errors.New("The passed packet is not of type MessageTypeRelayForward")
	}
	if msg == nil {
		return nil, errors.New("The passed message cannot be nil")
	}
	if msg.IsRelay() {
		return nil, errors.New("The passed message cannot be a relay")
	}
	for {
		linkAddr = append(linkAddr, relay.LinkAddr())
		peerAddr = append(peerAddr, relay.PeerAddr())
		optiid = append(optiid, relay.GetOneOption(OptionInterfaceID))
		optrid = append(optrid, relay.GetOneOption(OptionRemoteID))
		decap, err := DecapsulateRelay(relay)
		if err != nil {
			return nil, err
		}
		if decap.IsRelay() {
			relay = decap.(*DHCPv6Relay)
		} else {
			break
		}
	}
	for i := len(linkAddr) - 1; i >= 0; i-- {
		msg, err = EncapsulateRelay(msg, MessageTypeRelayReply, linkAddr[i], peerAddr[i])
		if err != nil {
			return nil, err
		}
		if opt := optiid[i]; opt != nil {
			msg.AddOption(opt)
		}
		if opt := optrid[i]; opt != nil {
			msg.AddOption(opt)
		}
	}
	return msg, nil
}
