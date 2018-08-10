package dhcpv6

import (
	"fmt"
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
	if vop := opt.VendorOpts(); !bytes.Equal(vop, vendorOpts) {
		t.Fatalf("Invalid VendorOption. Expected %v, got %v", expected, vop)
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

func TestOptVendorOptsSetEnterpriseNumber(t *testing.T) {
	opt := OptVendorOpts{}
	opt.SetEnterpriseNumber(uint32(3062))
	expected := uint32(3062)
	if opt.EnterpriseNumber() != expected {
		t.Fatalf("Invalid SetEnterpriseNumber result. Expected %v, got %v", expected, opt.EnterpriseNumber())
	}
}

func TestOptVendorOptsSetVendorOpts(t *testing.T) {
	opt := OptVendorOpts{}
	opt.SetVendorOpts([]byte("!Arista;DCS-7304;01.00;HSH14425148"))
	expected := []byte("!Arista;DCS-7304;01.00;HSH14425148")
	if !bytes.Equal(opt.VendorOpts(), expected) {
		t.Fatalf("Invalid SetEnterpriseNumber result. Expected %v, got %v", expected, opt.VendorOpts())
	}
}

func TestOptVendorOptsString(t *testing.T) {
	opt := OptVendorOpts{}
	opt.SetEnterpriseNumber(uint32(3062))
	opt.SetVendorOpts([]byte("!Arista;DCS-7304;01.00;HSH14425148"))
	expected := fmt.Sprintf("OptVendorOpts{enterprisenum=%v, vendorOpts=%s}",
		opt.enterpriseNumber, opt.vendorOpts,
	)
	if opt.String() != expected {
		t.Fatalf("Invalid SetEnterpriseNumber result. \nExpected %v \nGot \t %v", expected, opt.String())
	}
}
