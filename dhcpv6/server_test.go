package dhcpv6

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 0,
	}
	handler := func(conn net.PacketConn, peer net.Addr, m DHCPv6) {}
	s := NewServer(laddr, handler)
	require.NotNil(t, s)
	require.Nil(t, s.conn)
	require.Equal(t, laddr, s.LocalAddr)
	require.NotNil(t, s.Handler)
}
