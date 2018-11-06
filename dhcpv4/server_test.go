package dhcpv4

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

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
	offer, err := New()
	if err != nil {
		log.Printf("New Offer packet creation failed: %v", err)
		return
	}
	// TODO implement NewOfferFromDiscovery and NewAckFromRequest,
	//      and replace the bulk of code below that generates the
	//      response packet
	offer.SetOpcode(OpcodeBootReply)
	offer.SetTransactionID(m.TransactionID())
	offer.AddOption(&OptServerIdentifier{ServerID: net.IP{1, 2, 3, 4}})
	offer.AddOption(&OptMessageType{MessageType: MessageTypeOffer})
	if _, err := conn.WriteTo(offer.ToBytes(), peer); err != nil {
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

// utility function to return the loopback interface name
// TODO this is copied from dhcpv6/server_test.go , we should refactor common code in a separate package
func getLoopbackInterface() (string, error) {
	var ifaces []net.Interface
	var err error
	if ifaces, err = net.Interfaces(); err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Name[:2] == "lo" {
			return iface.Name, nil
		}
	}
	return "", errors.New("No loopback interface found")
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

	lo, err := getLoopbackInterface()
	require.NoError(t, err)

	hwaddr := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	discover, err := NewDiscovery(hwaddr)
	require.NoError(t, err)

	discover.SetTransactionID(0xaabbccdd)
	discover.SetUnicast()

	_, err = c.Exchange(lo, discover)
	require.NoError(t, err)
}
