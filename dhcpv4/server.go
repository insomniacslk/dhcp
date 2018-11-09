package dhcpv4

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

/*
  To use the DHCPv4 server code you have to call NewServer with two arguments:
  - a handler function, that will be called every time a valid DHCPv4 packet is
    received, and
  - an address to listen on.

  The handler is a function that takes as input a packet connection, that can be
  used to reply to the client; a peer address, that identifies the client sending
  the request, and the DHCPv4 packet itself. Just implement your custom logic in
  the handler.

  The address to listen on is used to know IP address, port and optionally the
  scope to create and UDP socket to listen on for DHCPv4 traffic.

  Example program:


package main

import (
	"log"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func handler(conn net.PacketConn, peer net.Addr, m dhcpv4.DHCPv4) {
	// this function will just print the received DHCPv4 message, without replying
	log.Print(m.Summary())
}

func main() {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 67,
	}
	server := dhcpv4.NewServer(laddr, handler)

	defer server.Close()
	if err := server.ActivateAndServe(); err != nil {
		log.Panic(err)
	}
	// Serve for 60 seconds
	time.Sleep(60 * time.Second)
}

*/

// Handler is a type that defines the handler function to be called every time a
// valid DHCPv4 message is received
type Handler func(conn net.PacketConn, peer net.Addr, m *DHCPv4)

// Server represents a DHCPv4 server object
type Server struct {
	conn       net.PacketConn
	connMutex  sync.Mutex
	shouldStop chan bool
	Handler    Handler
	localAddr  net.UDPAddr
}

// LocalAddr returns the local address of the listening socket, or nil if not
// listening
func (s *Server) LocalAddr() net.Addr {
	s.connMutex.Lock()
	defer s.connMutex.Unlock()
	if s.conn == nil {
		return nil
	}
	return s.conn.LocalAddr()
}

// ActivateAndServe starts the DHCPv4 server. The listener will run in
// background, and can be interrupted with `Server.Close`.
func (s *Server) ActivateAndServe() error {
	s.connMutex.Lock()
	// FIXME check if conn != nil instead, and close it if so. It may
	//       panic if it's in an invalid state, but better than assuming
	//       that the already open connection is usable.
	if s.conn == nil {
		conn, err := net.ListenUDP("udp4", &s.localAddr)
		if err != nil {
			s.connMutex.Unlock()
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
		return fmt.Errorf("error: not an UDPConn")
	}
	if pc == nil {
		return fmt.Errorf("ActivateAndServe: invalid nil PacketConn")
	}
	log.Printf("Server listening on %s", pc.LocalAddr())
	log.Print("Ready to handle requests")
	go func() {
		for {
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
					if !err.(net.Error).Timeout() {
						// TODO report errors through a channel
						log.Print(err)
						return
					}
					// if timeout, silently skip and continue
				default:
					// complain and continue
					log.Printf("Error reading from packet conn: %v", err)
				}
				continue
			}
			log.Printf("Handling request from %v", peer)
			m, err := FromBytes(rbuf[:n])
			if err != nil {
				log.Printf("Error parsing DHCPv4 request: %v", err)
				continue
			}
			go s.Handler(pc, peer, m)
		}
	}()
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
