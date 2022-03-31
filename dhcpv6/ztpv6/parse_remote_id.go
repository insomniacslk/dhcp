package ztpv6

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/insomniacslk/dhcp/dhcpv6"
)

var circuitRegexs = []*regexp.Regexp{
	// Arista Port, Vlan Pattern
	regexp.MustCompile("Ethernet(?P<port>[0-9]+):(?P<vlan>[0-9]+)"),
	// Arista Slot, Mod, Port Pattern
	regexp.MustCompile("Ethernet(?P<slot>[0-9]+)/(?P<mod>[0-9]+)/(?P<port>[0-9]+)"),
	// Juniper QFX et-0/0/0:0.0 and xe-0/0/0:0.0
	regexp.MustCompile("^(et|xe)-(?P<slot>[0-9]+)/(?P<mod>[0-9]+)/(?P<port>[0-9]+):(?P<subport>[0-9]+).*$"),
	// Juniper PTX et-0/0/0.0
	regexp.MustCompile("^et-(?P<slot>[0-9]+)/(?P<mod>[0-9]+)/(?P<port>[0-9]+).(?P<subport>[0-9]+)$"),
	// Juniper EX ge-0/0/0.0
	regexp.MustCompile("^ge-(?P<slot>[0-9]+)/(?P<mod>[0-9]+)/(?P<port>[0-9]+).(?P<subport>[0-9]+).*"),
	// Arista Ethernet3/17/1
	// Sometimes Arista prepend circuit id type(1 byte) and length(1 byte) not using ^
	regexp.MustCompile("Ethernet(?P<slot>[0-9]+)/(?P<mod>[0-9]+)/(?P<port>[0-9]+)$"),
	// Juniper QFX et-1/0/61
	regexp.MustCompile("^et-(?P<slot>[0-9]+)/(?P<mod>[0-9]+)/(?P<port>[0-9]+)$"),
	// Arista Ethernet14:Vlan2001
	// Arista Ethernet10:2020
	regexp.MustCompile("Ethernet(?P<port>[0-9]+):(?P<vlan>.*)$"),
	// Cisco Gi1/10:2020
	regexp.MustCompile("^Gi(?P<slot>[0-9]+)/(?P<port>[0-9]+):(?P<vlan>.*)$"),
	// Nexus Ethernet1/3
	regexp.MustCompile("^Ethernet(?P<slot>[0-9]+)/(?P<port>[0-9]+)$"),
	// Juniper bundle interface ae52.0
	regexp.MustCompile("^ae(?P<port>[0-9]+).(?P<subport>[0-9])$"),
}

// CircuitID represents the structure of network vendor interface formats
type CircuitID struct {
	Slot    string
	Module  string
	Port    string
	SubPort string
	Vlan    string
}

// ParseRemoteId will parse the RemoteId Option data for Vendor Specific data
func ParseRemoteID(packet dhcpv6.DHCPv6) (*CircuitID, error) {
	// Need to decapsulate the packet after multiple relays in order to reach RemoteId data
	inner, err := dhcpv6.DecapsulateRelayIndex(packet, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to decapsulate relay index: %v", err)
	}

	if rm, ok := inner.(*dhcpv6.RelayMessage); ok {
		if rid := rm.Options.RemoteID(); rid != nil {
			remoteID := string(rid.RemoteID)
			circ, err := matchCircuitId(remoteID)
			if err == nil {
				return circ, nil
			}
		}
		// if we fail to find circuit id from remote id try to use interface ID option
		if iid := rm.Options.InterfaceID(); iid != nil {
			interfaceID := string(iid)
			circ, err := matchCircuitId(interfaceID)
			if err == nil {
				return circ, nil
			}
		}
	}
	return nil, errors.New("failed to parse RemoteID and InterfaceID option data")
}

func matchCircuitId(circuitInfo string) (*CircuitID, error) {
	for _, re := range circuitRegexs {

		match := re.FindStringSubmatch(circuitInfo)
		if len(match) == 0 {
			continue
		}

		c := CircuitID{}
		for i, k := range re.SubexpNames() {
			switch k {
			case "slot":
				c.Slot = match[i]
			case "mod":
				c.Module = match[i]
			case "port":
				c.Port = match[i]
			case "subport":
				c.SubPort = match[i]
			case "vlan":
				c.Vlan = match[i]
			}
		}

		return &c, nil
	}
	return nil, fmt.Errorf("Unable to match circuit id : %s with listed regexes of interface types", circuitInfo)
}

// FormatCircuitID is the CircuitID format we send in our Bootfile URL for ZTP devices
func (c *CircuitID) FormatCircuitID() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v", c.Slot, c.Module, c.Port, c.SubPort, c.Vlan)
}
