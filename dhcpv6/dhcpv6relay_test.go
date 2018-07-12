package dhcpv6

import (
	"bytes"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDHCPv6Relay(t *testing.T) {
	ll := net.IP{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaa, 0xbb, 0xcc, 0xff, 0xfe, 0xdd, 0xee, 0xff}
	ma := net.IP{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	r := DHCPv6Relay{
		messageType: RELAY_FORW,
		hopCount:    10,
		linkAddr:    ll,
		peerAddr:    ma,
		// options is left empty here for testing purposes, even if it's
		// mandatory to have at least a relay message option
	}
	if mt := r.MessageType(); mt != RELAY_FORW {
		t.Fatalf("Invalid message type. Expected %v, got %v", RELAY_FORW, mt)
	}
	if hc := r.HopCount(); hc != 10 {
		t.Fatalf("Invalid hop count. Expected 10, got %v", hc)
	}
	if la := r.LinkAddr(); !bytes.Equal(la, ll) {
		t.Fatalf("Invalid link address. Expected %v, got %v", ll, la)
	}
	if pa := r.PeerAddr(); !bytes.Equal(pa, ma) {
		t.Fatalf("Invalid peer address. Expected %v, got %v", ma, pa)
	}
	if opts := r.Options(); len(opts) != 0 {
		t.Fatalf("Invalid options. Expected none, got %v", opts)
	}
}

func TestDHCPv6RelaySettersAndGetters(t *testing.T) {
	r := DHCPv6Relay{}
	// set and get message type
	r.SetMessageType(RELAY_REPL)
	if mt := r.MessageType(); mt != RELAY_REPL {
		t.Fatalf("Invalid message type. Expected %v, got %v", RELAY_REPL, mt)
	}
	// set and get hop count
	r.SetHopCount(5)
	if hc := r.HopCount(); hc != 5 {
		t.Fatalf("Invalid hop count. Expected 5, got %v", hc)
	}
	// set and get link address
	r.SetLinkAddr(net.IPv6loopback)
	if la := r.LinkAddr(); !bytes.Equal(la, net.IPv6loopback) {
		t.Fatalf("Invalid link address. Expected %v, got %v", net.IPv6loopback, la)
	}
	// set and get peer address
	r.SetPeerAddr(net.IPv6loopback)
	if pa := r.PeerAddr(); !bytes.Equal(pa, net.IPv6loopback) {
		t.Fatalf("Invalid peer address. Expected %v, got %v", net.IPv6loopback, pa)
	}
	// set and get options
	oneOpt := []Option{&OptRelayMsg{relayMessage: &DHCPv6Message{}}}
	r.SetOptions(oneOpt)
	if opts := r.Options(); len(opts) != 1 || opts[0] != oneOpt[0] {
		t.Fatalf("Invalid options. Expected %v, got %v", oneOpt, opts)
	}
}

func TestDHCPv6RelayToBytes(t *testing.T) {
	expected := []byte{
		12,                                                      // RELAY_FORW
		1,                                                       // hop count
		0xff, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01, // link addr
		0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, // peer addr
		// option relay message
		0, 9, // relay msg
		0, 10, // option length
		// inner dhcp solicit
		1,                // SOLICIT
		0xaa, 0xbb, 0xcc, // transaction ID
		// inner option - elapsed time
		0, 8, // elapsed time
		0, 2, // length
		0, 0,
	}
	r := DHCPv6Relay{
		messageType: RELAY_FORW,
		hopCount:    1,
		linkAddr:    net.IPv6interfacelocalallnodes,
		peerAddr:    net.IPv6linklocalallrouters,
	}
	opt := OptRelayMsg{
		relayMessage: &DHCPv6Message{
			messageType:   SOLICIT,
			transactionID: 0xaabbcc,
			options: []Option{
				&OptElapsedTime{
					ElapsedTime: 0,
				},
			},
		},
	}
	r.AddOption(&opt)
	relayBytes := r.ToBytes()
	if !bytes.Equal(expected, relayBytes) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, relayBytes)
	}
}

func TestGetInnerRelay(t *testing.T) {
	m := DHCPv6Message{}
	r1, err := EncapsulateRelay(&m, RELAY_FORW, net.IPv6linklocalallnodes, net.IPv6interfacelocalallnodes)
	require.NoError(t, err)
	r2, err := EncapsulateRelay(r1, RELAY_FORW, net.IPv6loopback, net.IPv6linklocalallnodes)
	require.NoError(t, err)
	r3, err := EncapsulateRelay(r2, RELAY_FORW, net.IPv6unspecified, net.IPv6linklocalallrouters)
	require.NoError(t, err)

	relay3, ok := r3.(*DHCPv6Relay)
	require.True(t, ok)

	ir, err := relay3.GetInnerRelay()
	require.NoError(t, err)
	relay, ok := ir.(*DHCPv6Relay)
	require.True(t, ok)
	require.Equal(t, relay.HopCount(), uint8(0))
	require.Equal(t, relay.LinkAddr(), net.IPv6linklocalallnodes)
	require.Equal(t, relay.PeerAddr(), net.IPv6interfacelocalallnodes)

	innerPeerAddr, err := relay3.GetInnerPeerAddr()
	require.NoError(t, err)
	require.Equal(t, innerPeerAddr, net.IPv6interfacelocalallnodes)
}

func TestNewRelayRepFromRelayForw(t *testing.T) {
	rf := DHCPv6Relay{}
	rf.SetMessageType(RELAY_FORW)
	rf.SetPeerAddr(net.IPv6linklocalallrouters)
	rf.SetLinkAddr(net.IPv6interfacelocalallnodes)
	oro := OptRelayMsg{}
	s := DHCPv6Message{}
	s.SetMessage(SOLICIT)
	cid := OptClientId{}
	s.AddOption(&cid)
	oro.SetRelayMessage(&s)
	rf.AddOption(&oro)

	a, err := NewAdvertiseFromSolicit(&s)
	require.NoError(t, err)
	rr, err := NewRelayReplFromRelayForw(&rf, a)
	require.NoError(t, err)
	relay := rr.(*DHCPv6Relay)
	require.Equal(t, rr.Type(), RELAY_REPL)
	require.Equal(t, relay.HopCount(), rf.HopCount())
	require.Equal(t, relay.PeerAddr(), rf.PeerAddr())
	require.Equal(t, relay.LinkAddr(), rf.LinkAddr())
	m, err := relay.GetInnerMessage()
	require.NoError(t, err)
	require.Equal(t, m, a)

	rr, err = NewRelayReplFromRelayForw(nil, a)
	require.Error(t, err)
	rr, err = NewRelayReplFromRelayForw(&rf, nil)
	require.Error(t, err)
}
