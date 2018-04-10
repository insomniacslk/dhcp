package dhcpv6

import (
	"fmt"
	"net"
)

const RelayHeaderSize = 34

type DHCPv6Relay struct {
	messageType MessageType
	hopCount    uint8
	linkAddr    net.IP
	peerAddr    net.IP
	options     []Option
}

func (r *DHCPv6Relay) Type() MessageType {
	return r.messageType
}

func (r *DHCPv6Relay) MessageTypeToString() string {
	return MessageTypeToString(r.messageType)
}

func (r *DHCPv6Relay) String() string {
	ret := fmt.Sprintf(
		"DHCPv6Relay(messageType=%v hopcount=%v, linkaddr=%v, peeraddr=%v, %d options)",
		r.MessageTypeToString(), r.hopCount, r.linkAddr, r.peerAddr, len(r.options),
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
		r.MessageTypeToString(),
		r.hopCount,
		r.linkAddr,
		r.peerAddr,
		r.options,
	)
	return ret
}

func (r *DHCPv6Relay) ToBytes() []byte {
	ret := make([]byte, RelayHeaderSize)
	ret[0] = byte(r.messageType)
	ret[1] = byte(r.hopCount)
	copy(ret[2:18], r.linkAddr)
	copy(ret[18:34], r.peerAddr)
	for _, opt := range r.options {
		ret = append(ret, opt.ToBytes()...)
	}

	return ret
}

func (r *DHCPv6Relay) MessageType() MessageType {
	return r.messageType
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

func (r *DHCPv6Relay) Length() int {
	mLen := RelayHeaderSize
	for _, opt := range r.options {
		mLen += opt.Length() + 4 // +4 for opt code and opt len
	}
	return mLen
}

func (r *DHCPv6Relay) Options() []Option {
	return r.options
}
func (r *DHCPv6Relay) GetOption(code OptionCode) []Option {
	return getOptions(r.options, code, false)
}

func (r *DHCPv6Relay) GetOneOption(code OptionCode) Option {
	return getOption(r.options, code)
}

func (r *DHCPv6Relay) SetOptions(options []Option) {
	r.options = options
}

func (r *DHCPv6Relay) AddOption(option Option) {
	r.options = append(r.options, option)
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

// GetInnerPeerAddr returns the peer address in the inner most relay info
// header, this is typically the IP address of the client making the request.
func (d *DHCPv6Relay) GetInnerPeerAddr() (net.IP, error) {
	var (
		p   DHCPv6
		err error
	)
	p = d
	hops := d.HopCount()
	addr := d.PeerAddr()
	for i := 0; i < int(hops); i++ {
		p, err = DecapsulateRelay(p)
		if err != nil {
			return nil, err
		}
		if p.IsRelay() {
			addr = p.(*DHCPv6Relay).PeerAddr()
		} else {
			return nil, fmt.Errorf("Wrong Hop count")
		}
	}
	return addr, nil
}
