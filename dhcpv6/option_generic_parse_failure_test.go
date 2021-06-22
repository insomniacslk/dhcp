package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptGenericParseFailure(t *testing.T) {
	buf := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // EnterpriseNumber
		12, 0, // length is little-endian - should be big-endian
		'd', 's', 'l', 'f', 'o', 'r', 'u', 'm', '.', 'o', 'r', 'g',
	}
	opt, _ := ParseOption(OptionVendorClass, buf)
	require.IsType(t, &OptionGenericParseFailure{}, opt)
	failedOpt, ok := opt.(*OptionGenericParseFailure)
	require.True(t, ok)
	require.Contains(
		t,
		failedOpt.Error.Error(),
		"buffer too short",
		"Error() should return the original parser error",
	)
	require.Contains(
		t,
		failedOpt.String(),
		"GenericParseFailure(Vendor Class)",
		"String() should include the Option Code",
	)
	require.Contains(
		t,
		failedOpt.String(),
		"enterprisenum=2864434397",
		"String() should include the details of the failed option",
	)
}
