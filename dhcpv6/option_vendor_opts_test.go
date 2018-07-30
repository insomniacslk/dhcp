package dhcpv6

import (
	"bytes"
	"testing"
)

func TestOptVendorOpts(t *testing.T) {
	expected := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	vendorOpts := []byte("!Arista;DCS-7304;01.00;HSH14425148")
	expected = append(expected, vendorOpts...)
	opt, err := ParseOptVendorOpts(expected)
	if err != nil {
		t.Fatal(err)
	}
	if optLen := opt.Length(); optLen != len(expected) {
		t.Fatalf("Invalid length. Expected %v, got %v", len(expected), optLen)
	}
	if en := opt.EnterpriseNumber(); en != 0xaabbccdd {
		t.Fatalf("Invalid Enterprise Number. Expected 0xaabbccdd, got %v", en)
	}
	if rid := opt.VendorOpts(); !bytes.Equal(rid, vendorOpts) {
		t.Fatalf("Invalid Remote ID. Expected %v, got %v", expected, rid)
	}
}

func TestOptVendorOptsToBytes(t *testing.T) {
	vendOpts := []byte("Arista;DCS-7304;01.00;HSH14425148")
	expected := []byte{00, 17, 00, byte(len(vendOpts) + 4), 00, 00, 00, 00}
	expected = append(expected, vendOpts...)
	opt := OptVendorOpts{
		vendorOpts: vendOpts,
	}
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, expected) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, toBytes)
	}
}
