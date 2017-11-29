package dhcpv6

import (
	"bytes"
	"github.com/insomniacslk/dhcp/dhcpv6/options"
	"testing"
)

func TestBytesToTransactionID(t *testing.T) {
	// only the first three bytes should be used
	tid, err := BytesToTransactionID([]byte{0x11, 0x22, 0x33, 0xaa})
	if err != nil {
		t.Fatal(err)
	}
	if tid == nil {
		t.Fatal("Invalid Transaction ID. Should not be nil")
	}
	if *tid != 0x112233 {
		t.Fatalf("Invalid Transaction ID. Expected 0x%x, got 0x%x", 0x112233, *tid)
	}
}

func TestBytesToTransactionIDShortData(t *testing.T) {
	// short sequence, less than three bytes
	tid, err := BytesToTransactionID([]byte{0x11, 0x22})
	if err == nil {
		t.Fatal("Expected non-nil error, got nil instead")
	}
	if tid != nil {
		t.Errorf("Expected nil Transaction ID, got %v instead", *tid)
	}
}

func TestGenerateTransactionID(t *testing.T) {
	tid, err := GenerateTransactionID()
	if err != nil {
		t.Fatal(err)
	}
	if tid == nil {
		t.Fatal("Expected non-nil Transaction ID, got nil instead")
	}
	if *tid > 0xffffff {
		// TODO this should be better tested by mocking the random generator
		t.Fatalf("Invalid Transaction ID: should be smaller than 0xffffff. Got 0x%x instead", *tid)
	}
}

func TestNew(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if d == nil {
		t.Fatal("Expected non-nil DHCPv6, got nil instead")
	}
	if d.message != SOLICIT {
		t.Fatalf("Invalid message type. Expected %v, got %v", SOLICIT, d.message)
	}
	if d.transactionID == 0 {
		t.Fatal("Invalid Transaction ID, expected non-zero, got zero")
	}
	if len(d.options) != 0 {
		t.Fatalf("Invalid options: expected none, got %v", len(d.options))
	}
}

func TestSettersAndGetters(t *testing.T) {
	d := DHCPv6{}
	// Message
	d.SetMessage(SOLICIT)
	msg := d.Message()
	if msg != SOLICIT {
		t.Fatalf("Invalid Message. Expected %v, got %v", SOLICIT, msg)
	}
	d.SetMessage(ADVERTISE)
	msg = d.Message()
	if msg != ADVERTISE {
		t.Fatalf("Invalid Message. Expected %v, got %v", ADVERTISE, msg)
	}
	// TransactionID
	d.SetTransactionID(12345)
	tid := d.TransactionID()
	if tid != 12345 {
		t.Fatalf("Invalid Transaction ID. Expected %v, got %v", 12345, tid)
	}
	// Options
	opts := d.Options()
	if len(opts) != 0 {
		t.Fatalf("Invalid Options. Expected empty array, got %v", opts)
	}
	opt := options.OptionGeneric{OptionCode: 0, OptionData: []byte{}}
	d.SetOptions([]options.Option{&opt})
	opts = d.Options()
	if len(opts) != 1 {
		t.Fatalf("Invalid Options. Expected one-element array, got %v", len(opts))
	}
	if _, ok := opts[0].(*options.OptionGeneric); !ok {
		t.Fatalf("Invalid Options. Expected one OptionGeneric, got %v", opts[0])
	}
}

func TestAddOption(t *testing.T) {
	d := DHCPv6{}
	opts := d.Options()
	if len(opts) != 0 {
		t.Fatalf("Invalid Options. Expected empty array, got %v", opts)
	}
	opt := options.OptionGeneric{OptionCode: 0, OptionData: []byte{}}
	d.AddOption(&opt)
	opts = d.Options()
	if len(opts) != 1 {
		t.Fatalf("Invalid Options. Expected one-element array, got %v", len(opts))
	}
	if _, ok := opts[0].(*options.OptionGeneric); !ok {
		t.Fatalf("Invalid Options. Expected one OptionGeneric, got %v", opts[0])
	}
}

func TestToBytes(t *testing.T) {
	d := DHCPv6{}
	d.SetMessage(SOLICIT)
	d.SetTransactionID(0xabcdef)
	opt := options.OptionGeneric{OptionCode: 0, OptionData: []byte{}}
	d.AddOption(&opt)
	toBytes := d.ToBytes()
	expected := []byte{01, 0xab, 0xcd, 0xef, 0x00, 0x00, 0x00, 0x00}
	if !bytes.Equal(toBytes, expected) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, toBytes)
	}
}

func TestFromAndToBytes(t *testing.T) {
	expected := []byte{01, 0xab, 0xcd, 0xef, 0x00, 0x00, 0x00, 0x00}
	d, err := FromBytes(expected)
	if err != nil {
		t.Fatal(err)
	}
	toBytes := d.ToBytes()
	if !bytes.Equal(toBytes, expected) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, toBytes)
	}
}

// TODO test NewSolicit
//      test String and Summary
