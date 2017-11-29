package dhcpv6

import (
	"fmt"
	"net"
)

func GetLinkLocalAddr(ifname string) (*net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Name != ifname {
			continue
		}
		ifaddrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, ifaddr := range ifaddrs {
			if ifaddr, ok := ifaddr.(*net.IPNet); ok {
				if ifaddr.IP.To4() == nil && ifaddr.IP.IsLinkLocalUnicast() {
					return &ifaddr.IP, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("No link-local address found for interface %v", ifname)
}
