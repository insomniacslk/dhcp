// +build darwin

package bsdp

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/require"
)

func TestParseBootImageSingleBootImage(t *testing.T) {
	input := []byte{
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                             // len(Name)
		98, 115, 100, 112, 45, 50, 49, // byte-encoding of Name
	}
	bs, err := ParseBootImagesFromOption(input)
	require.NoError(t, err)
	require.Equal(t, len(bs), 1, "parsing single boot image should return 1")
	b := bs[0]
	expectedBootImageID := BootImageID{
		IsInstall: false,
		ImageType: BootImageTypeMacOSX,
		Index:     0x1010,
	}
	require.Equal(t, expectedBootImageID, b.ID)
	require.Equal(t, b.Name, "bsdp-21")
}

func TestParseBootImageMultipleBootImage(t *testing.T) {
	input := []byte{
		// boot image 1
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                             // len(Name)
		98, 115, 100, 112, 45, 50, 49, // byte-encoding of Name

		// boot image 2
		0x82, 0, 0x11, 0x22, // boot image ID
		8,                                 // len(Name)
		98, 115, 100, 112, 45, 50, 50, 50, // byte-encoding of Name
	}
	bs, err := ParseBootImagesFromOption(input)
	require.NoError(t, err)
	require.Equal(t, len(bs), 2, "parsing 2 BootImages should return 2")
	b1 := bs[0]
	b2 := bs[1]
	expectedID1 := BootImageID{
		IsInstall: false,
		ImageType: BootImageTypeMacOSX,
		Index:     0x1010,
	}
	expectedID2 := BootImageID{
		IsInstall: true,
		ImageType: BootImageTypeMacOSXServer,
		Index:     0x1122,
	}
	require.Equal(t, expectedID1, b1.ID, "first BootImageID should be equal")
	require.Equal(t, expectedID2, b2.ID, "second BootImageID should be equal")
	require.Equal(t, "bsdp-21", b1.Name, "first BootImage name should be equal")
	require.Equal(t, "bsdp-222", b2.Name, "second BootImage name should be equal")
}

func TestParseBootImageFail(t *testing.T) {
	_, err := ParseBootImagesFromOption([]byte{})
	require.Error(t, err, "parseBootImages with empty arg")

	_, err = ParseBootImagesFromOption([]byte{1, 2, 3})
	require.Error(t, err, "parseBootImages with short arg")

	_, err = ParseBootImagesFromOption([]byte{
		// boot image 1
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                         // len(Name)
		98, 115, 100, 112, 45, 50, // byte-encoding of Name (intentionally shorter)

		// boot image 2
		0x82, 0, 0x11, 0x22, // boot image ID
		8,                                 // len(Name)
		98, 115, 100, 112, 45, 50, 50, 50, // byte-encoding of Name
	})
	require.Error(t, err, "parseBootImages with short arg")
}

/*
 * ParseVendorOptionsFromOptions
 */
func TestParseVendorOptions(t *testing.T) {
	expectedOpts := []dhcpv4.Option{
		&dhcpv4.OptionGeneric{
			OptionCode: OptionMessageType,
			Data:       []byte{byte(MessageTypeList)},
		},
		&dhcpv4.OptionGeneric{
			OptionCode: OptionVersion,
			Data:       Version1_0,
		},
	}
	var expectedOptsBytes []byte
	for _, opt := range expectedOpts {
		expectedOptsBytes = append(expectedOptsBytes, opt.ToBytes()...)
	}
	recvOpts := []dhcpv4.Option{
		&dhcpv4.OptionGeneric{
			OptionCode: dhcpv4.OptionDHCPMessageType,
			Data:       []byte{byte(dhcpv4.MessageTypeAck)},
		},
		&dhcpv4.OptionGeneric{
			OptionCode: dhcpv4.OptionBroadcastAddress,
			Data:       []byte{0xff, 0xff, 0xff, 0xff},
		},
		&dhcpv4.OptionGeneric{
			OptionCode: dhcpv4.OptionVendorSpecificInformation,
			Data:       expectedOptsBytes,
		},
	}
	opts := ParseVendorOptionsFromOptions(recvOpts)
	require.Equal(t, expectedOpts, opts, "Parsed vendorOpts should be the same")
}

func TestParseVendorOptionsFromOptionsNotPresent(t *testing.T) {
	expectedOpts := []dhcpv4.Option{
		dhcpv4.OptionGeneric{
			OptionCode: dhcpv4.OptionDHCPMessageType,
			Data:       []byte{byte(dhcpv4.MessageTypeAck)},
		},
		dhcpv4.OptionGeneric{
			OptionCode: dhcpv4.OptionBroadcastAddress,
			Data:       []byte{0xff, 0xff, 0xff, 0xff},
		},
	}
	opts := ParseVendorOptionsFromOptions(expectedOpts)
	require.Empty(t, opts, "empty vendor opts if not present in DHCP opts")
}

func TestParseVendorOptionsFromOptionsEmpty(t *testing.T) {
	opts := ParseVendorOptionsFromOptions([]dhcpv4.Option{})
	require.Empty(t, opts, "vendor opts should be empty if given an empty input")
}

func TestParseVendorOptionsFromOptionsFail(t *testing.T) {
	opts := []dhcpv4.Option{
		&dhcpv4.OptionGeneric{
			OptionCode: dhcpv4.OptionVendorSpecificInformation,
			Data: []byte{
				0x1, 0x1, 0x1, // Option 1: LIST
				0x2, 0x2, 0x01, // Option 2: Version (intentionally left short)
			},
		},
	}
	vendorOpts := ParseVendorOptionsFromOptions(opts)
	require.Empty(t, vendorOpts, "vendor opts should be empty on parse error")
}

/*
 * ParseBootImageListFromAck
 */
func TestParseBootImageListFromAck(t *testing.T) {
	expectedBootImages := []BootImage{
		BootImage{
			ID: BootImageID{
				IsInstall: true,
				ImageType: BootImageTypeMacOSX,
				Index:     0x1010,
			},
			Name: "bsdp-1",
		},
		BootImage{
			ID: BootImageID{
				IsInstall: false,
				ImageType: BootImageTypeMacOS9,
				Index:     0x1111,
			},
			Name: "bsdp-2",
		},
	}
	var bootImageBytes []byte
	for _, image := range expectedBootImages {
		bootImageBytes = append(bootImageBytes, image.ToBytes()...)
	}
	ack, _ := dhcpv4.New()
	bootImageListOpt := dhcpv4.OptionGeneric{
		OptionCode: OptionBootImageList,
		Data:       bootImageBytes,
	}
	ack.AddOption(&dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionVendorSpecificInformation,
		Data:       bootImageListOpt.ToBytes(),
	})

	images, err := ParseBootImageListFromAck(*ack)
	require.NoError(t, err)
	require.Equal(t, expectedBootImages, images, "should get same BootImages")
}

func TestParseBootImageListFromAckNoVendorOption(t *testing.T) {
	ack, _ := dhcpv4.New()
	ack.AddOption(dhcpv4.OptionGeneric{
		OptionCode: OptionMessageType,
		Data:       []byte{byte(dhcpv4.MessageTypeAck)},
	})
	images, err := ParseBootImageListFromAck(*ack)
	require.NoError(t, err, "no vendor extensions should not return error")
	require.Empty(t, images, "should not get images from ACK without Vendor extensions")
}

func TestParseBootImageListFromAckFail(t *testing.T) {
	ack, _ := dhcpv4.New()
	ack.AddOption(dhcpv4.OptionGeneric{
		OptionCode: OptionMessageType,
		Data:       []byte{byte(dhcpv4.MessageTypeAck)},
	})
	ack.AddOption(&dhcpv4.OptionGeneric{
		OptionCode: dhcpv4.OptionVendorSpecificInformation,
		Data: []byte{
			9,  // OptionBootImageList
			24, // length
			// boot image 1
			0x1, 0, 0x10, 0x10, // boot image ID
			7,                         // len(Name)
			98, 115, 100, 112, 45, 49, // byte-encoding of Name (intentionally short)

			// boot image 2
			0x82, 0, 0x11, 0x22, // boot image ID
			8,                                 // len(Name)
			98, 115, 100, 112, 45, 50, 50, 50, // byte-encoding of Name
		},
	},
	)

	images, err := ParseBootImageListFromAck(*ack)
	require.Nil(t, images, "should get nil on parse error")
	require.Error(t, err, "should get error on parse error")
}

/*
 * Private funcs
 */
func TestNeedsReplyPort(t *testing.T) {
	require.True(t, needsReplyPort(123))
	require.False(t, needsReplyPort(0))
	require.False(t, needsReplyPort(dhcpv4.ClientPort))
}
