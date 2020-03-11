package server6

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/nclient6"
	"github.com/insomniacslk/dhcp/interfaces"
	"github.com/stretchr/testify/require"
)

// Turns a connected UDP conn into an "unconnected" UDP conn.
type unconnectedConn struct {
	*net.UDPConn
}

func (f unconnectedConn) WriteTo(b []byte, _ net.Addr) (int, error) {
	return f.UDPConn.Write(b)
}

func (f unconnectedConn) ReadFrom(b []byte) (int, net.Addr, error) {
	n, err := f.Read(b)
	return n, nil, err
}

// utility function to set up a client and a server instance and run it in
// background. The caller needs to call Server.Close() once finished.
func setUpClientAndServer(handler Handler, logger *customLogger) (*nclient6.Client, *Server) {
	laddr := &net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 0,
	}
	s, err := NewServer("", laddr, handler)
	if err != nil {
		panic(err)
	}

	if logger != nil {
		s.logger = logger
	}

	go func() {
		_ = s.Serve()
	}()

	clientConn, err := net.DialUDP("udp6", &net.UDPAddr{IP: net.ParseIP("::1")}, s.conn.LocalAddr().(*net.UDPAddr))
	if err != nil {
		panic(err)
	}

	c, err := nclient6.NewWithConn(unconnectedConn{clientConn}, net.HardwareAddr{1, 2, 3, 4, 5, 6})
	if err != nil {
		panic(err)
	}
	return c, s
}

type customLogger struct {
	tb testing.TB
	called bool
	mux    sync.Mutex
}

func (s *customLogger) Printf(format string, v ...interface{}) {
	s.mux.Lock()
	s.called = true
	s.mux.Unlock()
	s.tb.Logf("===CustomLogger BEGIN===")
	s.tb.Logf(format, v...)
	s.tb.Logf("===CustomLogger END===")
}

func (s *customLogger) PrintMessage(prefix string, message *dhcpv6.Message) {
	s.mux.Lock()
	s.called = true
	s.mux.Unlock()
	s.tb.Logf("===CustomLogger BEGIN===")
	s.tb.Logf("%s: %s", prefix, message)
	s.tb.Logf("===CustomLogger END===")
}

func TestServer(t *testing.T) {
	handler := func(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
		msg := m.(*dhcpv6.Message)
		adv, err := dhcpv6.NewAdvertiseFromSolicit(msg)
		if err != nil {
			log.Printf("NewAdvertiseFromSolicit failed: %v", err)
			return
		}
		if _, err := conn.WriteTo(adv.ToBytes(), peer); err != nil {
			log.Printf("Cannot reply to client: %v", err)
		}
	}

	c, s := setUpClientAndServer(handler, nil)
	defer s.Close()

	ifaces, err := interfaces.GetLoopbackInterfaces()
	require.NoError(t, err)
	require.NotEqual(t, 0, len(ifaces))

	_, err = c.Solicit(context.Background(), dhcpv6.WithRapidCommit)
	require.NoError(t, err)
}

func TestCustomLoggerForServer(t *testing.T) {
	handler := func(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
		msg := m.(*dhcpv6.Message)
		adv, err := dhcpv6.NewAdvertiseFromSolicit(msg)
		if err != nil {
			log.Printf("NewAdvertiseFromSolicit failed: %v", err)
			return
		}
		if _, err := conn.WriteTo(adv.ToBytes(), peer); err != nil {
			log.Printf("Cannot reply to client: %v", err)
		}
	}

	c, s := setUpClientAndServer(handler, &customLogger{
		tb: t,
	})
	defer s.Close()

	ifaces, err := interfaces.GetLoopbackInterfaces()
	require.NoError(t, err)
	require.NotEqual(t, 0, len(ifaces))

	_, err = c.Solicit(context.Background(), dhcpv6.WithRapidCommit)
	require.NoError(t, err)
	go func() {
		time.Sleep(time.Second * 5)
		require.Equal(t, true, s.logger.(*customLogger).called)
	}()
}

func TestServerInstantiationWithCustomLogger(t *testing.T) {
	handler := func(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
		msg := m.(*dhcpv6.Message)
		adv, err := dhcpv6.NewAdvertiseFromSolicit(msg)
		if err != nil {
			log.Printf("NewAdvertiseFromSolicit failed: %v", err)
			return
		}
		if _, err := conn.WriteTo(adv.ToBytes(), peer); err != nil {
			log.Printf("Cannot reply to client: %v", err)
		}
	}

	laddr := &net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 0,
	}
	s, err := NewServer("", laddr, handler, WithLogger(&customLogger{
		tb: t,
	}))
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, s)
}

func TestServerInstantiationWithSummaryLogger(t *testing.T) {
	handler := func(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
		msg := m.(*dhcpv6.Message)
		adv, err := dhcpv6.NewAdvertiseFromSolicit(msg)
		if err != nil {
			log.Printf("NewAdvertiseFromSolicit failed: %v", err)
			return
		}
		if _, err := conn.WriteTo(adv.ToBytes(), peer); err != nil {
			log.Printf("Cannot reply to client: %v", err)
		}
	}

	laddr := &net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 0,
	}
	s, err := NewServer("", laddr, handler, WithSummaryLogger())
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, s)
}

func TestServerInstantiationWithDebugLogger(t *testing.T) {
	handler := func(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
		msg := m.(*dhcpv6.Message)
		adv, err := dhcpv6.NewAdvertiseFromSolicit(msg)
		if err != nil {
			log.Printf("NewAdvertiseFromSolicit failed: %v", err)
			return
		}
		if _, err := conn.WriteTo(adv.ToBytes(), peer); err != nil {
			log.Printf("Cannot reply to client: %v", err)
		}
	}

	laddr := &net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 0,
	}
	s, err := NewServer("", laddr, handler, WithDebugLogger())
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, s)
}
