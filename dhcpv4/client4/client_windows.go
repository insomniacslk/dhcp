//go:build windows

// Package client4 is deprecated. Use "nclient4" instead.
package client4

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// MaxUDPReceivedPacketSize is the (arbitrary) maximum UDP packet size supported
// by this library. Theoretically could be up to 65kb.
const (
	MaxUDPReceivedPacketSize = 8192
)

var (
	// DefaultReadTimeout is the time to wait after listening in which the
	// exchange is considered failed.
	DefaultReadTimeout = 3 * time.Second

	// DefaultWriteTimeout is the time to wait after sending in which the
	// exchange is considered failed.
	DefaultWriteTimeout = 3 * time.Second
)

// Client is the object that actually performs the DHCP exchange. It currently
// only has read and write timeout values, plus (optional) local and remote
// addresses.
type Client struct {
	ReadTimeout, WriteTimeout time.Duration
	RemoteAddr                net.Addr
	LocalAddr                 net.Addr
}

// NewClient generates a new client to perform a DHCP exchange with, setting the
// read and write timeout fields to defaults.
func NewClient() *Client {
	return &Client{
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
	}
}

// MakeRawUDPPacket is not supported on Windows.
func MakeRawUDPPacket(payload []byte, serverAddr, clientAddr net.UDPAddr) ([]byte, error) {
	return nil, errors.New("MakeRawUDPPacket is not supported on Windows")
}

// MakeBroadcastSocket is not supported on Windows.
func MakeBroadcastSocket(ifname string) (int, error) {
	return -1, errors.New("MakeBroadcastSocket is not supported on Windows")
}

// MakeListeningSocket is not supported on Windows.
func MakeListeningSocket(ifname string) (int, error) {
	return -1, errors.New("MakeListeningSocket is not supported on Windows")
}

func toUDPAddr(addr net.Addr, defaultAddr *net.UDPAddr) (*net.UDPAddr, error) {
	var uaddr *net.UDPAddr
	if addr == nil {
		uaddr = defaultAddr
	} else {
		if a, ok := addr.(*net.UDPAddr); ok {
			uaddr = a
		} else {
			return nil, fmt.Errorf("could not convert to net.UDPAddr, got %T instead", addr)
		}
	}
	if uaddr.IP.To4() == nil {
		return nil, fmt.Errorf("'%s' is not a valid IPv4 address", uaddr.IP)
	}
	return uaddr, nil
}

func (c *Client) getLocalUDPAddr() (*net.UDPAddr, error) {
	defaultLocalAddr := &net.UDPAddr{IP: net.IPv4zero, Port: dhcpv4.ClientPort}
	laddr, err := toUDPAddr(c.LocalAddr, defaultLocalAddr)
	if err != nil {
		return nil, fmt.Errorf("Invalid local address: %s", err)
	}
	return laddr, nil
}

func (c *Client) getRemoteUDPAddr() (*net.UDPAddr, error) {
	defaultRemoteAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: dhcpv4.ServerPort}
	raddr, err := toUDPAddr(c.RemoteAddr, defaultRemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("Invalid remote address: %s", err)
	}
	return raddr, nil
}

// Exchange runs a full DORA transaction: Discover, Offer, Request, Acknowledge,
// over UDP. Does not retry in case of failures.
//
// On Windows, this uses standard UDP sockets. Note that on Windows,
// binding to a specific interface is not supported, so the client will
// listen on all interfaces and filter by transaction ID.
func (c *Client) Exchange(ifname string, modifiers ...dhcpv4.Modifier) ([]*dhcpv4.DHCPv4, error) {
	conversation := make([]*dhcpv4.DHCPv4, 0)
	raddr, err := c.getRemoteUDPAddr()
	if err != nil {
		return nil, err
	}
	laddr, err := c.getLocalUDPAddr()
	if err != nil {
		return nil, err
	}

	// On Windows, we use standard UDP socket
	// Listen on all interfaces since we can't bind to a specific one
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: laddr.Port})
	if err != nil {
		return conversation, fmt.Errorf("failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	// Discover
	discover, err := dhcpv4.NewDiscoveryForInterface(ifname, modifiers...)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, discover)

	// Offer
	offer, err := c.sendReceive(conn, discover, raddr, dhcpv4.MessageTypeOffer)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, offer)

	// Request
	request, err := dhcpv4.NewRequestFromOffer(offer, modifiers...)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, request)

	// Ack
	ack, err := c.sendReceive(conn, request, raddr, dhcpv4.MessageTypeAck)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, ack)

	return conversation, nil
}

// sendReceive sends a DHCP packet and waits for a response
func (c *Client) sendReceive(conn *net.UDPConn, packet *dhcpv4.DHCPv4, raddr *net.UDPAddr, messageType dhcpv4.MessageType) (*dhcpv4.DHCPv4, error) {
	// Set write deadline
	if err := conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout)); err != nil {
		return nil, err
	}

	// Send the packet
	if _, err := conn.WriteTo(packet.ToBytes(), raddr); err != nil {
		return nil, fmt.Errorf("failed to send DHCP packet: %v", err)
	}

	// Set read deadline
	if err := conn.SetReadDeadline(time.Now().Add(c.ReadTimeout)); err != nil {
		return nil, err
	}

	// Receive response
	for {
		buf := make([]byte, MaxUDPReceivedPacketSize)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			return nil, fmt.Errorf("failed to receive DHCP response: %v", err)
		}

		response, err := dhcpv4.FromBytes(buf[:n])
		if err != nil {
			log.Printf("Error parsing DHCPv4 response: %v", err)
			continue
		}

		// Check transaction ID
		if response.TransactionID != packet.TransactionID {
			continue
		}

		// Check opcode
		if response.OpCode != dhcpv4.OpcodeBootReply {
			continue
		}

		// Check message type if specified
		if messageType != dhcpv4.MessageTypeNone && response.MessageType() != messageType {
			continue
		}

		return response, nil
	}
}

// SendReceive is deprecated on Windows.
func (c *Client) SendReceive(sendFd, recvFd int, packet *dhcpv4.DHCPv4, messageType dhcpv4.MessageType) (*dhcpv4.DHCPv4, error) {
	return nil, errors.New("SendReceive with file descriptors is not supported on Windows")
}
