package ztpv6

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/stretchr/testify/require"
)

func TestCircuitID(t *testing.T) {
	tt := []struct {
		name    string
		circuit string
		want    *CircuitID
		fail    bool
	}{
		{name: "Bogus string", circuit: "ope/1/2/3:ope", fail: true, want: nil},
		{name: "Arista Port Vlan Pattern", circuit: "Ethernet13:2001", want: &CircuitID{Port: "13", Vlan: "2001"}},
		{name: "Arista Slot Module Port Pattern", circuit: "Ethernet1/3/4", want: &CircuitID{Slot: "1", Module: "3", Port: "4"}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			circuit, err := matchCircuitId(tc.circuit)
			if err != nil && !tc.fail {
				t.Errorf("unexpected failure: %v", err)
			}
			if circuit != nil {
				require.Equal(t, *tc.want, *circuit, "comparing remoteID data")
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
		{name: "Arista format Port/Vlan", circuit: &CircuitID{Port: "13", Vlan: "2001"}, want: ",,13,,2001"},
		{name: "Arista format Slot/Module/Port", circuit: &CircuitID{Slot: "1", Module: "3", Port: "4"}, want: "1,3,4,,"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			circuit := tc.circuit.FormatCircuitID()
			require.Equal(t, tc.want, circuit, "FormatRemoteID data")
		})
	}

}

func TestParseRemoteID(t *testing.T) {
	tt := []struct {
		name    string
		circuit []byte
		want    *CircuitID
		fail    bool
	}{
		{name: "Bogus string", circuit: []byte("ope/1/2/3:ope.1"), fail: true, want: nil},
		{name: "Arista Port Vlan Pattern", circuit: []byte("Ethernet13:2001"), want: &CircuitID{Port: "13", Vlan: "2001"}},
		{name: "Arista Slot Module Port Pattern", circuit: []byte("Ethernet1/3/4"), want: &CircuitID{Slot: "1", Module: "3", Port: "4"}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			packet, err := dhcpv6.NewMessage()
			if err != nil {
				t.Fatalf("failed to creat dhcpv6 packet object: %v", err)
			}
			opt := dhcpv6.OptRemoteId{}
			opt.SetRemoteID(tc.circuit)
			opt.SetEnterpriseNumber(1234)
			packet.AddOption(&opt)

			circuit, err := ParseRemoteId(packet)
			if err != nil && !tc.fail {
				t.Errorf("unexpected failure: %v", err)
			}
			if circuit != nil {
				require.Equal(t, *tc.want, *circuit, "ZTPRemoteID data")
			} else {
				require.Equal(t, tc.want, circuit, "ZTPRemoteID data")
			}
		})
	}
}
