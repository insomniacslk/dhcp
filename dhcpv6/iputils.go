package dhcpv6

import (
	"fmt"
	"net"
)

// InterfaceAddresses is used to fetch addresses of an interface with given name
var InterfaceAddresses func(string) ([]net.Addr, error) = interfaceAddresses

func interfaceAddresses(ifname string) ([]net.Addr, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	return iface.Addrs()
}

func getMatchingAddr(ifname string, matches func(net.IP) bool) (net.IP, error) {
	ifaddrs, err := InterfaceAddresses(ifname)
	if err != nil {
		return nil, err
	}
	for _, ifaddr := range ifaddrs {
		if ifaddr, ok := ifaddr.(*net.IPNet); ok && matches(ifaddr.IP) {
			return ifaddr.IP, nil
		}
	}
	return nil, fmt.Errorf("no matching address found for interface %s", ifname)
}

// GetLinkLocalAddr returns a link-local address for the interface
func GetLinkLocalAddr(ifname string) (net.IP, error) {
	return getMatchingAddr(ifname, func(ip net.IP) bool {
		return ip.To4() == nil && ip.IsLinkLocalUnicast()
	})
}

// GetGlobalAddr returns a global address for the interface
func GetGlobalAddr(ifname string) (net.IP, error) {
	return getMatchingAddr(ifname, func(ip net.IP) bool {
		return ip.To4() == nil && ip.IsGlobalUnicast()
	})
}
