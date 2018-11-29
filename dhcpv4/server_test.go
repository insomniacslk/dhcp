// +build integration

package dhcpv4

import (
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/insomniacslk/dhcp/interfaces"
	"github.com/stretchr/testify/require"
)

func init() {
	// initialize seed. This is generally bad, but "good enough"
	// to generate random ports for these tests
	rand.Seed(time.Now().UTC().UnixNano())
}

func randPort() int {
	// can't use port 0 with raw sockets, so until we implement
	// a non-raw-sockets client for non-static ports, we have to
	// deal with this "randomness"
	return 1024 + rand.Intn(65536-1024)
}

// DORAHandler is a server handler suitable for DORA transactions
func DORAHandler(conn net.PacketConn, peer net.Addr, m *DHCPv4) {
	if m == nil {
		log.Printf("Packet is nil!")
		return
	}
	if m.Opcode() != OpcodeBootRequest {
		log.Printf("Not a BootRequest!")
		return
	}
	reply, err := NewReplyFromRequest(m)
	if err != nil {
		log.Printf("NewReplyFromRequest failed: %v", err)
		return
	}
	reply.AddOption(&OptServerIdentifier{ServerID: net.IP{1, 2, 3, 4}})
	opt := m.GetOneOption(OptionDHCPMessageType)
	if opt == nil {
		log.Printf("No message type found!")
		return
	}
	switch opt.(*OptMessageType).MessageType {
	case MessageTypeDiscover:
		reply.AddOption(&OptMessageType{MessageType: MessageTypeOffer})
	case MessageTypeRequest:
		reply.AddOption(&OptMessageType{MessageType: MessageTypeAck})
	default:
		log.Printf("Unhandled message type: %v", opt.(*OptMessageType).MessageType)
		return
	}

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		log.Printf("Cannot reply to client: %v", err)
	}
}

// utility function to set up a client and a server instance and run it in
// background. The caller needs to call Server.Close() once finished.
func setUpClientAndServer(handler Handler) (*Client, *Server) {
	// strong assumption, I know
	loAddr := net.ParseIP("127.0.0.1")
	laddr := net.UDPAddr{
		IP:   loAddr,
		Port: randPort(),
	}
	s := NewServer(laddr, handler)
	go s.ActivateAndServe()

	c := NewClient()
	// FIXME this doesn't deal well with raw sockets, the actual 0 will be used
	// in the UDP header as source port
	c.LocalAddr = &net.UDPAddr{IP: loAddr, Port: randPort()}
	for {
		if s.LocalAddr() != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
		log.Printf("Waiting for server to run...")
	}
	c.RemoteAddr = s.LocalAddr()
	log.Printf("Client.RemoteAddr: %s", c.RemoteAddr)

	return c, s
}

func TestNewServer(t *testing.T) {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
	}
	s := NewServer(laddr, DORAHandler)
	defer s.Close()

	require.NotNil(t, s)
	require.Nil(t, s.conn)
	require.Equal(t, laddr, s.localAddr)
	require.NotNil(t, s.Handler)
}

func TestServerActivateAndServe(t *testing.T) {
	c, s := setUpClientAndServer(DORAHandler)
	defer s.Close()

	ifaces, err := interfaces.GetLoopbackInterfaces()
	require.NoError(t, err)
	require.NotEqual(t, 0, len(ifaces))

	xid := uint32(0xaabbccdd)
	hwaddr := [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}

	modifiers := []Modifier{
		WithTransactionID(xid),
		WithHwAddr(hwaddr[:]),
	}

	conv, err := c.Exchange(ifaces[0].Name, modifiers...)
	require.NoError(t, err)
	require.Equal(t, 4, len(conv))
	for _, p := range conv {
		require.Equal(t, xid, p.TransactionID())
		require.Equal(t, [16]byte(hwaddr), p.ClientHwAddr())
	}
}
