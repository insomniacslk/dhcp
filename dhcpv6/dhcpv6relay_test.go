package dhcpv6

import (
	"bytes"
	"net"
	"testing"
)

func TestDHCPv6Relay(t *testing.T) {
	ll := net.IP{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaa, 0xbb, 0xcc, 0xff, 0xfe, 0xdd, 0xee, 0xff}
	ma := net.IP{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	d := DHCPv6Relay{
		messageType: RELAY_FORW,
		hopCount:    10,
		linkAddr:    ll,
		peerAddr:    ma,
		// options is left empty here for testing purposes, even if it's
		// mandatory to have at least a relay message option
	}
	if mt := d.MessageType(); mt != RELAY_FORW {
		t.Fatalf("Invalid message type. Expected %v, got %v", RELAY_FORW, mt)
	}
	if hc := d.HopCount(); hc != 10 {
		t.Fatalf("Invalid hop count. Expected 10, got %v", hc)
	}
	if la := d.LinkAddr(); !bytes.Equal(la, ll) {
		t.Fatalf("Invalid link address. Expected %v, got %v", ll, la)
	}
	if pa := d.PeerAddr(); !bytes.Equal(pa, ma) {
		t.Fatalf("Invalid peer address. Expected %v, got %v", ma, pa)
	}
	if opts := d.Options(); len(opts) != 0 {
		t.Fatalf("Invalid options. Expected none, got %v", opts)
	}
}
