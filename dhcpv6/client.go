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
	Dialer       *net.Dialer
	ReadTimeout  *time.Duration
	WriteTimeout *time.Duration
	LocalAddr    net.Addr
	RemoteAddr   net.Addr
}

// Make a stateful DHCPv6 request
func (c *Client) Exchange(ifname string, d *DHCPv6) ([]DHCPv6, error) {
	conversation := make([]DHCPv6, 1)
	var err error

	// Solicit
	if d == nil {
		d, err = NewSolicitForInterface(ifname)
		if err != nil {
			return conversation, err
		}
	}
	conversation[0] = *d
	advertise, err := c.ExchangeSolicitAdvertise(ifname, d)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, *advertise)

	// TODO request/reply
	return conversation, nil
}

func (c *Client) ExchangeSolicitAdvertise(ifname string, d *DHCPv6) (*DHCPv6, error) {
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

	// send the SOLICIT packet out
	_, err = conn.WriteTo(d.ToBytes(), &raddr)
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
