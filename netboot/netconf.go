package netboot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/vishvananda/netlink"
)

// AddrConf holds a single IP address configuration for a NIC
type AddrConf struct {
	IPNet             net.IPNet
	PreferredLifetime int
	ValidLifetime     int
}

// NetConf holds multiple IP configuration for a NIC, and DNS configuration
type NetConf struct {
	Addresses     []AddrConf
	DNSServers    []net.IP
	DNSSearchList []string
}

// GetNetConfFromPacket extracts network configuration information from a DHCPv6
// Reply packet and returns a populated NetConf structure
func GetNetConfFromPacketv6(d *dhcpv6.DHCPv6Message) (*NetConf, error) {
	opt := d.GetOneOption(dhcpv6.OPTION_IA_NA)
	if opt == nil {
		return nil, errors.New("No option IA NA found")
	}
	netconf := NetConf{}
	// get IP configuration
	oiana := opt.(*dhcpv6.OptIANA)
	iaaddrs := make([]*dhcpv6.OptIAAddress, 0)
	for _, o := range oiana.Options() {
		if o.Code() == dhcpv6.OPTION_IAADDR {
			iaaddrs = append(iaaddrs, o.(*dhcpv6.OptIAAddress))
		}
	}
	netmask := net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::"))
	for _, iaaddr := range iaaddrs {
		netconf.Addresses = append(netconf.Addresses, AddrConf{
			IPNet: net.IPNet{
				IP:   iaaddr.IPv6Addr(),
				Mask: netmask,
			},
			PreferredLifetime: int(iaaddr.PreferredLifetime()),
			ValidLifetime:     int(iaaddr.ValidLifetime()),
		})
	}
	// get DNS configuration
	opt = d.GetOneOption(dhcpv6.DNS_RECURSIVE_NAME_SERVER)
	if opt == nil {
		return nil, errors.New("No option DNS Recursive Name Servers found ")
	}
	odnsserv := opt.(*dhcpv6.OptDNSRecursiveNameServer)
	// TODO should this be copied?
	netconf.DNSServers = odnsserv.NameServers()

	opt = d.GetOneOption(dhcpv6.DOMAIN_SEARCH_LIST)
	if opt == nil {
		return nil, errors.New("No option DNS Domain Search List found")
	}
	odomains := opt.(*dhcpv6.OptDomainSearchList)
	// TODO should this be copied?
	netconf.DNSSearchList = odomains.DomainSearchList()

	return &netconf, nil
}

// IfUp brings up an interface by name, and waits for it to come up until a timeout expires
func IfUp(ifname string, timeout time.Duration) (netlink.Link, error) {
	start := time.Now()
	for time.Since(start) < timeout {
		iface, err := netlink.LinkByName(ifname)
		if err != nil {
			return nil, fmt.Errorf("cannot get interface %q by name: %v", ifname, err)
		}

		// if the interface is up, return
		if iface.Attrs().OperState == netlink.OperUp {
			// XXX despite the OperUp state, upon the first attempt I
			// consistently get a "cannot assign requested address" error. This
			// may be a bug in the netlink library. Need to investigate more.
			time.Sleep(time.Second)
			return iface, nil
		}
		// otherwise try to bring it up
		if err := netlink.LinkSetUp(iface); err != nil {
			return nil, fmt.Errorf("interface %q: %v can't make it up: %v", ifname, iface, err)
		}
	}

	return nil, fmt.Errorf("timed out while waiting for %s to come up", ifname)

}

// ConfigureInterface configures a network interface with the configuration held by a
// NetConf structure
func ConfigureInterface(ifname string, netconf *NetConf) error {
	iface, err := netlink.LinkByName(ifname)
	if err != nil {
		return fmt.Errorf("error getting interface information for %s: %v", ifname, err)
	}
	// configure interfaces
	for _, addr := range netconf.Addresses {
		dest := &netlink.Addr{
			IPNet:       &addr.IPNet,
			PreferedLft: addr.PreferredLifetime,
			ValidLft:    addr.ValidLifetime,
		}
		if err := netlink.AddrReplace(iface, dest); err != nil {
			if os.IsExist(err) {
				return fmt.Errorf("cannot configure %s on %s,%d,%d: %v", ifname, addr.IPNet, addr.PreferredLifetime, addr.ValidLifetime, err)
			}
		}
	}
	// configure /etc/resolv.conf
	resolvconf := ""
	for _, ns := range netconf.DNSServers {
		resolvconf += fmt.Sprintf("nameserver %s\n", ns)
	}
	resolvconf += fmt.Sprintf("search %s\n", strings.Join(netconf.DNSSearchList, " "))
	return ioutil.WriteFile("/etc/resolv.conf", []byte(resolvconf), 0644)
}
