package dhcpv6

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewAsyncClient(t *testing.T) {
	c := NewAsyncClient()
	require.NotNil(t, c)
	require.Equal(t, c.ReadTimeout, DefaultReadTimeout)
	require.Equal(t, c.ReadTimeout, DefaultWriteTimeout)
}

func TestOpenInvalidAddrFailes(t *testing.T) {
	c := NewAsyncClient()
	err := c.Open(512)
	require.Error(t, err)
}

// This test uses port 15438 so please make sure its not used before running
func TestOpenClose(t *testing.T) {
	c := NewAsyncClient()
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
	c := NewAsyncClient()
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
	m, err := NewMessage()
	require.NoError(t, err)
	_, err = c.Send(m).WaitTimeout(200 * time.Millisecond)
	require.Error(t, err)
}
