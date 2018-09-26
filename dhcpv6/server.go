package dhcpv6

import (
	"fmt"
	"log"
	"net"
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

	if err := server.ActivateAndServe(); err != nil {
		log.Fatal(err)
	}
}

*/

type Handler func(conn net.PacketConn, peer net.Addr, m DHCPv6)

type Server struct {
	conn      net.PacketConn
	LocalAddr net.UDPAddr
	Handler   Handler
}

func (s *Server) ActivateAndServe() error {
	if s.conn == nil {
		conn, err := net.ListenUDP("udp6", &s.LocalAddr)
		if err != nil {
			return err
		}
		s.conn = conn
	}
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
		log.Printf("Waiting..")
		rbuf := make([]byte, 4096) // FIXME this is bad
		n, peer, err := pc.ReadFrom(rbuf)
		if err != nil {
			log.Printf("Error reading from packet conn: %v", err)
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

func NewServer(addr net.UDPAddr, handler Handler) *Server {
	return &Server{
		LocalAddr: addr,
		Handler:   handler,
	}
}
