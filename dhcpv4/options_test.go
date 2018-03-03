package dhcpv4

import (
	"bytes"
	"testing"
)

func TestParseOption(t *testing.T) {
	option := []byte{5, 4, 192, 168, 1, 254} // DNS option
	opt, err := ParseOption(option)
	if err != nil {
		t.Fatal(err)
	}
	if opt.Code != OptionNameServer {
		t.Fatalf("Invalid option code. Expected 5, got %v", opt.Code)
	}
	if !bytes.Equal(opt.Data, option[2:]) {
		t.Fatalf("Invalid option data. Expected %v, got %v", option[2:], opt.Data)
	}
}

func TestParseOptionPad(t *testing.T) {
	option := []byte{0}
	opt, err := ParseOption(option)
	if err != nil {
		t.Fatal(err)
	}
	if opt.Code != OptionPad {
		t.Fatalf("Invalid option code. Expected %v, got %v", OptionPad, opt.Code)
	}
	if len(opt.Data) != 0 {
		t.Fatalf("Invalid option data. Expected empty slice, got %v", opt.Data)
	}
}

func TestParseOptionZeroLength(t *testing.T) {
	option := []byte{}
	_, err := ParseOption(option)
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}

func TestParseOptionShortOption(t *testing.T) {
	option := []byte{53, 1}
	_, err := ParseOption(option)
	if err == nil {
		t.Fatal(err)
	}
}

func TestOptionsFromBytes(t *testing.T) {
	options := []byte{
		99, 130, 83, 99, // Magic Cookie
		5, 4, 192, 168, 1, 1, // DNS
		255,     // end
		0, 0, 0, //padding
	}
	opts, err := OptionsFromBytesWithMagicCookie(options)
	if err != nil {
		t.Fatal(err)
	}
	// each padding byte counts as an option. Magic Cookie doesn't add up
	if len(opts) != 5 {
		t.Fatalf("Invalid options length. Expected 5, got %v", len(opts))
	}
	if opts[0].Code != OptionNameServer {
		t.Fatalf("Invalid option code. Expected %v, got %v", OptionNameServer, opts[0].Code)
	}
	if !bytes.Equal(opts[0].Data, options[6:10]) {
		t.Fatalf("Invalid option data. Expected %v, got %v", options[6:10], opts[0].Data)
	}
	if opts[1].Code != OptionEnd {
		t.Fatalf("Invalid option code. Expected %v, got %v", OptionEnd, opts[1].Code)
	}
	if opts[2].Code != OptionPad {
		t.Fatalf("Invalid option code. Expected %v, got %v", OptionPad, opts[2].Code)
	}
}

func TestOptionsFromBytesZeroLength(t *testing.T) {
	options := []byte{}
	_, err := OptionsFromBytesWithMagicCookie(options)
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}

func TestOptionsFromBytesBadMagicCookie(t *testing.T) {
	options := []byte{1, 2, 3, 4}
	_, err := OptionsFromBytesWithMagicCookie(options)
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}

func TestOptionsToBytesWithMagicCookie(t *testing.T) {
	originalOptions := []byte{
		99, 130, 83, 99, // Magic Cookie
		5, 4, 192, 168, 1, 1, // DNS
		255,     // end
		0, 0, 0, //padding
	}
	options, err := OptionsFromBytesWithMagicCookie(originalOptions)
	if err != nil {
		t.Fatal(err)
	}
	finalOptions := OptionsToBytesWithMagicCookie(options)
	if !bytes.Equal(originalOptions, finalOptions) {
		t.Fatalf("Invalid options. Expected %v, got %v", originalOptions, finalOptions)
	}
}

func TestOptionsToBytesWithMagicCookieEmpty(t *testing.T) {
	originalOptions := []byte{99, 130, 83, 99}
	options, err := OptionsFromBytesWithMagicCookie(originalOptions)
	if err != nil {
		t.Fatal(err)
	}
	finalOptions := OptionsToBytesWithMagicCookie(options)
	if !bytes.Equal(originalOptions, finalOptions) {
		t.Fatalf("Invalid options. Expected %v, got %v", originalOptions, finalOptions)
	}
}

func TestOptionsToStringPad(t *testing.T) {
	option := []byte{0}
	opt, err := ParseOption(option)
	if err != nil {
		t.Fatal(err)
	}
	stropt := opt.String()
	if stropt != "Pad -> []" {
		t.Fatalf("Invalid string representation: %v", stropt)
	}
}

func TestOptionsToStringDHCPMessageType(t *testing.T) {
	option := []byte{53, 1, 5}
	opt, err := ParseOption(option)
	if err != nil {
		t.Fatal(err)
	}
	stropt := opt.String()
	if stropt != "DHCP Message Type -> [5]" {
		t.Fatalf("Invalid string representation: %v", stropt)
	}
}

func TestBSDPOptionToString(t *testing.T) {
	// Parse message type
	option := Option{
		Code: BSDPOptionMessageType,
		Data: []byte{BSDPMessageTypeList},
	}
	stropt := option.BSDPString()
	AssertEqual(t, stropt, "BSDP Message Type -> [1]", "BSDP string representation")

	// Parse failure
	option = Option{
		Code: OptionCode(12), // invalid BSDP Opcode
		Data: []byte{1, 2, 3},
	}
	stropt = option.BSDPString()
	AssertEqual(t, stropt, "Unknown -> [1 2 3]", "BSDP string representation")
}
