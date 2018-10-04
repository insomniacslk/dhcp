package dhcpv4

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// utility function to set up a client and a server instance and run it in
// background. The caller needs to call Server.Close() once finished.
func setUpClientAndServer(handler Handler) (*Client, *Server) {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
		Zone: "lo",
	}
	s := NewServer(laddr, handler)
	go s.ActivateAndServe()

	c := NewClient()
	for {
		if s.LocalAddr() != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
		log.Printf("Waiting for server to run...")
	}
	c.RemoteAddr = &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: s.LocalAddr().(*net.UDPAddr).Port,
		Zone: "lo",
	}

	return c, s
}

func TestNewServer(t *testing.T) {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
		Zone: "lo",
	}
	handler := func(conn net.PacketConn, peer net.Addr, m *DHCPv4) {}
	s := NewServer(laddr, handler)
	defer s.Close()

	require.NotNil(t, s)
	require.Nil(t, s.conn)
	require.Equal(t, laddr, s.localAddr)
	require.NotNil(t, s.Handler)
}

func TestServerActivateAndServe(t *testing.T) {
	handler := func(conn net.PacketConn, peer net.Addr, m *DHCPv4) {
		offer, err := New()
		if err != nil {
			log.Printf("New offer packet creation failed: %v", err)
			return
		}
		offer.SetOpcode(OpcodeBootReply)
		if _, err := conn.WriteTo(offer.ToBytes(), peer); err != nil {
			log.Printf("Cannot reply to client: %v", err)
		}
	}
	c, s := setUpClientAndServer(handler)
	defer s.Close()

	_, _, err := c.Solicit("lo", nil)

	require.NoError(t, err)
}
