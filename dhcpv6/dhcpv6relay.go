package dhcpv6

import (
	"errors"
	"fmt"
	"net"

	"github.com/u-root/u-root/pkg/uio"
)

const RelayHeaderSize = 34

// RelayMessage is a DHCPv6 relay agent message as defined by RFC 3315 Section
// 7.
type RelayMessage struct {
	messageType MessageType
	hopCount    uint8
	linkAddr    net.IP
	peerAddr    net.IP
	options     Options
}

// Type is this relay message's types.
func (r *RelayMessage) Type() MessageType {
	return r.messageType
}

// String prints a short human-readable relay message.
func (r *RelayMessage) String() string {
	ret := fmt.Sprintf(
		"RelayMessage(messageType=%v hopcount=%v, linkaddr=%v, peeraddr=%v, %d options)",
		r.Type().String(), r.hopCount, r.linkAddr, r.peerAddr, len(r.options),
	)
	return ret
}

// Summary prints all options associated with this relay message.
func (r *RelayMessage) Summary() string {
	ret := fmt.Sprintf(
		"RelayMessage\n"+
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
// RFC 3315, Section 7.
func (r *RelayMessage) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(make([]byte, 0, RelayHeaderSize))
	buf.Write8(byte(r.messageType))
	buf.Write8(byte(r.hopCount))
	buf.WriteBytes(r.linkAddr.To16())
	buf.WriteBytes(r.peerAddr.To16())
	buf.WriteBytes(r.options.ToBytes())
	return buf.Data()
}

// SetMessageType sets the message type of this relay message.
func (r *RelayMessage) SetMessageType(messageType MessageType) {
	// not enforcing if message type is not a RELAY_FORW or a RELAY_REPL message
	r.messageType = messageType
}

// HopCount returns the hop count.
func (r *RelayMessage) HopCount() uint8 {
	return r.hopCount
}

// SetHopCount sets the hop count.
func (r *RelayMessage) SetHopCount(hopCount uint8) {
	r.hopCount = hopCount
}

// LinkAddr returns the link address for this relay message.
func (r *RelayMessage) LinkAddr() net.IP {
	return r.linkAddr
}

// SetLinkAddr sets the link address.
func (r *RelayMessage) SetLinkAddr(linkAddr net.IP) {
	r.linkAddr = linkAddr
}

// PeerAddr returns the peer address for this relay message.
func (r *RelayMessage) PeerAddr() net.IP {
	return r.peerAddr
}

// SetPeerAddr sets the peer address.
func (r *RelayMessage) SetPeerAddr(peerAddr net.IP) {
	r.peerAddr = peerAddr
}

// Options returns the current set of options associated with this message.
func (r *RelayMessage) Options() []Option {
	return r.options
}

// GetOption returns the options associated with the code.
func (r *RelayMessage) GetOption(code OptionCode) []Option {
	return r.options.Get(code)
}

// GetOneOption returns the first associated option with the code from this
// message.
func (r *RelayMessage) GetOneOption(code OptionCode) Option {
	return r.options.GetOne(code)
}

// SetOptions replaces this message's options.
func (r *RelayMessage) SetOptions(options []Option) {
	r.options = options
}

// AddOption adds an option to this message.
func (r *RelayMessage) AddOption(option Option) {
	r.options.Add(option)
}

// UpdateOption replaces the first option of the same type as the specified one.
func (r *RelayMessage) UpdateOption(option Option) {
	r.options.Update(option)
}

// IsRelay returns whether this is a relay message or not.
func (r *RelayMessage) IsRelay() bool {
	return true
}

// GetInnerMessage recurses into a relay message and extract and return the
// inner Message.  Return nil if none found (e.g. not a relay message).
func (r *RelayMessage) GetInnerMessage() (DHCPv6, error) {
	var (
		p   DHCPv6
		err error
	)
	p = r
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
	relay, ok := relayForw.(*RelayMessage)
	if !ok {
		return nil, errors.New("Not a RelayMessage")
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
			relay = decap.(*RelayMessage)
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
