// +build integration

package netboot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The test assumes that the interface exists and is configurable.
// If you are running this test locally, you may need to adjust this value.
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
	// Linux-only. `netboot.ConfigureInterface` writes to /etc/resolv.conf when
	// `NetConf.DNSServers` is set. In this test we make a backup of resolv.conf
	// and subsequently restore it. This is really ugly, and not safe if
	// multiple tests do the same.
	resolvconf, err := ioutil.ReadFile("/etc/resolv.conf")
	if err != nil {
		panic(fmt.Sprintf("Failed to read /etc/resolv.conf: %v", err))
	}
	type testCase struct {
		Name    string
		NetConf *NetConf
	}
	testCases := []testCase{
		{
			Name: "just IP addr",
			NetConf: &NetConf{
				Addresses: []AddrConf{
					AddrConf{IPNet: net.IPNet{IP: net.ParseIP("10.20.30.40")}},
				},
			},
		},
		{
			Name: "IP addr, DNS, and routers",
			NetConf: &NetConf{
				Addresses: []AddrConf{
					AddrConf{IPNet: net.IPNet{IP: net.ParseIP("10.20.30.40")}},
				},
				DNSServers:    []net.IP{net.ParseIP("8.8.8.8")},
				DNSSearchList: []string{"slackware.it"},
				Routers:       []net.IP{net.ParseIP("10.20.30.254")},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			require.NoError(t, ConfigureInterface(ifname, tc.NetConf))

			// after the test, restore the content of /etc/resolv.conf . The permissions
			// are used only if it didn't exist.
			if err = ioutil.WriteFile("/etc/resolv.conf", resolvconf, 0644); err != nil {
				panic(fmt.Sprintf("Failed to restore /etc/resolv.conf: %v", err))
			}
			log.Printf("Restored /etc/resolv.conf")
		})
	}
}
