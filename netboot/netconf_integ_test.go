// +build integration

package netboot

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// this assumes that eth0 exists and is configurable
var ifname = "eth0"

func TestIfUp(t *testing.T) {
	iface, err := IfUp(ifname, 2*time.Second)
	require.NoError(t, err)
	assert.Equal(t, ifname, iface.Name)
}

func TestIfUpTimeout(t *testing.T) {
	_, err := IfUp(ifname, 0*time.Second)
	require.Error(t, err)
}

func TestConfigureInterface(t *testing.T) {
	nc := NetConf{
		Addresses: []AddrConf{
			AddrConf{IPNet: net.IPNet{IP: net.ParseIP("10.20.30.40")}},
		},
	}
	err := ConfigureInterface(ifname, &nc)
	require.NoError(t, err)
}

func TestConfigureInterfaceWithRouteAndDNS(t *testing.T) {
	nc := NetConf{
		Addresses: []AddrConf{
			AddrConf{IPNet: net.IPNet{IP: net.ParseIP("10.20.30.40")}},
		},
		DNSServers:    []net.IP{net.ParseIP("8.8.8.8")},
		DNSSearchList: []string{"slackware.it"},
		Routers:       []net.IP{net.ParseIP("10.20.30.254")},
	}
	err := ConfigureInterface(ifname, &nc)
	require.NoError(t, err)
}
