package async

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/iana"
)


// solicit creates new solicit based on the mac address
func solicit(input string) (dhcpv6.DHCPv6, error) {
	mac, err := net.ParseMAC(input)
	if err != nil {
		return nil, err
	}
	duid := dhcpv6.Duid{
		Type:          dhcpv6.DUID_LLT,
		HwType:        iana.HwTypeEthernet,
		Time:          dhcpv6.GetTime(),
		LinkLayerAddr: mac,
	}
	return dhcpv6.NewSolicitWithCID(duid)
}

// server creates a server which responds with predefined answers
func serve(t *testing.T, addr *net.UDPAddr, responses ...dhcpv6.DHCPv6) {
	conn, err := net.ListenUDP("udp6", addr)
	require.NoError(t, err)
	defer conn.Close()
	oobdata := []byte{}
	buffer := make([]byte, dhcpv6.MaxUDPReceivedPacketSize)
	for _, packet := range responses {
		n, _, _, src, err := conn.ReadMsgUDP(buffer, oobdata)
		require.NoError(t, err)
		_, err = dhcpv6.FromBytes(buffer[:n])
		require.NoError(t, err)
		_, err = conn.WriteTo(packet.ToBytes(), src)
		require.NoError(t, err)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	require.NotNil(t, c)
	require.Equal(t, c.ReadTimeout, dhcpv6.DefaultReadTimeout)
	require.Equal(t, c.ReadTimeout, dhcpv6.DefaultWriteTimeout)
}

func TestOpenInvalidAddrFailes(t *testing.T) {
	c := NewClient()
	err := c.Open(512)
	require.Error(t, err)
}

// This test uses port 15438 so please make sure its not used before running
func TestOpenClose(t *testing.T) {
	c := NewClient()
	addr, err := net.ResolveUDPAddr("udp6", ":15438")
	require.NoError(t, err)
	c.LocalAddr = addr
	err = c.Open(512)
	require.NoError(t, err)
	defer c.Close()
}

// This test uses ports 15438 and 15439 so please make sure they are not used
// before running
func TestSendTimeout(t *testing.T) {
	c := NewClient()
	addr, err := net.ResolveUDPAddr("udp6", ":15438")
	require.NoError(t, err)
	remote, err := net.ResolveUDPAddr("udp6", ":15439")
	require.NoError(t, err)
	c.ReadTimeout = 50 * time.Millisecond
	c.WriteTimeout = 50 * time.Millisecond
	c.LocalAddr = addr
	c.RemoteAddr = remote
	err = c.Open(512)
	require.NoError(t, err)
	defer c.Close()
	m, err := dhcpv6.NewMessage()
	require.NoError(t, err)
	_, err, timeout := c.Send(m).GetOrTimeout(200)
	require.NoError(t, err)
	require.True(t, timeout)
}

// This test uses ports 15438 and 15439 so please make sure they are not used
// before running
func TestSend(t *testing.T) {
	s, err := solicit("c8:6c:2c:47:96:fd")
	require.NoError(t, err)
	require.NotNil(t, s)

	a, err := dhcpv6.NewAdvertiseFromSolicit(s)
	require.NoError(t, err)
	require.NotNil(t, a)

	c := NewClient()
	addr, err := net.ResolveUDPAddr("udp6", ":15438")
	require.NoError(t, err)
	remote, err := net.ResolveUDPAddr("udp6", ":15439")
	require.NoError(t, err)
	c.LocalAddr = addr
	c.RemoteAddr = remote

	go serve(t, remote, a)

	err = c.Open(16)
	require.NoError(t, err)
	defer c.Close()

	f := c.Send(s)
	response, err, timeout := f.GetOrTimeout(1000)
	require.False(t, timeout)
	require.NoError(t, err)
	require.Equal(t, a, response)
}
