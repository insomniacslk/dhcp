package ztpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchCircuitID(t *testing.T) {
	tt := []struct {
		name    string
		circuit string
		want    *CircuitID
		fail    bool 
	}{
		{name: "Bogus string", circuit: "bogus_interface", fail: true, want: nil},
		{name: "juniperQFX pattern", circuit: "et-0/0/0:0.0", want: &CircuitID{Slot: "0", Module: "0", Port: "0", SubPort: "0"}},
		{name: "juniperPTX pattern", circuit: "et-0/0/0.0", want: &CircuitID{Slot: "0", Module: "0", Port: "0"}},
		{name: "Arista pattern", circuit: "Ethernet3/17/1", want: &CircuitID{Slot: "3", Module: "17", Port: "1"}},
		{name: "Juniper QFX pattern", circuit: "et-1/0/61", want: &CircuitID{Slot: "1", Module: "0", Port: "61"}},
		{name: "Arista Vlan pattern 1", circuit: "Ethernet14:Vlan2001", want: &CircuitID{Port: "14", Vlan: "Vlan2001"}},
		{name: "Arista Vlan pattern 2", circuit: "Ethernet10:2020", want: &CircuitID{Port: "10", Vlan: "2020"}},
		{name: "Cisco pattern", circuit: "Gi1/10:2020", want: &CircuitID{Slot: "1", Port: "10", Vlan: "2020"}},
		{name: "Cisco Nexus pattern", circuit: "Ethernet1/3", want: &CircuitID{Slot: "1", Port: "3"}},
		{name: "Juniper Bundle Pattern", circuit: "ae52.0", want: &CircuitID{Port: "52", SubPort: "0"}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			circuit, err := matchCircuitId(tc.circuit)
			if err != nil && !tc.fail {
				t.Errorf("unexpected failure: %v", err)
			}
			if circuit != nil {
				require.Equal(t, tc.want, circuit, "comparing remoteID data")
			}
		})
	}
}

func TestFormatCircuitID(t *testing.T) {
	tt := []struct {
		name    string
		circuit *CircuitID
		want    string
		fail    bool
	}{
		{name: "empty", circuit: &CircuitID{}, want: ",,,,"},
		{name: "juniperQFX pattern", circuit: &CircuitID{Slot: "0", Module: "0", Port: "0", SubPort: "0"}, want: "0,0,0,0,"},
		{name: "juniperPTX pattern", circuit: &CircuitID{Slot: "0", Module: "0", Port: "0"}, want: "0,0,0,,"},
		{name: "Arista pattern", circuit: &CircuitID{Slot: "3", Module: "17", Port: "1"}, want: "3,17,1,,"},
		{name: "Juniper QFX pattern", circuit: &CircuitID{Slot: "1", Module: "0", Port: "61"}, want: "1,0,61,,"},
		{name: "Arista Vlan pattern 1", circuit: &CircuitID{Port: "14", Vlan: "Vlan2001"}, want: ",,14,,Vlan2001"},
		{name: "Arista Vlan pattern 2", circuit: &CircuitID{Port: "10", Vlan: "2020"}, want: ",,10,,2020"},
		{name: "Cisco Nexus pattern", circuit: &CircuitID{Slot: "1", Port: "3"}, want: "1,,3,,"},
		{name: "Juniper Bundle Pattern", circuit: &CircuitID{Port: "52", SubPort: "0"}, want: ",,52,0,"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			circuit := tc.circuit.FormatCircuitID()
			require.Equal(t, tc.want, circuit, "FormatRemoteID data")
		})
	}

}

