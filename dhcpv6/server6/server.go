package server6

import (
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv6"
)

// Handler is a type that defines the handler function to be called every time a
// valid DHCPv6 message is received
type Handler func(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6)

// Server represents a DHCPv6 server object
type Server struct {
	conn    net.PacketConn
	handler Handler
}

// Serve starts the DHCPv6 server. The listener will run in background, and can
// be interrupted with `Server.Close`.
func (s *Server) Serve() {
	log.Printf("Server listening on %s", s.conn.LocalAddr())
	log.Print("Ready to handle requests")

	for {
		rbuf := make([]byte, 4096) // FIXME this is bad
		n, peer, err := s.conn.ReadFrom(rbuf)
		if err != nil {
			log.Printf("Error reading from packet conn: %v", err)
			return
		}
		log.Printf("Handling request from %v", peer)

		d, err := dhcpv6.FromBytes(rbuf[:n])
		if err != nil {
			log.Printf("Error parsing DHCPv6 request: %v", err)
			continue
		}

		go s.handler(s.conn, peer, d)
	}
}

// Close sends a termination request to the server, and closes the UDP listener
func (s *Server) Close() error {
	return s.conn.Close()
}

// A ServerOpt configures a Server.
type ServerOpt func(s *Server)

// WithConn configures a server with the given connection.
func WithConn(conn net.PacketConn) ServerOpt {
	return func(s *Server) {
		s.conn = conn
	}
}

// NewServer initializes and returns a new Server object
func NewServer(addr *net.UDPAddr, handler Handler, opt ...ServerOpt) (*Server, error) {
	s := &Server{
		handler: handler,
	}

	for _, o := range opt {
		o(s)
	}

	if s.conn == nil {
		conn, err := net.ListenUDP("udp6", addr)
		if err != nil {
			return nil, err
		}
		s.conn = conn
	}
	return s, nil
}
