package dhcpv6

import (
	"fmt"
	"net"
	"time"
)

// Client constants
const (
	DefaultWriteTimeout       = 3 * time.Second // time to wait for write calls
	DefaultReadTimeout        = 3 * time.Second // time to wait for read calls
	DefaultInterfaceUpTimeout = 3 * time.Second // time to wait before a network interface goes up
	maxUDPReceivedPacketSize  = 8192            // arbitrary size. Theoretically could be up to 65kb
)

// Broadcast destination IP addresses as defined by RFC 3315
var (
	AllDHCPRelayAgentsAndServers = net.ParseIP("ff02::1:2")
	AllDHCPServers               = net.ParseIP("ff05::1:3")
)

// Client implements a DHCPv6 client
type Client struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LocalAddr    net.Addr
	RemoteAddr   net.Addr
}

// NewClient returns a Client with default settings
func NewClient() *Client {
	return &Client{
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
	}
}

// Exchange executes a 4-way DHCPv6 request (SOLICIT, ADVERTISE, REQUEST,
// REPLY). If the SOLICIT packet is nil, defaults are used. The modifiers will
// be applied to the Request packet. A common use is to make sure that the
// Request packet has the right options, see modifiers.go
func (c *Client) Exchange(ifname string, solicit DHCPv6, modifiers ...Modifier) ([]DHCPv6, error) {
	conversation := make([]DHCPv6, 0)
	var err error

	// Solicit
	if solicit == nil {
		solicit, err = NewSolicitForInterface(ifname)
		if err != nil {
			return conversation, err
		}
	}
	solicit, advertise, err := c.Solicit(ifname, solicit, modifiers...)
	conversation = append(conversation, solicit)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, advertise)

  // Decapsulate advertise if it's relayed before passing it to Request
  if advertise.IsRelay() {
    advertiseRelay := advertise.(*DHCPv6Relay)
    advertise, err = advertiseRelay.GetInnerMessage()
    if err != nil {
      return conversation, err
    }
  }
	request, reply, err := c.Request(ifname, advertise, nil, modifiers...)
	if request != nil {
		conversation = append(conversation, request)
	}
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, reply)
	return conversation, nil
}

func (c *Client) sendReceive(ifname string, packet DHCPv6, expectedType MessageType) (DHCPv6, error) {
	if packet == nil {
		return nil, fmt.Errorf("Packet to send cannot be nil")
	}
	if expectedType == MSGTYPE_NONE {
		// infer the expected type from the packet being sent
		if packet.Type() == SOLICIT {
			expectedType = ADVERTISE
		} else if packet.Type() == REQUEST {
			expectedType = REPLY
		} else if packet.Type() == RELAY_FORW {
			expectedType = RELAY_REPL
		} else if packet.Type() == LEASEQUERY {
			expectedType = LEASEQUERY_REPLY
		} // and probably more
	}
	// if no LocalAddr is specified, get the interface's link-local address
	var laddr net.UDPAddr
	if c.LocalAddr == nil {
		llAddr, err := GetLinkLocalAddr(ifname)
		if err != nil {
			return nil, err
		}
		laddr = net.UDPAddr{IP: *llAddr, Port: DefaultClientPort, Zone: ifname}
	} else {
		if addr, ok := c.LocalAddr.(*net.UDPAddr); ok {
			laddr = *addr
		} else {
			return nil, fmt.Errorf("Invalid local address: not a net.UDPAddr: %v", c.LocalAddr)
		}
	}

	// if no RemoteAddr is specified, use AllDHCPRelayAgentsAndServers
	var raddr net.UDPAddr
	if c.RemoteAddr == nil {
		raddr = net.UDPAddr{IP: AllDHCPRelayAgentsAndServers, Port: DefaultServerPort}
	} else {
		if addr, ok := c.RemoteAddr.(*net.UDPAddr); ok {
			raddr = *addr
		} else {
			return nil, fmt.Errorf("Invalid remote address: not a net.UDPAddr: %v", c.RemoteAddr)
		}
	}

	// prepare the socket to listen on for replies
	conn, err := net.ListenUDP("udp6", &laddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// send the packet out
	conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	_, err = conn.WriteTo(packet.ToBytes(), &raddr)
	if err != nil {
		return nil, err
	}

	// wait for a reply
	oobdata := []byte{} // ignoring oob data
	conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	var (
		adv       DHCPv6
		isMessage bool
	)
	defer conn.Close()
	msg, ok := packet.(*DHCPv6Message)
	if ok {
		isMessage = true
	}
	for {
		buf := make([]byte, maxUDPReceivedPacketSize)
		n, _, _, _, err := conn.ReadMsgUDP(buf, oobdata)
		if err != nil {
			return nil, err
		}
		adv, err = FromBytes(buf[:n])
		if err != nil {
			// skip non-DHCP packets
			continue
		}
		if recvMsg, ok := adv.(*DHCPv6Message); ok && isMessage {
			// if a regular message, check the transaction ID first
			// XXX should this unpack relay messages and check the XID of the
			// inner packet too?
			if msg.TransactionID() != recvMsg.TransactionID() {
				// different XID, we don't want this packet for sure
				continue
			}
		}
		if expectedType == MSGTYPE_NONE {
			// just take whatever arrived
			break
		} else if adv.Type() == expectedType {
			break
		}
	}
	return adv, nil
}

// Solicit sends a SOLICIT, return the solicit, an ADVERTISE (if not nil), and
// an error if any
func (c *Client) Solicit(ifname string, solicit DHCPv6, modifiers ...Modifier) (DHCPv6, DHCPv6, error) {
	var err error
	if solicit == nil {
		solicit, err = NewSolicitForInterface(ifname)
		if err != nil {
			return nil, nil, err
		}
	}
	for _, mod := range modifiers {
		solicit = mod(solicit)
	}
	advertise, err := c.sendReceive(ifname, solicit, MSGTYPE_NONE)
	return solicit, advertise, err
}

// Request sends a REQUEST built from an ADVERTISE if no REQUEST is specified.
// It returns the request, a reply if not nil, and an error if any
func (c *Client) Request(ifname string, advertise, request DHCPv6, modifiers ...Modifier) (DHCPv6, DHCPv6, error) {
	if request == nil {
		var err error
		request, err = NewRequestFromAdvertise(advertise)
		if err != nil {
			return nil, nil, err
		}
	}
	for _, mod := range modifiers {
		request = mod(request)
	}
	reply, err := c.sendReceive(ifname, request, MSGTYPE_NONE)
	return request, reply, err
}
