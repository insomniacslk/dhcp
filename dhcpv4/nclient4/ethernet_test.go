//go:build go1.12
// +build go1.12

package nclient4

import (
	"bytes"
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/u-root/uio/uio"
)

var (
	testMac = net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	// Test payload 123456789abcdefghijklmnopqrstvwxyz
	// This length is required to avoid zero-padding the Ethernet frame from the right side
	testPayload = gopacket.Payload{0x54, 0x65, 0x73, 0x74, 0x20, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x20, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x76, 0x77, 0x78, 0x79, 0x7a}

	ethHdrIPv4 = &layers.Ethernet{
		DstMAC:       BroadcastMac,
		SrcMAC:       testMac,
		EthernetType: layers.EthernetTypeIPv4,
	}
	ethHdrVLAN = &layers.Ethernet{
		DstMAC:       BroadcastMac,
		SrcMAC:       testMac,
		EthernetType: layers.EthernetTypeDot1Q,
	}
	ethHdrARP = &layers.Ethernet{
		DstMAC:       BroadcastMac,
		SrcMAC:       testMac,
		EthernetType: layers.EthernetTypeARP,
	}
	vlanTagOuter = &layers.Dot1Q{
		Priority:       0,
		DropEligible:   false,
		VLANIdentifier: 100,
		Type:           layers.EthernetTypeDot1Q,
	}
	vlanTagInner = &layers.Dot1Q{
		Priority:       0,
		DropEligible:   false,
		VLANIdentifier: 200,
		Type:           layers.EthernetTypeIPv4,
	}
	ipv4Hdr = &layers.IPv4{
		SrcIP: net.IP{1, 2, 3, 4},
		DstIP: net.IP{5, 6, 7, 8},
	}
	opts = gopacket.SerializeOptions{}
)

func TestProcessVLANStack(t *testing.T) {
	for _, tt := range []struct {
		name        string
		inputBytes  []byte
		vlanConfig  []uint16
		wantSuccess bool
	}{
		{"no VLANs + no VLANs configured", []byte{0x08, 0x00}, []uint16{}, true},
		{"no VLANs + VLAN configured", []byte{0x08, 0x00}, []uint16{100}, false},
		{"valid VLAN stack (single)", []byte{0x81, 0x00, 0x01, 0x00, 0x08, 0x00}, []uint16{0x100}, true},
		{"invalid VLAN stack (single)", []byte{0x81, 0x00, 0x01, 0xFF, 0x08, 0x00}, []uint16{0x100}, false},
		{"valid VLAN stack (double)", []byte{0x81, 0x00, 0x01, 0x00, 0x81, 0x00, 0x02, 0x00, 0x08, 0x00}, []uint16{0x100, 0x200}, true},
		{"invalid VLAN stack (double)", []byte{0x81, 0x00, 0x01, 0x00, 0x81, 0x00, 0x02, 0xFF, 0x08, 0x00}, []uint16{0x100, 0x200}, false},
		{"invalid VLAN stack (too short)", []byte{0x81, 0x00, 0x01, 0x00, 0x08, 0x00}, []uint16{0x100, 0x200}, false},
		{"invalid VLAN stack (too long)", []byte{0x81, 0x00, 0x01, 0x00, 0x81, 0x00, 0x02, 0x00, 0x08, 0x00}, []uint16{0x100}, false},
		{"invalid packet (ARP)", []byte{0x08, 0x06}, []uint16{}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testBuf := uio.NewBigEndianBuffer(tt.inputBytes)
			testSuccess := processVLANStack(testBuf, tt.vlanConfig)

			if testSuccess != tt.wantSuccess {
				t.Errorf("got %v, want %v", testSuccess, tt.wantSuccess)
			}
		})
	}
}

func TestCreateVLANTag(t *testing.T) {
	// Gopacket builds VLAN tags the other way around: first VLAN ID/TCI, then TPID, due to their different layered approach
	// Since a VLAN tag is only 4 bytes, and the value is well-known, it makes sense to just construct the packet by hand.
	want := []byte{0x81, 0x00, 0x01, 0x23}

	test := createVLANTag(0x0123)

	if !bytes.Equal(test, want) {
		t.Errorf("got %v, want %v", test, want)
	}
}

func TestGetEthernetPayload(t *testing.T) {
	for _, tt := range []struct {
		name       string
		testLayers []gopacket.SerializableLayer
		wantLayers []gopacket.SerializableLayer
		vlanConfig []uint16
	}{
		{"no VLAN", []gopacket.SerializableLayer{ethHdrIPv4, ipv4Hdr, testPayload}, []gopacket.SerializableLayer{ipv4Hdr, testPayload}, []uint16{}},
		{"single VLAN", []gopacket.SerializableLayer{ethHdrVLAN, vlanTagInner, ipv4Hdr, testPayload}, []gopacket.SerializableLayer{ipv4Hdr, testPayload}, []uint16{200}},
		{"QinQ", []gopacket.SerializableLayer{ethHdrVLAN, vlanTagOuter, vlanTagInner, ipv4Hdr, testPayload}, []gopacket.SerializableLayer{ipv4Hdr, testPayload}, []uint16{100, 200}},
		{"invalid VLAN", []gopacket.SerializableLayer{ethHdrVLAN, vlanTagInner, ipv4Hdr, testPayload}, nil, []uint16{300}},
		{"invalid packet (ARP)", []gopacket.SerializableLayer{ethHdrARP}, nil, []uint16{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testBuf := gopacket.NewSerializeBuffer()
			err := gopacket.SerializeLayers(testBuf, opts, tt.testLayers...)
			if err != nil {
				t.Errorf("error serializing in gopacket (not our bug!)")
			}

			var want []byte
			if tt.wantLayers == nil {
				want = nil
			} else {
				wantBuf := gopacket.NewSerializeBuffer()
				err := gopacket.SerializeLayers(wantBuf, opts, tt.wantLayers...)
				if err != nil {
					t.Errorf("error serializing in gopacket (not our bug!)")
				}
				want = wantBuf.Bytes()
			}

			testBytes := testBuf.Bytes()
			test := getEthernetPayload(testBytes, tt.vlanConfig)

			if !bytes.Equal(test, want) {
				t.Errorf("got %v, want %v", test, want)
			}
		})
	}
}

func TestAddEthernetHdrTwo(t *testing.T) {
	for _, tt := range []struct {
		name       string
		testLayers []gopacket.SerializableLayer
		wantLayers []gopacket.SerializableLayer
		vlanConfig []uint16
	}{
		{"no VLAN", []gopacket.SerializableLayer{ipv4Hdr, testPayload}, []gopacket.SerializableLayer{ethHdrIPv4, ipv4Hdr, testPayload}, []uint16{}},
		{"single VLAN", []gopacket.SerializableLayer{ipv4Hdr, testPayload}, []gopacket.SerializableLayer{ethHdrVLAN, vlanTagInner, ipv4Hdr, testPayload}, []uint16{200}},
		{"QinQ", []gopacket.SerializableLayer{ipv4Hdr, testPayload}, []gopacket.SerializableLayer{ethHdrVLAN, vlanTagOuter, vlanTagInner, ipv4Hdr, testPayload}, []uint16{100, 200}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testBuf := gopacket.NewSerializeBuffer()
			err := gopacket.SerializeLayers(testBuf, opts, tt.testLayers...)
			if err != nil {
				t.Errorf("error serializing in gopacket (not our bug!)")
			}

			var want []byte
			if tt.wantLayers == nil {
				want = nil
			} else {
				wantBuf := gopacket.NewSerializeBuffer()
				err := gopacket.SerializeLayers(wantBuf, opts, tt.wantLayers...)
				if err != nil {
					t.Errorf("error serializing in gopacket (not our bug!)")
				}
				want = wantBuf.Bytes()
			}

			testBytes := testBuf.Bytes()
			test := addEthernetHdr(testBytes, BroadcastMac, testMac, etherIPv4Proto, tt.vlanConfig)

			if !bytes.Equal(test, want) {
				t.Errorf("got %v, want %v", test, want)
			}
		})
	}
}
