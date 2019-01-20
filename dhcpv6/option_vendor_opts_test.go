package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptVendorOpts(t *testing.T) {
	optData := []byte("Arista;DCS-7304;01.00;HSH14425148")
	// NOTE: this should be aware of endianness
	expected := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	expected = append(expected, []byte{0, 1, //code
		0, byte(len(optData)), //length
	}...)
	expected = append(expected, optData...)
	expectedOpts := OptVendorOpts{}
	var vendorOpts []Option
	expectedOpts.VendorOpts = append(vendorOpts, &OptionGeneric{OptionCode: 1, OptionData: optData})
	opt, err := ParseOptVendorOpts(expected)
	require.NoError(t, err)
	require.Equal(t, len(expected), opt.Length())
	require.Equal(t, uint32(0xaabbccdd), opt.EnterpriseNumber)
	require.Equal(t, expectedOpts.VendorOpts, opt.VendorOpts)

	shortData := make([]byte, 1)
	_, err = ParseOptVendorOpts(shortData)
	require.Error(t, err)
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
	require.Equal(t, expected, toBytes)
}
