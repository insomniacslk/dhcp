package dhcpv4

import (
	"encoding/binary"
	"errors"
	"net"
	"os"
	"time"
	"fmt"
	"reflect"
	"log"

	"golang.org/x/net/ipv4"
	"golang.org/x/sys/unix"
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
	RemoteAddr   net.Addr
	LocalAddr    net.Addr
}

// NewClient generates a new client to perform a DHCP exchange with, setting the
// read and write timeout fields to defaults.
func NewClient() *Client {
	return &Client{
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
	}
}

// MakeRawBroadcastPacket leverages MakeRawPacket to create a raw packet suitable
// for UDP broadcast.
func MakeRawBroadcastPacket(payload []byte) ([]byte, error) {
	log.Printf("Warning: dhcpv4.MakeRawBroadcastPacket() is deprecated and will be removed.")
	serverAddr := net.UDPAddr{IP: net.IPv4bcast, Port: ServerPort}
	clientAddr := net.UDPAddr{IP: net.IPv4zero, Port: ClientPort}
	return MakeRawUDPPacket(payload, serverAddr, clientAddr)
}

// MakeRawPacket converts a payload (a serialized DHCPv4 packet) into a
// raw packet for the specified serverAddr from the specified clientAddr.
func MakeRawUDPPacket(payload []byte, serverAddr, clientAddr net.UDPAddr) ([]byte, error) {
	udp := make([]byte, 8)
	binary.BigEndian.PutUint16(udp[:2], uint16(clientAddr.Port))
	binary.BigEndian.PutUint16(udp[2:4], uint16(serverAddr.Port))
	binary.BigEndian.PutUint16(udp[4:6], uint16(8+len(payload)))
	binary.BigEndian.PutUint16(udp[6:8], 0) // try to offload the checksum

	h := ipv4.Header{
		Version:  4,
		Len:      20,
		TotalLen: 20 + len(udp) + len(payload),
		TTL:      64,
		Protocol: 17, // UDP
		Dst:      serverAddr.IP,
		Src:      clientAddr.IP,
	}
	ret, err := h.Marshal()
	if err != nil {
		return nil, err
	}
	ret = append(ret, udp...)
	ret = append(ret, payload...)
	return ret, nil
}

// makeRawSocket creates a socket that can be passed to unix.Sendto.
func makeRawSocket() (int, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_RAW)
	if err != nil {
		return fd, err
	}
	err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	if err != nil {
		return fd, err
	}
	err = unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_HDRINCL, 1)
	if err != nil {
		return fd, err
	}
	return fd, nil
}

// MakeBroadcastSocket creates a socket that can be passed to unix.Sendto
// that will send packets out to the broadcast address.
func MakeBroadcastSocket(ifname string) (int, error) {
	fd, err := makeRawSocket()
	if err != nil {
		return fd, err
	}
	err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_BROADCAST, 1)
	if err != nil {
		return fd, err
	}
	return fd, nil
}

// MakeListeningSocket creates a listening socket on 0.0.0.0 for the DHCP client
// port and returns it.
func MakeListeningSocket(ifname string) (int, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_UDP)
	if err != nil {
		return fd, err
	}
	err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	if err != nil {
		return fd, err
	}
	var addr [4]byte
	copy(addr[:], net.IPv4zero.To4())
	if err = unix.Bind(fd, &unix.SockaddrInet4{Port: ClientPort, Addr: addr}); err != nil {
		return fd, err
	}
	err = BindToInterface(fd, ifname)
	if err != nil {
		return fd, err
	}
	return fd, nil
}

func (c *Client) getUDPAddresses() (*net.UDPAddr, *net.UDPAddr, error) {
	var laddr *net.UDPAddr
	if c.LocalAddr == nil {
		laddr = &net.UDPAddr{IP: net.IPv4bcast, Port: ServerPort}
	} else {
		if addr, ok := c.LocalAddr.(*net.UDPAddr); ok {
			laddr = addr
		} else {
			return nil, nil, fmt.Errorf("Invalid local address: only net.UDPAddr is supported, got %v instead", reflect.TypeOf(c.LocalAddr))
		}
	}
	var raddr *net.UDPAddr
	if c.RemoteAddr == nil {
		raddr = &net.UDPAddr{IP: net.IPv4bcast, Port: ServerPort}
	} else {
		if addr, ok := c.RemoteAddr.(*net.UDPAddr); ok {
			raddr = addr
		} else {
			return nil, nil, fmt.Errorf("Invalid remote address: only net.UDPAddr is supported, got %v instead", reflect.TypeOf(c.RemoteAddr))
		}
	}

	return raddr, laddr, nil
}

// Exchange runs a full DORA transaction: Discover, Offer, Request, Acknowledge,
// over UDP. Does not retry in case of failures. Returns a list of DHCPv4
// structures representing the exchange. It can contain up to four elements,
// ordered as Discovery, Offer, Request and Acknowledge. In case of errors, an
// error is returned, and the list of DHCPv4 objects will be shorted than 4,
// containing all the sent and received DHCPv4 messages.
func (c *Client) Exchange(ifname string, discover *DHCPv4, modifiers ...Modifier) ([]*DHCPv4, error) {
	conversation := make([]*DHCPv4, 0)
	var err error

	// Get our file descriptor for the raw socket we need.
	var sfd int
	// TODO: we should check the IP is actually not a broadcast address for the
	// subnet here.
	if c.RemoteAddr != nil {
		sfd, err = makeRawSocket()
	} else {
		sfd, err = MakeBroadcastSocket(ifname)
	}
	if err != nil {
		return conversation, err
	}
	rfd, err := MakeListeningSocket(ifname)
	if err != nil {
		return conversation, err
	}

	// Discover
	if discover == nil {
		discover, err = NewDiscoveryForInterface(ifname)
		if err != nil {
			return conversation, err
		}
	}
	for _, mod := range modifiers {
		discover = mod(discover)
	}
	conversation = append(conversation, discover)

	// Offer
	offer, err := c.sendReceive(sfd, rfd, discover, MessageTypeOffer)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, offer)

	// Request
	request, err := NewRequestFromOffer(offer, modifiers...)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, request)

	// Ack
	ack, err := c.sendReceive(sfd, rfd, request, MessageTypeAck)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, ack)
	return conversation, nil
}

// SendReceive sends a packet (with some write timeout) and waits for a
// response up to some read timeout value. If the message type is not
// MessageTypeNone, it will wait for a specific message type
func (c *Client) sendReceive(sendFd, recvFd int, packet *DHCPv4, messageType MessageType) (*DHCPv4, error) {
	raddr, laddr, err := c.getUDPAddresses()
	if err != nil {
		return nil, err
	}
	packetBytes, err := MakeRawUDPPacket(packet.ToBytes(), *raddr, *laddr)
	if err != nil {
		return nil, err
	}

	// Create a goroutine to perform the blocking send, and time it out after
	// a certain amount of time.
	var (
		destination [4]byte
		response    *DHCPv4
	)
	copy(destination[:], raddr.IP.To4())
	remoteAddr := unix.SockaddrInet4{Port: laddr.Port, Addr: destination}
	recvErrors := make(chan error, 1)
	go func(errs chan<- error) {
		conn, innerErr := net.FileConn(os.NewFile(uintptr(recvFd), ""))
		if err != nil {
			errs <- innerErr
			return
		}
		defer conn.Close()
		conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))

		for {
			buf := make([]byte, MaxUDPReceivedPacketSize)
			n, _, _, _, innerErr := conn.(*net.UDPConn).ReadMsgUDP(buf, []byte{})
			if innerErr != nil {
				errs <- innerErr
				return
			}

			response, innerErr = FromBytes(buf[:n])
			if err != nil {
				errs <- innerErr
				return
			}
			// check that this is a response to our message
			if response.TransactionID() != packet.TransactionID() {
				continue
			}
			// wait for a response message
			if response.Opcode() != OpcodeBootReply {
				continue
			}
			// if we are not requested to wait for a specific message type,
			// return what we have
			if messageType == MessageTypeNone {
				break
			}
			// break if it's a reply of the desired type, continue otherwise
			if response.MessageType() != nil && *response.MessageType() == messageType {
				break
			}
		}
		recvErrors <- nil
	}(recvErrors)
	if err = unix.Sendto(sendFd, packetBytes, 0, &remoteAddr); err != nil {
		return nil, err
	}

	select {
	case err = <-recvErrors:
		if err != nil {
			return nil, err
		}
	case <-time.After(c.ReadTimeout):
		return nil, errors.New("timed out while listening for replies")
	}

	return response, nil
}

// BroadcastSendReceive broadcasts packet (with some write timeout) and waits for a
// response up to some read timeout value. If the message type is not
// MessageTypeNone, it will wait for a specific message type
func BroadcastSendReceive(sendFd, recvFd int, packet *DHCPv4, readTimeout, writeTimeout time.Duration, messageType MessageType) (*DHCPv4, error) {
	log.Printf("Warning: dhcpv4.BroadcastSendAndReceive() is deprecated and will be removed. You should use dhcpv4.client.Exchange() instead.")
	packetBytes, err := MakeRawBroadcastPacket(packet.ToBytes())
	if err != nil {
		return nil, err
	}

	// Create a goroutine to perform the blocking send, and time it out after
	// a certain amount of time.
	var (
		destination [4]byte
		response    *DHCPv4
	)
	copy(destination[:], net.IPv4bcast.To4())
	remoteAddr := unix.SockaddrInet4{Port: ClientPort, Addr: destination}
	recvErrors := make(chan error, 1)
	go func(errs chan<- error) {
		conn, innerErr := net.FileConn(os.NewFile(uintptr(recvFd), ""))
		if err != nil {
			errs <- innerErr
			return
		}
		defer conn.Close()
		conn.SetReadDeadline(time.Now().Add(readTimeout))

		for {
			buf := make([]byte, MaxUDPReceivedPacketSize)
			n, _, _, _, innerErr := conn.(*net.UDPConn).ReadMsgUDP(buf, []byte{})
			if innerErr != nil {
				errs <- innerErr
				return
			}

			response, innerErr = FromBytes(buf[:n])
			if err != nil {
				errs <- innerErr
				return
			}
			// check that this is a response to our message
			if response.TransactionID() != packet.TransactionID() {
				continue
			}
			// wait for a response message
			if response.Opcode() != OpcodeBootReply {
				continue
			}
			// if we are not requested to wait for a specific message type,
			// return what we have
			if messageType == MessageTypeNone {
				break
			}
			// break if it's a reply of the desired type, continue otherwise
			if response.MessageType() != nil && *response.MessageType() == messageType {
				break
			}
		}
		recvErrors <- nil
	}(recvErrors)
	if err = unix.Sendto(sendFd, packetBytes, 0, &remoteAddr); err != nil {
		return nil, err
	}

	select {
	case err = <-recvErrors:
		if err != nil {
			return nil, err
		}
	case <-time.After(readTimeout):
		return nil, errors.New("timed out while listening for replies")
	}

	return response, nil
}
