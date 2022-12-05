package nclient4

import (
	"net"
	"testing"

	"github.com/vishvananda/netlink"
)

const (
	linkName = "neigh0"
	ipStr    = "10.99.0.1"
	macStr   = "aa:bb:cc:dd:00:01"
)

func TestGetHwAddrFromLocalCache(t *testing.T) {
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		t.Fatal(err)
	}
	ip := net.ParseIP(ipStr)

	if err := addNeigh(ip, mac); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := delNeigh(ip, mac); err != nil {
			t.Fatal(err)
		}
	}()

	_, err = net.InterfaceByName(linkName)
	if err != nil {
		t.Fatal(err)
	}

	if hw, err := getHwAddr(ip); err != nil && hw != nil && hw.String() == macStr {
		t.Fatal(err)
	}
}

func addNeigh(ip net.IP, mac net.HardwareAddr) error {
	dummy := netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: linkName}}
	if err := netlink.LinkAdd(&dummy); err != nil {
		return err
	}
	newlink, err := netlink.LinkByName(dummy.Name)
	if err != nil {
		return err
	}
	dummy.Index = newlink.Attrs().Index

	return netlink.NeighAdd(&netlink.Neigh{
		LinkIndex:    dummy.Index,
		State:        netlink.NUD_REACHABLE,
		IP:           ip,
		HardwareAddr: mac,
	})
}

func delNeigh(ip net.IP, mac net.HardwareAddr) error {
	dummy, err := netlink.LinkByName(linkName)
	if err != nil {
		return err
	}

	return netlink.NeighDel(&netlink.Neigh{
		LinkIndex:    dummy.Attrs().Index,
		State:        netlink.NUD_REACHABLE,
		IP:           ip,
		HardwareAddr: mac,
	})
}

