package server6

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv6"
)

// Handler is a type that defines the handler function to be called every time a
// valid DHCPv6 message is received
type Handler func(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6)

// Server represents a DHCPv6 server object
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

// ActivateAndServe starts the DHCPv6 server. The listener will run in
// background, and can be interrupted with `Server.Close`.
func (s *Server) ActivateAndServe() error {
	s.connMutex.Lock()
	if s.conn != nil {
		// this may panic if s.conn is closed but not reset properly. For that
		// you should use `Server.Close`.
		s.Close()
	}
	conn, err := net.ListenUDP("udp6", &s.localAddr)
	if err != nil {
		s.connMutex.Unlock()
		return err
	}
	s.conn = conn
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
					return err
				}
				// if timeout, silently skip and continue
			default:
				// complain and continue
				log.Printf("Error reading from packet conn: %v", err)
			}
			continue
		}
		log.Printf("Handling request from %v", peer)
		m, err := dhcpv6.FromBytes(rbuf[:n])
		if err != nil {
			log.Printf("Error parsing DHCPv6 request: %v", err)
			continue
		}
		go s.Handler(pc, peer, m)
	}
}

// Close sends a termination request to the server, and closes the UDP listener
func (s *Server) Close() error {
	s.shouldStop <- true
	s.connMutex.Lock()
	defer s.connMutex.Unlock()
	if s.conn != nil {
		ret := s.conn.Close()
		s.conn = nil
		return ret
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
