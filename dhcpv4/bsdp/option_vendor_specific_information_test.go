package bsdp

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/require"
)

func TestOptVendorSpecificInformationInterfaceMethods(t *testing.T) {
	messageTypeOpt := &OptMessageType{MessageTypeList}
	versionOpt := &OptVersion{Version1_1}
	o := &OptVendorSpecificInformation{[]dhcpv4.Option{messageTypeOpt, versionOpt}}
	require.Equal(t, dhcpv4.OptionVendorSpecificInformation, o.Code(), "Code")
	require.Equal(t, 2+messageTypeOpt.Length()+2+versionOpt.Length(), o.Length(), "Length")

	expectedBytes := []byte{
		43,      // code
		7,       // length
		1, 1, 1, // List option
		2, 2, 1, 1, // Version option
	}
	o = &OptVendorSpecificInformation{
		[]dhcpv4.Option{
			&OptMessageType{MessageTypeList},
			&OptVersion{Version1_1},
		},
	}
	require.Equal(t, expectedBytes, o.ToBytes(), "ToBytes")
}

func TestParseOptVendorSpecificInformation(t *testing.T) {
	var (
		o   *OptVendorSpecificInformation
		err error
	)
	o, err = ParseOptVendorSpecificInformation([]byte{})
	require.Error(t, err, "empty byte stream")

	o, err = ParseOptVendorSpecificInformation([]byte{1, 2})
	require.Error(t, err, "short byte stream")

	o, err = ParseOptVendorSpecificInformation([]byte{53, 2, 1, 1})
	require.Error(t, err, "wrong option code")

	// Good byte stream
	data := []byte{
		43,      // code
		7,       // length
		1, 1, 1, // List option
		2, 2, 1, 1, // Version option
	}
	o, err = ParseOptVendorSpecificInformation(data)
	require.NoError(t, err)
	expected := &OptVendorSpecificInformation{
		[]dhcpv4.Option{
			&OptMessageType{MessageTypeList},
			&OptVersion{Version1_1},
		},
	}
	require.Equal(t, 2, len(o.Options), "number of parsed suboptions")
	require.Equal(t, expected.Options[0].Code(), o.Options[0].Code())
	require.Equal(t, expected.Options[1].Code(), o.Options[1].Code())

	// Short byte stream (length and data mismatch)
	data = []byte{
		43,      // code
		7,       // length
		1, 1, 1, // List option
		2, 2, 1, // Version option
	}
	o, err = ParseOptVendorSpecificInformation(data)
	require.Error(t, err)
}

func TestOptVendorSpecificInformationString(t *testing.T) {
	o := &OptVendorSpecificInformation{
		[]dhcpv4.Option{
			&OptMessageType{MessageTypeList},
			&OptVersion{Version1_1},
		},
	}
	expectedString := "Vendor Specific Information ->\n  BSDP Message Type -> LIST\n  BSDP Version -> 1.1"
	require.Equal(t, expectedString, o.String())

	// Test more complicated string - sub options of sub options.
	o = &OptVendorSpecificInformation{
		[]dhcpv4.Option{
			&OptMessageType{MessageTypeList},
			&OptBootImageList{
				[]BootImage{
					BootImage{
						ID: BootImageID{
							IsInstall: false,
							ImageType: BootImageTypeMacOSX,
							Index:     1001,
						},
						Name: "bsdp-1",
					},
					BootImage{
						ID: BootImageID{
							IsInstall: true,
							ImageType: BootImageTypeMacOS9,
							Index:     9009,
						},
						Name: "bsdp-2",
					},
				},
			},
		},
	}
	expectedString = "Vendor Specific Information ->\n" +
		"  BSDP Message Type -> LIST\n" +
		"  BSDP Boot Image List ->\n" +
		"    bsdp-1 [1001] uninstallable macOS image\n" +
		"    bsdp-2 [9009] installable macOS 9 image"
	require.Equal(t, expectedString, o.String())
}

func TestOptVendorSpecificInformationGetOptions(t *testing.T) {
	// No option
	o := &OptVendorSpecificInformation{
		[]dhcpv4.Option{
			&OptMessageType{MessageTypeList},
			&OptVersion{Version1_1},
		},
	}
	foundOpts := o.GetOptions(OptionBootImageList)
	require.Empty(t, foundOpts, "should not get any options")

	// One option
	o = &OptVendorSpecificInformation{
		[]dhcpv4.Option{
			&OptMessageType{MessageTypeList},
			&OptVersion{Version1_1},
		},
	}
	foundOpts = o.GetOptions(OptionMessageType)
	require.Equal(t, 1, len(foundOpts), "should only get one option")
	require.Equal(t, MessageTypeList, foundOpts[0].(*OptMessageType).Type)

	// Multiple options
	o = &OptVendorSpecificInformation{
		[]dhcpv4.Option{
			&OptMessageType{MessageTypeList},
			&OptVersion{Version1_1},
			&OptVersion{Version1_0},
		},
	}
	foundOpts = o.GetOptions(OptionVersion)
	require.Equal(t, 2, len(foundOpts), "should get two options")
	require.Equal(t, Version1_1, foundOpts[0].(*OptVersion).Version)
	require.Equal(t, Version1_0, foundOpts[1].(*OptVersion).Version)
}

func TestOptVendorSpecificInformationGetOneOption(t *testing.T) {
	// No option
	o := &OptVendorSpecificInformation{
		[]dhcpv4.Option{
			&OptMessageType{MessageTypeList},
			&OptVersion{Version1_1},
		},
	}
	foundOpt := o.GetOneOption(OptionBootImageList)
	require.Nil(t, foundOpt, "should not get options")

	// One option
	o = &OptVendorSpecificInformation{
		[]dhcpv4.Option{
			&OptMessageType{MessageTypeList},
			&OptVersion{Version1_1},
		},
	}
	foundOpt = o.GetOneOption(OptionMessageType)
	require.Equal(t, MessageTypeList, foundOpt.(*OptMessageType).Type)

	// Multiple options
	o = &OptVendorSpecificInformation{
		[]dhcpv4.Option{
			&OptMessageType{MessageTypeList},
			&OptVersion{Version1_1},
			&OptVersion{Version1_0},
		},
	}
	foundOpt = o.GetOneOption(OptionVersion)
	require.Equal(t, Version1_1, foundOpt.(*OptVersion).Version)
}
