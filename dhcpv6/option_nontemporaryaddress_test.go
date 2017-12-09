package dhcpv6

import (
	"testing"
)

func TestOptIANAParseOptIANA(t *testing.T) {
	data := []byte{
		02,            // advertise
		0, 0x80, 0x8b, // transaction ID
		// IA_NA option
		0, 3, // option code
		0, 40, // option length
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, 0, 0, 0xb2, 0x7a, 0, 0, 0xc0, 0x8a, // options
	}
	opt, err := ParseOptIANA(data)
	if err != nil {
		t.Fatal(err)
	}
	if oLen := opt.Length(); oLen != len(data) {
		t.Fatalf("Invalid IANA option length. Expected %v, got %v", len(data), oLen)
	}
}
