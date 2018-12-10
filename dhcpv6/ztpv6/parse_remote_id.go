package ztpv6

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/insomniacslk/dhcp/dhcpv6"
)

var (
	// Arista Port, Vlan Pattern
	aristaPVPattern = regexp.MustCompile("Ethernet(?P<port>[0-9]+):(?P<vlan>[0-9]+)")
	// Arista Slot, Mod, Port Pattern
	aristaSMPPattern = regexp.MustCompile("Ethernet(?P<slot>[0-9]+)/(?P<module>[0-9]+)/(?P<port>[0-9]+)")
)

// CircuitID represents the structure of network vendor interface formats
type CircuitID struct {
	Slot    string
	Module  string
	Port    string
	SubPort string
	Vlan    string
}

// ParseRemoteId will parse the RemoteId Option data for Vendor Specific data
func ParseRemoteId(packet dhcpv6.DHCPv6) (*CircuitID, error) {
	// Need to decapsulate the packet after multiple relays in order to reach RemoteId data
	inner, err := dhcpv6.DecapsulateRelayIndex(packet, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to decapsulate relay index: %v", err)
	}

	if rid := inner.GetOneOption(dhcpv6.OptionRemoteID); rid != nil {
		remoteID := string(rid.(*dhcpv6.OptRemoteId).RemoteID())
		circ, err := matchCircuitId(remoteID)
		if err != nil {
			return nil, err
		}
		return circ, nil
	}
	return nil, errors.New("failed to parse RemoteID option data")
}

func matchCircuitId(remoteID string) (*CircuitID, error) {
	var names, matches []string

	switch {
	case aristaPVPattern.MatchString(remoteID):
		matches = aristaPVPattern.FindStringSubmatch(remoteID)
		names = aristaPVPattern.SubexpNames()
	case aristaSMPPattern.MatchString(remoteID):
		matches = aristaSMPPattern.FindStringSubmatch(remoteID)
		names = aristaSMPPattern.SubexpNames()
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no circuitId regex matches for %v", remoteID)
	}

	var circuit CircuitID
	for i, match := range matches {
		switch names[i] {
		case "port":
			circuit.Port = match
		case "slot":
			circuit.Slot = match
		case "module":
			circuit.Module = match
		case "subport":
			circuit.SubPort = match
		case "vlan":
			circuit.Vlan = match
		}
	}

	return &circuit, nil
}

// FormatCircuitID is the CircuitID format we send in our Bootfile URL for ZTP devices
func (c *CircuitID) FormatCircuitID() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v", c.Slot, c.Module, c.Port, c.SubPort, c.Vlan)
}
