//+build integration

package netboot

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// integration tests that require Linux with a properly working rtnetlink
// interface, and the existence of an "eth0" interface.
// WARNING: these tests may improperly configure your network interfaces and
// routing, so be careful before running them. Privileged access and integration
// build tag required to run them.

var (
	testIfname = "eth0"
)

func TestInit(t *testing.T) {
	var r RTNL
	err := r.init()
	assert.NoError(t, err)
	require.NotNil(t, r.conn)
	r.Close()
}

func TestClose(t *testing.T) {
	var r RTNL
	err := r.init()
	assert.NoError(t, err)
	require.NotNil(t, r.conn)
	r.Close()
	require.Nil(t, r.conn)
}

func TestGetLinkState(t *testing.T) {
	var r RTNL
	defer r.Close()

	iface, err := net.InterfaceByName(testIfname)
	require.NoError(t, err)
	_, err = r.GetLinkState(iface.Index)
	require.NoError(t, err)
}

func TestSetLinkState(t *testing.T) {
	var r RTNL
	defer r.Close()

	iface, err := net.InterfaceByName(testIfname)
	require.NoError(t, err)
	err = r.SetLinkState(iface.Index, true)
	require.NoError(t, err)
}

func Test_getFamily(t *testing.T) {
	require.Equal(t, unix.AF_INET, getFamily(net.IPv4zero))
	require.Equal(t, unix.AF_INET, getFamily(net.IPv4bcast))
	require.Equal(t, unix.AF_INET, getFamily(net.IPv4allrouter))

	require.Equal(t, unix.AF_INET6, getFamily(net.IPv6zero))
	require.Equal(t, unix.AF_INET6, getFamily(net.IPv6loopback))
	require.Equal(t, unix.AF_INET6, getFamily(net.IPv6linklocalallrouters))
}

func TestSetAddr(t *testing.T) {
	var r RTNL
	defer r.Close()
	iface, err := net.InterfaceByName(testIfname)
	require.NoError(t, err)

	a := net.IPNet{IP: net.ParseIP("10.0.123.1"), Mask: net.IPv4Mask(255, 255, 255, 0)}
	err = r.SetAddr(iface.Index, a)
	require.NoError(t, err)
	// TODO implement GetAddr to further validate this, and minimize the effect
	// of concurrent tests that may invalidate this check.
}
