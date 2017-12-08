package dhcpv6

import (
	"reflect"
	"testing"
)

func TestRelayMsgParseOptRelayMsg(t *testing.T) {
	opt, err := ParseOptRelayMsg([]byte{
		1,                // SOLICIT
		0xaa, 0xbb, 0xcc, // transaction ID
		0, 8, // option: elapsed time
		0, 2, // option length
		0, 0, // option value
	})
	if err != nil {
		t.Fatal(err)
	}
	if code := opt.Code(); code != OPTION_RELAY_MSG {
		t.Fatalf("Invalid option code. Expected OPTION_RELAY_MSG (%v), got %v",
			OPTION_RELAY_MSG, code,
		)
	}
}

func TestRelayMsgOptionsFromBytes(t *testing.T) {
	opts, err := OptionsFromBytes([]byte{
		0, 9, // option: relay message
		0, 10, // relayed message length
		1,                // SOLICIT
		0xaa, 0xbb, 0xcc, // transaction ID
		0, 8, // option: elapsed time
		0, 2, // option length
		0, 0, // option value
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts) != 1 {
		t.Fatalf("Invalid number of options. Expected 1, got %v", len(opts))
	}
	opt := opts[0]
	if code := opt.Code(); code != OPTION_RELAY_MSG {
		t.Fatalf("Invalid option code. Expected OPTION_RELAY_MSG (%v), got %v",
			OPTION_RELAY_MSG, code,
		)
	}
}

func TestRelayMsgParseOptRelayMsgSingleEncapsulation(t *testing.T) {
	d, err := FromBytes([]byte{
		12,                                             // RELAY-FORW
		0,                                              // hop count
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // linkAddr
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, // peerAddr
		0, 9, // option: relay message
		0, 10, // relayed message length
		1,                // SOLICIT
		0xaa, 0xbb, 0xcc, // transaction ID
		0, 8, // option: elapsed time
		0, 2, // option length
		0x11, 0x22, // option value
	})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := d.(*DHCPv6Relay)
	if !ok {
		t.Fatalf("Invalid DHCPv6 type. Expected DHCPv6Relay, got %v",
			reflect.TypeOf(d),
		)
	}
	if mType := r.Type(); mType != RELAY_FORW {
		t.Fatalf("Invalid messge type for relay. Expected %v, got %v", RELAY_FORW, mType)
	}
	if len(r.options) != 1 {
		t.Fatalf("Invalid number of options. Expected 1, got %v", len(r.options))
	}
	if code := r.options[0].Code(); code != OPTION_RELAY_MSG {
		t.Fatalf("Invalid option code. Expected OPTION_RELAY_MSG (%v), got %v",
			OPTION_RELAY_MSG, code,
		)
	}
	opt := r.options[0]
	ro, ok := opt.(*OptRelayMsg)
	if !ok {
		t.Fatalf("Invalid option type. Expected OptRelayMsg, got %v",
			reflect.TypeOf(ro),
		)
	}
	innerDHCP, ok := ro.RelayMessage().(*DHCPv6Message)
	if !ok {
		t.Fatalf("Invalid relay message type. Expected DHCPv6Message, got %v",
			reflect.TypeOf(innerDHCP),
		)
	}
	if dType := innerDHCP.Type(); dType != SOLICIT {
		t.Fatal("Invalid inner DHCP type. Expected SOLICIT (%v), got %v",
			SOLICIT, dType,
		)
	}
	if tID := innerDHCP.TransactionID(); tID != 0xaabbcc {
		t.Fatal("Invalid inner DHCP transaction ID. Expected 0xaabbcc, got %v", tID)
	}
	if len(innerDHCP.options) != 1 {
		t.Fatal("Invalid inner DHCP options length. Expected 1, got %v", len(innerDHCP.options))
	}
	innerOpt := innerDHCP.options[0]
	eto, ok := innerOpt.(*OptElapsedTime)
	if !ok {
		t.Fatal("Invalid inner option type. Expected OptElapsedTime, got %v",
			reflect.TypeOf(innerOpt),
		)
	}
	if eTime := eto.ElapsedTime(); eTime != 0x1122 {
		t.Fatal("Invalid elapsed time. Expected 0x1122, got 0x%04x", eTime)
	}
}

func TestSample(t *testing.T) {
	// Nested relay message. This test only checks if it parses correctly, but
	// could/should be extended to check all the fields like done in other tests
	buf := []byte{
		12,                                             // relay
		1,                                              // hop count
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // linkAddr
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // peerAddr
		// relay msg
		0, 9, // opt relay msg
		0, 66, // opt len
		// relay fwd
		12,                                             // relay
		0,                                              // hop count
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // linkAddr
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // peerAddr
		// opt interface ID
		0, 18, // interface id
		0, 6, // opt len
		0xba, 0xbe, 0xb1, 0xb0, 0xbe, 0xbe, // opt value
		// relay msg
		0, 9, // relay msg
		0, 18, // msg len
		// dhcpv6 msg
		1,                // solicit
		0xaa, 0xbb, 0xcc, // transaction ID
		// client ID
		0, 1, // opt client id
		0, 10, // opt len
		0, 3, // duid type
		0, 1, // hw type
		5, 6, 7, 8, 9, 10, // duid value
	}
	_, err := FromBytes(buf)
	if err != nil {
		t.Fatal(err)
	}
}
