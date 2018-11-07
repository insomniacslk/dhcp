package dhcpv6

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestOptVendorOpts(t *testing.T) {
	optData := []byte("Arista;DCS-7304;01.00;HSH14425148")
	expected := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	expected = append(expected, []byte{0, 1, //code
		0, byte(len(optData)), //length
	}...)
	expected = append(expected, optData...)
	expectedOpts := OptVendorOpts{}
	var vendorOpts []Option
	expectedOpts.VendorOpts = append(vendorOpts, &OptionGeneric{OptionCode: 1, OptionData: optData})
	opt, err := ParseOptVendorOpts(expected)
	if err != nil {
		t.Fatal(err)
	}

	if optLen := opt.Length(); optLen != len(expected) {
		t.Fatalf("Invalid length. Expected %v, got %v", len(expected), optLen)
	}
	if en := opt.EnterpriseNumber; en != 0xaabbccdd {
		t.Fatalf("Invalid Enterprise Number. Expected 0xaabbccdd, got %v", en)
	}
	if !reflect.DeepEqual(opt.VendorOpts, expectedOpts.VendorOpts) {
		t.Fatalf("Invalid Vendor Option Data. Expected %v, got %v", expected, expectedOpts.VendorOpts)
	}

	shortData := make([]byte, 1)
	opt, err = ParseOptVendorOpts(shortData)
	if err == nil {
		t.Fatalf("Short data (<4 bytes) did not cause an error when it should have")
	}

}

func TestOptVendorOptsToBytes(t *testing.T) {
	optData := []byte("Arista;DCS-7304;01.00;HSH14425148")
	var opts []Option
	opts = append(opts, &OptionGeneric{OptionCode: 1, OptionData: optData})

	var expected []byte
	expected = append(expected, []byte{0, 17, // VendorOption Code 17
		0, byte(len(optData) + 8), // Length of optionData + 4 (code & length of sub-option) + 4 for EnterpriseNumber Length
		0, 0, 0, 0, // EnterpriseNumber
		0, 1, // Sub-Option code from vendor
		0, byte(len(optData))}...) // Length of optionData only
	expected = append(expected, optData...)

	opt := OptVendorOpts{
		EnterpriseNumber: 0000,
		VendorOpts:       opts,
	}
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, expected) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, toBytes)
	}
}

func TestVendParseOption(t *testing.T) {
	var buf []byte
	buf = append(buf, []byte{00, 1, 00, 33}...)
	buf = append(buf, []byte("Arista;DCS-7304;01.00;HSH14425148")...)

	expected := &OptionGeneric{OptionCode: 1, OptionData: []byte("Arista;DCS-7304;01.00;HSH14425148")}
	opt, err := vendParseOption(buf)
	if err != nil {
		fmt.Println(err)
	}
	if !reflect.DeepEqual(opt, expected) {
		t.Fatalf("Invalid Vendor Parse Option result. Expected %v, got %v", expected, opt)
	}


	shortData := make([]byte, 1) // data length too small
	opt, err = vendParseOption(shortData)
	if err == nil {
		t.Fatalf("Short data (<4 bytes) did not cause an error when it should have")
	}

	shortData = []byte{0, 0, 0, 0} // missing actual vendor data.
	opt, err = vendParseOption(shortData)
	if err == nil {
		t.Fatalf("Missing VendorData option. An error should have been returned but wasn't")
	}

	shortData = []byte{0, 0,
		0, 4, // declared length
		0} // data starts here, length of 1
	opt, err = vendParseOption(shortData)
	if err == nil {
		t.Fatalf("Declared length does not match actual data length. An error should have been returned but wasn't")
	}

}
