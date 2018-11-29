package dhcpv6

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/insomniacslk/dhcp/interfaces"
	"github.com/stretchr/testify/require"
)

// utility function to set up a client and a server instance and run it in
// background. The caller needs to call Server.Close() once finished.
func setUpClientAndServer(handler Handler) (*Client, *Server) {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 0,
	}
	s := NewServer(laddr, handler)
	go s.ActivateAndServe()

	c := NewClient()
	c.LocalAddr = &net.UDPAddr{
		IP: net.ParseIP("::1"),
	}
	for {
		if s.LocalAddr() != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
		log.Printf("Waiting for server to run...")
	}
	c.RemoteAddr = &net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: s.LocalAddr().(*net.UDPAddr).Port,
	}

	return c, s
}

func TestNewServer(t *testing.T) {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 0,
	}
	handler := func(conn net.PacketConn, peer net.Addr, m DHCPv6) {}
	s := NewServer(laddr, handler)
	defer s.Close()

	require.NotNil(t, s)
	require.Nil(t, s.conn)
	require.Equal(t, laddr, s.localAddr)
	require.NotNil(t, s.Handler)
}

func TestServerActivateAndServe(t *testing.T) {
	handler := func(conn net.PacketConn, peer net.Addr, m DHCPv6) {
		adv, err := NewAdvertiseFromSolicit(m)
		if err != nil {
			log.Printf("NewAdvertiseFromSolicit failed: %v", err)
			return
		}
		if _, err := conn.WriteTo(adv.ToBytes(), peer); err != nil {
			log.Printf("Cannot reply to client: %v", err)
		}
	}
	c, s := setUpClientAndServer(handler)
	defer s.Close()

	ifaces, err := interfaces.GetLoopbackInterfaces()
	require.NoError(t, err)
	require.NotEqual(t, 0, len(ifaces))

	_, _, err = c.Solicit(ifaces[0].Name)
	require.NoError(t, err)
}
