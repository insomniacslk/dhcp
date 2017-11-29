package dhcpv6

import (
	"fmt"
	"log"
	"net"
)

type ResponseWriter interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	WriteMsg(DHCPv6) error
	Write([]byte) (int, error)
	Close() error
}

type Handler interface {
	ServeDHCP(w ResponseWriter, m *DHCPv6)
}

type Server struct {
	PacketConn net.PacketConn
	Handler    Handler
}

func (s *Server) ActivateAndServe() error {
	if s.PacketConn == nil {
		return fmt.Errorf("Error: no packet connection specified")
	}
	var pc *net.UDPConn
	var ok bool
	if pc, ok = s.PacketConn.(*net.UDPConn); !ok {
		return fmt.Errorf("Error: not an UDPConn")
	}
	if pc == nil {
		return fmt.Errorf("ActivateAndServe: Invalid nil PacketConn")
	}
	log.Print("Handling requests")
	for {
		rbuf := make([]byte, 1024) // FIXME this is bad
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
		log.Print(m.Summary())
		// FIXME use s.Handler
	}
	return nil
}
