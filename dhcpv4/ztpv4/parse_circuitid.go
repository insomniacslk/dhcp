package ztpv4

import (
	"fmt"
	"regexp"

)

// CircuitID represents the structure of network vendor interface formats
type CircuitID struct {
	Slot		string
	Module	string
	Port		string
	SubPort string
	Vlan		string
}

var circuitRegexs = []*regexp.Regexp{
	// Juniper QFX et-0/0/0:0.0
	regexp.MustCompile(".*(?P<slot>[0-9]+)/(?P<mod>[0-9]+)/(?P<port>[0-9]+):(?P<subport>[0-9]+).*$"),
	// Juniper PTX et-0/0/0.0
	regexp.MustCompile(".*(?P<slot>[0-9]+)/(?P<mod>[0-9]+)/(?P<port>[0-9]+).*$"),
	// Arista Ethernet3/17/1
	// Juniper QFX et-1/0/61
	regexp.MustCompile(".*(?P<slot>[0-9]+)/(?P<mod>[0-9]+)/(?P<port>[0-9]+)$"),
	// Arista Ethernet14:Vlan2001
	// Arista Ethernet10:2020
	regexp.MustCompile(".*Ethernet(?P<port>[0-9]+):(?P<vlan>.*)$"),
	// Cisco Gi1/10:2020
	regexp.MustCompile(".*(?P<slot>[0-9]+)/(?P<port>[0-9]+):(?P<vlan>.*)$"),
	// Nexus Ethernet1/3
	regexp.MustCompile(".*(?P<slot>[0-9]+)/(?P<port>[0-9]+)$"),
	// Juniper bundle interface ae52.0
	regexp.MustCompile("^ae(?P<port>[0-9]+).(?P<subport>[0-9])$"),
}

func matchCircuitId(circuitId string) (*CircuitID, error) {

	for _, re := range circuitRegexs {

		match := re.FindStringSubmatch(circuitId)
		if len(match) == 0 {
			continue
		}

		//matchMap := map[string]string{}
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
	return nil, fmt.Errorf("Unable to match circuit id : %s with listed regexes of interface types", circuitId)
}

// FormatCircuitID is the CircuitID format we send in our Bootfile URL for ZTP devices 
func (c *CircuitID) FormatCircuitID() string {                                         
	return fmt.Sprintf("%v,%v,%v,%v,%v", c.Slot, c.Module, c.Port, c.SubPort, c.Vlan)    
}                                                                                      
