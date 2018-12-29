package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	sampleVIVCOpt = VIVCIdentifiers{
		VIVCIdentifier{EntID: 9, Data: []byte("CiscoIdentifier")},
		VIVCIdentifier{EntID: 18, Data: []byte("WellfleetIdentifier")},
	}
	sampleVIVCOptRaw = []byte{
		0x0, 0x0, 0x0, 0x9, // enterprise id 9
		0xf, // length
		'C', 'i', 's', 'c', 'o', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r',
		0x0, 0x0, 0x0, 0x12, // enterprise id 18
		0x13, // length
		'W', 'e', 'l', 'l', 'f', 'l', 'e', 'e', 't', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r',
	}
)

func TestOptVIVCInterfaceMethods(t *testing.T) {
	opt := OptVIVC(sampleVIVCOpt...)
	require.Equal(t, OptionVendorIdentifyingVendorClass, opt.Code, "Code")
	require.Equal(t, sampleVIVCOptRaw, opt.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Vendor-Identifying Vendor Class: 9:'CiscoIdentifier', 18:'WellfleetIdentifier'",
		opt.String())
}

func TestParseOptVICO(t *testing.T) {
	options := Options{OptionVendorIdentifyingVendorClass.Code(): sampleVIVCOptRaw}
	o := GetVIVC(options)
	require.Equal(t, sampleVIVCOpt, o)

	// Identifier len too long
	data := make([]byte, len(sampleVIVCOptRaw))
	copy(data, sampleVIVCOptRaw)
	data[4] = 40
	options = Options{OptionVendorIdentifyingVendorClass.Code(): data}
	o = GetVIVC(options)
	require.Nil(t, o, "should get error from bad length")

	// Longer than length
	data[4] = 5
	options = Options{OptionVendorIdentifyingVendorClass.Code(): data[:10]}
	o = GetVIVC(options)
	require.Equal(t, o[0].Data, []byte("Cisco"))

	require.Equal(t, VIVCIdentifiers(nil), GetVIVC(Options{}))
}
