package dhcpv4

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOption(t *testing.T) {
	// Generic
	option := []byte{5, 4, 192, 168, 1, 254} // DNS option
	opt, err := ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	generic := opt.(*OptionGeneric)
	require.Equal(t, OptionNameServer, generic.Code())
	require.Equal(t, []byte{192, 168, 1, 254}, generic.Data)
	require.Equal(t, 4, generic.Length())
	require.Equal(t, "Name Server -> [192 168 1 254]", generic.String())

	// Option subnet mask
	option = []byte{1, 4, 255, 255, 255, 0}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionSubnetMask, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option router
	option = []byte{3, 4, 192, 168, 1, 1}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionRouter, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option domain name server
	option = []byte{6, 4, 192, 168, 1, 1}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionDomainNameServer, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option host name
	option = []byte{12, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionHostName, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option domain name
	option = []byte{15, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionDomainName, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option root path
	option = []byte{17, 4, '/', 'f', 'o', 'o'}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionRootPath, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option broadcast address
	option = []byte{28, 4, 255, 255, 255, 255}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionBroadcastAddress, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option NTP servers
	option = []byte{42, 4, 10, 10, 10, 10}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionNTPServers, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Requested IP address
	option = []byte{50, 4, 1, 2, 3, 4}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionRequestedIPAddress, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Requested IP address lease time
	option = []byte{51, 4, 0, 0, 0, 0}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionIPAddressLeaseTime, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Message type
	option = []byte{53, 1, 1}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionDHCPMessageType, opt.Code(), "Code")
	require.Equal(t, 1, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option server ID
	option = []byte{54, 4, 1, 2, 3, 4}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionServerIdentifier, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Parameter request list
	option = []byte{55, 3, 5, 53, 61}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionParameterRequestList, opt.Code(), "Code")
	require.Equal(t, 3, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option max message size
	option = []byte{57, 2, 1, 2}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionMaximumDHCPMessageSize, opt.Code(), "Code")
	require.Equal(t, 2, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option class identifier
	option = []byte{60, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionClassIdentifier, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option TFTP server name
	option = []byte{66, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionTFTPServerName, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option Bootfile name
	option = []byte{67, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't'}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionBootfileName, opt.Code(), "Code")
	require.Equal(t, 9, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option user class information
	option = []byte{77, 5, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionUserClassInformation, opt.Code(), "Code")
	require.Equal(t, 5, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option relay agent information
	option = []byte{82, 6, 1, 4, 129, 168, 0, 1}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionRelayAgentInformation, opt.Code(), "Code")
	require.Equal(t, 6, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")

	// Option client system architecture type option
	option = []byte{93, 4, 't', 'e', 's', 't'}
	opt, err = ParseOption(OptionCode(option[0]), option[2:])
	require.NoError(t, err)
	require.Equal(t, OptionClientSystemArchitectureType, opt.Code(), "Code")
	require.Equal(t, 4, opt.Length(), "Length")
	require.Equal(t, option, opt.ToBytes(), "ToBytes")
}

func TestOptionsUnmarshal(t *testing.T) {
	for i, tt := range []struct {
		input     []byte
		want      Options
		wantError bool
	}{
		{
			// Buffer missing data.
			input: []byte{
				3 /* key */, 3 /* length */, 1,
			},
			wantError: true,
		},
		{
			input: []byte{
				// This may look too long, but 0 is padding.
				// The issue here is the missing OptionEnd.
				3, 3, 0, 0, 0, 0, 0, 0, 0,
			},
			wantError: true,
		},
		{
			// Only OptionPad and OptionEnd can stand on their own
			// without a length field. So this is too short.
			input: []byte{
				3,
			},
			wantError: true,
		},
		{
			// Option present after the End is a nono.
			input:     []byte{byte(OptionEnd), 3},
			wantError: true,
		},
		{
			input: []byte{byte(OptionEnd)},
			want:  Options{},
		},
		{
			input: []byte{
				3, 2, 5, 6,
				byte(OptionEnd),
			},
			want: Options{
				&OptionGeneric{
					OptionCode: 3,
					Data:       []byte{5, 6},
				},
			},
		},
		{
			// Test RFC 3396.
			input: append(
				append([]byte{3, math.MaxUint8}, bytes.Repeat([]byte{10}, math.MaxUint8)...),
				3, 5, 10, 10, 10, 10, 10,
				byte(OptionEnd),
			),
			want: Options{
				&OptionGeneric{
					OptionCode: 3,
					Data:       bytes.Repeat([]byte{10}, math.MaxUint8+5),
				},
			},
		},
		{
			input: []byte{
				10, 2, 255, 254,
				11, 3, 5, 5, 5,
				byte(OptionEnd),
			},
			want: Options{
				&OptionGeneric{
					OptionCode: 10,
					Data:       []byte{255, 254},
				},
				&OptionGeneric{
					OptionCode: 11,
					Data:       []byte{5, 5, 5},
				},
			},
		},
		{
			input: append(
				append([]byte{10, 2, 255, 254}, bytes.Repeat([]byte{byte(OptionPad)}, 255)...),
				byte(OptionEnd),
			),
			want: Options{
				&OptionGeneric{
					OptionCode: 10,
					Data:       []byte{255, 254},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test %02d", i), func(t *testing.T) {
			opt, err := OptionsFromBytesWithParser(tt.input, ParseOptionGeneric, true)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, opt, tt.want)
			}
		})
	}
}
