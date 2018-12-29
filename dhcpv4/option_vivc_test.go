package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	sampleVIVCOpt = OptVIVC{
		Identifiers: []VIVCIdentifier{
			{EntID: 9, Data: []byte("CiscoIdentifier")},
			{EntID: 18, Data: []byte("WellfleetIdentifier")},
		},
	}
	sampleVIVCOptRaw = []byte{
		byte(OptionVendorIdentifyingVendorClass), 44, // option header
		0x0, 0x0, 0x0, 0x9, // enterprise id 9
		0xf, // length
		'C', 'i', 's', 'c', 'o', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r',
		0x0, 0x0, 0x0, 0x12, // enterprise id 18
		0x13, // length
		'W', 'e', 'l', 'l', 'f', 'l', 'e', 'e', 't', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r',
	}
)

func TestOptVIVCInterfaceMethods(t *testing.T) {
	require.Equal(t, OptionVendorIdentifyingVendorClass, sampleVIVCOpt.Code(), "Code")
	require.Equal(t, 44, sampleVIVCOpt.Length(), "Length")
	require.Equal(t, sampleVIVCOptRaw, sampleVIVCOpt.ToBytes(), "ToBytes")
}

func TestParseOptVICO(t *testing.T) {
	o, err := ParseOptVIVC(sampleVIVCOptRaw[2:])
	require.NoError(t, err)
	require.Equal(t, &sampleVIVCOpt, o)

	// Identifier len too long
	data := make([]byte, len(sampleVIVCOptRaw[2:]))
	copy(data, sampleVIVCOptRaw[2:])
	data[4] = 40
	_, err = ParseOptVIVC(data)
	require.Error(t, err, "should get error from bad length")

	// Longer than length
	data[4] = 5
	o, err = ParseOptVIVC(data[:10])
	require.NoError(t, err)
	require.Equal(t, o.Identifiers[0].Data, []byte("Cisco"))
}

func TestOptVIVCString(t *testing.T) {
	require.Equal(t, "Vendor-Identifying Vendor Class -> 9:'CiscoIdentifier', 18:'WellfleetIdentifier'",
		sampleVIVCOpt.String())
}
