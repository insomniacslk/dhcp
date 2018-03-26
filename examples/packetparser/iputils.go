package main

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/milosgajdos83/tenus"
)

// GetLinkLocalAddr returns the link-local address and the network for a given
// network interface, or an error if any.
func GetLinkLocalAddr(ifname string) (*net.IP, *net.IPNet, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}
	var iface *net.Interface
	var linkLocalAddr *net.IP
	for _, ifi := range ifaces {
		if ifi.Name == ifname {
			iface = &ifi
			break
		}
	}
	// build the addr from the interface
	hwa := iface.HardwareAddr
	linkLocalAddr = &net.IP{
		0xfe, 0x80, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		hwa[0] ^ 2, hwa[1], hwa[2], 0xff,
		0xfe, hwa[3], hwa[4], hwa[5],
	}
	m := net.CIDRMask(64, 128)
	linkLocalNet := net.IPNet{IP: linkLocalAddr.Mask(m), Mask: m} // a /64
	return linkLocalAddr, &linkLocalNet, nil
}

// WaitForInterfaceStatusUp waits until a network interface is UP and ready to
// be used. If the interface is not ready within the given timeout, an error is
// returned.
func WaitForInterfaceStatusUp(ifname string, timeout time.Duration) error {
	// FIXME should use netlink events rather than polling like this
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	deadline := time.Now().Add(timeout)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("Timed out while waiting for interface to be up")
		}
		for _, ifi := range ifaces {
			if ifi.Flags&net.FlagUp != 0 {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// ConfigureLinkLocalAddress brings up an interface and configures its
// link-local address.
func ConfigureLinkLocalAddress(ifname string) (*net.IP, error) {
	// Configure the link-local address for the given interface, via Linux's
	// netlink, and bring the interface up
	llAddr, llNet, err := GetLinkLocalAddr(ifname)
	if err != nil {
		return nil, err
	}
	dl, err := tenus.NewLinkFrom(ifname)
	if err != nil {
		return nil, err
	}
	addrs, err := dl.NetInterface().Addrs()
	if err != nil {
		return nil, err
	}
	found := false
	for _, addr := range addrs {
		if uAddr, ok := addr.(*net.IPNet); ok {
			if uAddr.IP.To16() != nil && bytes.Equal(uAddr.IP, *llAddr) {
				found = true
			}
		}
	}
	if !found {
		if err = dl.SetLinkIp(*llAddr, llNet); err != nil {
			return nil, err
		}
		if err = dl.SetLinkUp(); err != nil {
			return nil, err
		}
	}
	return llAddr, nil
}
