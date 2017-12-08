package dhcpv6

import (
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
