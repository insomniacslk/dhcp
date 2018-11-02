package netboot

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
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
	Routers       []net.IP
}

// GetNetConfFromPacketv6 extracts network configuration information from a DHCPv6
// Reply packet and returns a populated NetConf structure
func GetNetConfFromPacketv6(d *dhcpv6.DHCPv6Message) (*NetConf, error) {
	opt := d.GetOneOption(dhcpv6.OptionIANA)
	if opt == nil {
		return nil, errors.New("No option IA NA found")
	}
	netconf := NetConf{}
	// get IP configuration
	oiana := opt.(*dhcpv6.OptIANA)
	iaaddrs := make([]*dhcpv6.OptIAAddress, 0)
	for _, o := range oiana.Options {
		if o.Code() == dhcpv6.OptionIAAddr {
			iaaddrs = append(iaaddrs, o.(*dhcpv6.OptIAAddress))
		}
	}
	netmask := net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::"))
	for _, iaaddr := range iaaddrs {
		netconf.Addresses = append(netconf.Addresses, AddrConf{
			IPNet: net.IPNet{
				IP:   iaaddr.IPv6Addr,
				Mask: netmask,
			},
			PreferredLifetime: int(iaaddr.PreferredLifetime),
			ValidLifetime:     int(iaaddr.ValidLifetime),
		})
	}
	// get DNS configuration
	opt = d.GetOneOption(dhcpv6.OptionDNSRecursiveNameServer)
	if opt == nil {
		return nil, errors.New("No option DNS Recursive Name Servers found ")
	}
	odnsserv := opt.(*dhcpv6.OptDNSRecursiveNameServer)
	// TODO should this be copied?
	netconf.DNSServers = odnsserv.NameServers

	opt = d.GetOneOption(dhcpv6.OptionDomainSearchList)
	if opt == nil {
		return nil, errors.New("No option DNS Domain Search List found")
	}
	odomains := opt.(*dhcpv6.OptDomainSearchList)
	// TODO should this be copied?
	netconf.DNSSearchList = odomains.DomainSearchList

	return &netconf, nil
}

// GetNetConfFromPacketv4 extracts network configuration information from a DHCPv4
// Reply packet and returns a populated NetConf structure
func GetNetConfFromPacketv4(d *dhcpv4.DHCPv4) (*NetConf, error) {
	// extract the address from the DHCPv4 address
	ipAddr := d.YourIPAddr()
	if ipAddr.Equal(net.IPv4zero) {
		return nil, errors.New("ip address is null (0.0.0.0)")
	}
	netconf := NetConf{}

	// get the subnet mask from OptionSubnetMask. If the netmask is not defined
	// in the packet, an error is returned
	netmaskOption := d.GetOneOption(dhcpv4.OptionSubnetMask)
	if netmaskOption == nil {
		return nil, errors.New("no netmask option in response packet")
	}
	netmask := netmaskOption.(*dhcpv4.OptSubnetMask).SubnetMask
	if binary.LittleEndian.Uint32(netmask) == 0 {
		return nil, errors.New("netmask extracted from OptSubnetMask options is null")
	}

	// netconf struct requires a valid lifetime to be specified. ValidLifetime is a dhcpv6
	// concept, the closest mapping in dhcpv4 world is "IP Address Lease Time". If the lease
	// time option is nil, we set it to 0
	leaseTimeOption := d.GetOneOption(dhcpv4.OptionIPAddressLeaseTime)
	leaseTime := uint32(0)
	if leaseTimeOption != nil {
		leaseTime = leaseTimeOption.(*dhcpv4.OptIPAddressLeaseTime).LeaseTime
	}

	if int(leaseTime) < 0 {
		return nil, fmt.Errorf("lease time overflow, Original lease time: %d", leaseTime)
	}

	netconf.Addresses = append(netconf.Addresses, AddrConf{
		IPNet: net.IPNet{
			IP:   ipAddr,
			Mask: netmask,
		},
		PreferredLifetime: 0,
		ValidLifetime:     int(leaseTime),
	})
	if bytes.Equal(netmask, []byte{0, 0, 0, 0}) {
		return nil, errors.New("netmask in response packet is null")
	}

	// get DNS configuration
	dnsServersOption := d.GetOneOption(dhcpv4.OptionDomainNameServer)
	if dnsServersOption == nil {
		return nil, errors.New("no dns server option in response packet")
	}
	dnsServers := dnsServersOption.(*dhcpv4.OptDomainNameServer).NameServers
	if len(dnsServers) == 0 {
		return nil, errors.New("no dns servers options in response packet")
	}
	netconf.DNSServers = dnsServers

	// get domain search list
	dnsDomainSearchListOption := d.GetOneOption(dhcpv4.OptionDNSDomainSearchList)
	if dnsDomainSearchListOption == nil {
		return nil, errors.New("no domain search list option in response packet")

	}
	dnsSearchList := dnsDomainSearchListOption.(*dhcpv4.OptDomainSearch).DomainSearch
	if len(dnsSearchList) == 0 {
		return nil, errors.New("dns search list is empty")
	}
	netconf.DNSSearchList = dnsSearchList

	// get default gateway
	routerOption := d.GetOneOption(dhcpv4.OptionRouter)
	if routerOption == nil {
		return nil, errors.New("no router option specified in reponse packet")
	}

	routersList := routerOption.(*dhcpv4.OptRouter).Routers
	if len(routersList) == 0 {
		return nil, errors.New("no routers specified in the corresponding option")
	}

	netconf.Routers = routersList

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
	if err = ioutil.WriteFile("/etc/resolv.conf", []byte(resolvconf), 0644); err != nil {
		return fmt.Errorf("could not write resolv.conf file %v", err)
	}

	// add default route information for v4 space. only one default route is allowed
	// so ignore the others if there are multiple ones
	if len(netconf.Routers) > 0 {
		iface, err = netlink.LinkByName(ifname)
		if err != nil {
			return fmt.Errorf("could not obtain interface when adding default route: %v", err)
		}
		// if there is a default v4 route, remove it, as we want to add the one we just got during
		// the dhcp transaction. if the route is not present, which is the final state we want,
		// an error is returned so ignore it
		dst := &net.IPNet{
			IP:   net.IPv4(0, 0, 0, 0),
			Mask: net.CIDRMask(0, 32),
		}
		route := netlink.Route{LinkIndex: iface.Attrs().Index, Dst: dst, Src: net.IPv4(0, 0, 0, 0)}
		netlink.RouteDel(&route)

		src := netconf.Addresses[0].IPNet.IP
		route = netlink.Route{LinkIndex: iface.Attrs().Index, Dst: dst, Src: src, Gw: netconf.Routers[0]}
		err = netlink.RouteAdd(&route)
		if err != nil {
			return fmt.Errorf("could not add default route: %v", err)
		}
	}

	return nil
}
