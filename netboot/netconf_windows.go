//go:build windows

package netboot

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
)

// AddrConf holds a single IP address configuration for a NIC
type AddrConf struct {
	IPNet             net.IPNet
	PreferredLifetime time.Duration
	ValidLifetime     time.Duration
}

// NetConf holds multiple IP configuration for a NIC, and DNS configuration
type NetConf struct {
	Addresses     []AddrConf
	DNSServers    []net.IP
	DNSSearchList []string
	Routers       []net.IP
	NTPServers    []net.IP
}

// GetNetConfFromPacketv6 extracts network configuration information from a DHCPv6
// Reply packet and returns a populated NetConf structure
func GetNetConfFromPacketv6(d *dhcpv6.Message) (*NetConf, error) {
	iana := d.Options.OneIANA()
	if iana == nil {
		return nil, errors.New("no option IA NA found")
	}
	netconf := NetConf{}

	for _, iaaddr := range iana.Options.Addresses() {
		netconf.Addresses = append(netconf.Addresses, AddrConf{
			IPNet: net.IPNet{
				IP:   iaaddr.IPv6Addr,
				Mask: net.CIDRMask(128, 128),
			},
			PreferredLifetime: iaaddr.PreferredLifetime,
			ValidLifetime:     iaaddr.ValidLifetime,
		})
	}
	// get DNS configuration
	netconf.DNSServers = d.Options.DNS()

	// get domain search list
	domains := d.Options.DomainSearchList()
	if domains != nil {
		netconf.DNSSearchList = domains.Labels
	}

	// get NTP servers
	netconf.NTPServers = d.Options.NTPServers()

	return &netconf, nil
}

// GetNetConfFromPacketv4 extracts network configuration information from a DHCPv4
// Reply packet and returns a populated NetConf structure
func GetNetConfFromPacketv4(d *dhcpv4.DHCPv4) (*NetConf, error) {
	// extract the address from the DHCPv4 address
	ipAddr := d.YourIPAddr
	if ipAddr == nil || ipAddr.Equal(net.IPv4zero) {
		return nil, errors.New("ip address is null (0.0.0.0)")
	}
	netconf := NetConf{}

	// get the subnet mask from OptionSubnetMask
	netmask := d.SubnetMask()
	if netmask == nil {
		return nil, errors.New("no netmask option in response packet")
	}
	ones, _ := netmask.Size()
	if ones == 0 {
		return nil, errors.New("netmask extracted from OptSubnetMask options is null")
	}

	leaseTime := d.IPAddressLeaseTime(0)

	netconf.Addresses = append(netconf.Addresses, AddrConf{
		IPNet: net.IPNet{
			IP:   ipAddr,
			Mask: netmask,
		},
		PreferredLifetime: 0,
		ValidLifetime:     leaseTime,
	})

	// get DNS configuration
	netconf.DNSServers = d.DNS()

	// get domain search list
	dnsSearchList := d.DomainSearch()
	if dnsSearchList != nil {
		if len(dnsSearchList.Labels) == 0 {
			return nil, errors.New("dns search list is empty")
		}
		netconf.DNSSearchList = dnsSearchList.Labels
	}

	// get default gateway
	netconf.Routers = d.Router()

	// get NTP servers
	netconf.NTPServers = d.NTPServers()

	return &netconf, nil
}

// IfUp brings up an interface by name, and waits for it to come up until a timeout expires.
// On Windows, this simply checks if the interface exists and is ready.
func IfUp(ifname string, timeout time.Duration) (*net.Interface, error) {
	start := time.Now()
	for time.Since(start) < timeout {
		iface, err := net.InterfaceByName(ifname)
		if err != nil {
			return nil, err
		}
		// Check if interface is up
		if iface.Flags&net.FlagUp != 0 {
			return iface, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil, fmt.Errorf("timed out while waiting for %s to come up", ifname)
}

// ConfigureInterface configures a network interface with the configuration held by a
// NetConf structure. On Windows, this is not implemented.
func ConfigureInterface(ifname string, netconf *NetConf) error {
	return errors.New("ConfigureInterface is not implemented on Windows; use netsh or Windows API instead")
}
