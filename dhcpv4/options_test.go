package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOption(t *testing.T) {
	// Generic
	option := []byte{5, 4, 192, 168, 1, 254} // DNS option
	opt, err := ParseOption(option)
	require.NoError(t, err)
	generic := opt.(*OptionGeneric)
	require.Equal(t, OptionNameServer, generic.Code())
	require.Equal(t, []byte{192, 168, 1, 254}, generic.Data)
	require.Equal(t, 4, generic.Length())
	require.Equal(t, "Name Server -> [192 168 1 254]", generic.String())

	// Option subnet mask
	option = []byte{1, 4, 255, 255, 255, 0}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionSubnetMask, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option router
	option = []byte{3, 4, 192, 168, 1, 1}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionRouter, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option domain name server
	option = []byte{6, 4, 192, 168, 1, 1}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionDomainNameServer, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option host name
	option = []byte{12, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionHostName, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option domain name
	option = []byte{15, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionDomainName, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option NTP servers
	option = []byte{42, 4, 10, 10, 10, 10}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionNTPServers, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Requested IP address
	option = []byte{50, 4, 1, 2, 3, 4}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionRequestedIPAddress, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Message type
	option = []byte{53, 1, 1}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionDHCPMessageType, opt.Code(), "Code")
	require.Equal(t, 1, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option server ID
	option = []byte{54, 4, 1, 2, 3, 4}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionServerIdentifier, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Parameter request list
	option = []byte{55, 3, 5, 53, 61}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionParameterRequestList, opt.Code(), "Code")
	require.Equal(t, 3, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option max message size
	option = []byte{57, 2, 1, 2}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionMaximumDHCPMessageSize, opt.Code(), "Code")
	require.Equal(t, 2, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option class identifier
	option = []byte{60, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionClassIdentifier, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option TFTP server name
	option = []byte{66, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionTFTPServerName, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option Bootfile name
	option = []byte{67, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't'}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionBootfileName, opt.Code(), "Code")
	require.Equal(t, 9, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option user class information
	option = []byte{77, 5, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionUserClassInformation, opt.Code(), "Code")
	require.Equal(t, 5, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option client system architecture type option
	option = []byte{93, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionClientSystemArchitectureType, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	option = []byte{17, 4, '/', 'f', 'o', 'o'}
	opt, err = ParseOption(option)
	require.NoError(t, err)
	require.Equal(t, OptionRootPath, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")
}

func TestParseOptionZeroLength(t *testing.T) {
	option := []byte{}
	_, err := ParseOption(option)
	require.Error(t, err, "should get error from zero-length options")
}

func TestParseOptionShortOption(t *testing.T) {
	option := []byte{53, 1}
	_, err := ParseOption(option)
	require.Error(t, err, "should get error from short options")
}

func TestOptionsFromBytes(t *testing.T) {
	options := []byte{
		99, 130, 83, 99, // Magic Cookie
		5, 4, 192, 168, 1, 1, // DNS
		255,     // end
		0, 0, 0, //padding
	}
	opts, err := OptionsFromBytes(options)
	require.NoError(t, err)
	require.Equal(t, 2, len(opts))
	require.Equal(t, opts[0].(*OptionGeneric), &OptionGeneric{OptionCode: OptionNameServer, Data: []byte{192, 168, 1, 1}})
	require.Equal(t, opts[1].(*OptionGeneric), &OptionGeneric{OptionCode: OptionEnd})
}

func TestOptionsFromBytesZeroLength(t *testing.T) {
	options := []byte{}
	_, err := OptionsFromBytes(options)
	require.Error(t, err)
}

func TestOptionsFromBytesBadMagicCookie(t *testing.T) {
	options := []byte{1, 2, 3, 4}
	_, err := OptionsFromBytes(options)
	require.Error(t, err)
}

func TestOptionsFromBytesShortOption(t *testing.T) {
	options := []byte{
		99, 130, 83, 99, // Magic Cookie
		5, 4, 192, 168, // DNS
	}
	_, err := OptionsFromBytes(options)
	require.Error(t, err)
}
