package dhcpv6

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

/*
  To use the DHCPv6 server code you have to call NewServer with two arguments:
  - a handler function, that will be called every time a valid DHCPv6 packet is
      received, and
  - an address to listen on.

  The handler is a function that takes as input a packet connection, that can be
  used to reply to the client; a peer address, that identifies the client sending
  the request, and the DHCPv6 packet itself. Just implement your custom logic in
  the handler.

  The address to listen on is used to know IP address, port and optionally the
  scope to create and UDP6 socket to listen on for DHCPv6 traffic.

  Example program:


package main

import (
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv6"
)

func handler(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
	// this function will just print the received DHCPv6 message, without replying
	log.Print(m.Summary())
}

func main() {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 547,
	}
	server := dhcpv6.NewServer(laddr, handler)

	defer server.Close()
	if err := server.ActivateAndServe(); err != nil {
		log.Panic(err)
	}
}

*/

// Handler is a type that defines the handler function to be called every time a
// valid DHCPv6 message is received
type Handler func(conn net.PacketConn, peer net.Addr, m DHCPv6)

// Server represents a DHCPv6 server object
type Server struct {
	conn       net.PacketConn
	connMutex  sync.Mutex
	shouldStop chan bool
	Handler    Handler
	localAddr  net.UDPAddr
}

// LocalAddr returns the local address of the listening socked, or nil if not
// listening
func (s *Server) LocalAddr() net.Addr {
	s.connMutex.Lock()
	defer s.connMutex.Unlock()
	if s.conn == nil {
		return nil
	}
	return s.conn.LocalAddr()
}

// ActivateAndServe starts the DHCPv6 server
func (s *Server) ActivateAndServe() error {
	s.connMutex.Lock()
	if s.conn == nil {
		conn, err := net.ListenUDP("udp6", &s.localAddr)
		if err != nil {
			return err
		}
		s.conn = conn
	}
	s.connMutex.Unlock()
	var (
		pc *net.UDPConn
		ok bool
	)
	if pc, ok = s.conn.(*net.UDPConn); !ok {
		return fmt.Errorf("Error: not an UDPConn")
	}
	if pc == nil {
		return fmt.Errorf("ActivateAndServe: Invalid nil PacketConn")
	}
	log.Printf("Server listening on %s", pc.LocalAddr())
	log.Print("Ready to handle requests")
	for {
		log.Printf("CHECK")
		select {
		case <-s.shouldStop:
			break
		case <-time.After(time.Millisecond):
		}
		pc.SetReadDeadline(time.Now().Add(time.Second))
		rbuf := make([]byte, 4096) // FIXME this is bad
		n, peer, err := pc.ReadFrom(rbuf)
		if err != nil {
			switch err.(type) {
			case net.Error:
				// silently skip and continue
			default:
				//complain and continue
				log.Printf("Error reading from packet conn: %v", err)
			}
			continue
		}
		log.Printf("Handling request from %v", peer)
		m, err := FromBytes(rbuf[:n])
		if err != nil {
			log.Printf("Error parsing DHCPv6 request: %v", err)
			continue
		}
		s.Handler(pc, peer, m)
	}
	s.conn.Close()
	return nil
}

// Close sends a termination request to the server, and closes the UDP listener
func (s *Server) Close() error {
	s.shouldStop <- true
	s.connMutex.Lock()
	defer s.connMutex.Unlock()
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// NewServer initializes and returns a new Server object
func NewServer(addr net.UDPAddr, handler Handler) *Server {
	return &Server{
		localAddr:  addr,
		Handler:    handler,
		shouldStop: make(chan bool, 1),
	}
}
