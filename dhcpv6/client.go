package dhcpv6

import (
	"fmt"
	"net"
	"time"
)

const (
	DefaultWriteTimeout       = 3 * time.Second // time to wait for write calls
	DefaultReadTimeout        = 3 * time.Second // time to wait for read calls
	DefaultInterfaceUpTimeout = 3 * time.Second // time to wait before a network interface goes up
	maxUDPReceivedPacketSize  = 8192            // arbitrary size. Theoretically could be up to 65kb
)

var AllDHCPRelayAgentsAndServers = net.ParseIP("ff02::1:2")
var AllDHCPServers = net.ParseIP("ff05::1:3")

type Client struct {
	ReadTimeout  *time.Duration
	WriteTimeout *time.Duration
	LocalAddr    net.Addr
	RemoteAddr   net.Addr
}

// Make a stateful DHCPv6 request
func (c *Client) Exchange(ifname string, solicit DHCPv6) ([]DHCPv6, error) {
	conversation := make([]DHCPv6, 0)
	var err error

	// Solicit
	solicit, advertise, err := c.Solicit(ifname, solicit)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, solicit)
	conversation = append(conversation, advertise)

	request, reply, err := c.Request(ifname, advertise, nil)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, request)
	conversation = append(conversation, reply)
	return conversation, nil
}

func (c *Client) sendReceive(ifname string, packet DHCPv6) (DHCPv6, error) {
	if packet == nil {
		return nil, fmt.Errorf("Packet to send cannot be nil")
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

	// set WriteTimeout to DefaultWriteTimeout if no other timeout is specified
	var wtimeout time.Duration
	if c.WriteTimeout == nil {
		wtimeout = DefaultWriteTimeout
	} else {
		wtimeout = *c.WriteTimeout
	}
	conn.SetWriteDeadline(time.Now().Add(wtimeout))

	// send the packet out
	_, err = conn.WriteTo(packet.ToBytes(), &raddr)
	if err != nil {
		return nil, err
	}

	// set ReadTimeout to DefaultReadTimeout if no other timeout is specified
	var rtimeout time.Duration
	if c.ReadTimeout == nil {
		rtimeout = DefaultReadTimeout
	} else {
		rtimeout = *c.ReadTimeout
	}
	conn.SetReadDeadline(time.Now().Add(rtimeout))

	// wait for an ADVERTISE response
	buf := make([]byte, maxUDPReceivedPacketSize)
	oobdata := []byte{} // ignoring oob data
	n, _, _, _, err := conn.ReadMsgUDP(buf, oobdata)
	if err != nil {
		return nil, err
	}
	adv, err := FromBytes(buf[:n])
	if err != nil {
		return nil, err
	}
	return adv, nil
}

// send a SOLICIT, return the solicit, an ADVERTISE (if not nil), and an error if
// any
func (c *Client) Solicit(ifname string, solicit DHCPv6) (DHCPv6, DHCPv6, error) {
	var err error
	if solicit == nil {
		solicit, err = NewSolicitForInterface(ifname)
		if err != nil {
			return nil, nil, err
		}
	}
	advertise, err := c.sendReceive(ifname, solicit)
	return solicit, advertise, err
}

// send a REQUEST built from an ADVERTISE if no REQUEST is specified, return a
// the request, a reply if not nil, and an error if any
func (c *Client) Request(ifname string, advertise, request DHCPv6) (DHCPv6, DHCPv6, error) {
	if request == nil {
		var err error
		request, err = NewRequestFromAdvertise(advertise)
		if err != nil {
			return nil, nil, err
		}
	}
	reply, err := c.sendReceive(ifname, request)
	return request, reply, err
}
