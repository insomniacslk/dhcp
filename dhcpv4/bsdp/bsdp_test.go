package bsdp

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/assert"
)

/*
 * BootImageID
 */
func TestBootImageIDToBytes(t *testing.T) {
	b := BootImageID{
		IsInstall: true,
		ImageType: BootImageTypeMacOSX,
		Index:     0x1000,
	}
	actual := b.ToBytes()
	expected := []byte{0x81, 0, 0x10, 0}
	assert.Equal(t, actual, expected, "serialized BootImageID should be equal")

	b.IsInstall = false
	actual = b.ToBytes()
	expected = []byte{0x01, 0, 0x10, 0}
	assert.Equal(t, actual, expected, "serialized BootImageID should be equal")
}

func TestBootImageIDFromBytes(t *testing.T) {
	b := BootImageID{
		IsInstall: false,
		ImageType: BootImageTypeMacOSX,
		Index:     0x1000,
	}
	newBootImage, err := BootImageIDFromBytes(b.ToBytes())
	assert.Nil(t, err, "error from BootImageIDFromBytes")
	assert.Equal(t, b, *newBootImage, "deserialized BootImage should be equal")

	b = BootImageID{
		IsInstall: true,
		ImageType: BootImageTypeMacOSX,
		Index:     0x1011,
	}
	newBootImage, err = BootImageIDFromBytes(b.ToBytes())
	assert.Nil(t, err, "error from BootImageIDFromBytes")
	assert.Equal(t, b, *newBootImage, "deserialized BootImage should be equal")
}

func TestBootImageIDFromBytesFail(t *testing.T) {
	serialized := []byte{0x81, 0, 0x10} // intentionally left short
	deserialized, err := BootImageIDFromBytes(serialized)
	assert.Nil(t, deserialized, "BootImageIDFromBytes should return nil on failed deserialization")
	assert.NotNil(t, err, "BootImageIDFromBytes should return err on failed deserialization")
}

/*
 * BootImage
 */
func TestBootImageToBytes(t *testing.T) {
	b := BootImage{
		ID: BootImageID{
			IsInstall: true,
			ImageType: BootImageTypeMacOSX,
			Index:     0x1000,
		},
		Name: "bsdp-1",
	}
	expected := []byte{
		0x81, 0, 0x10, 0, // boot image ID
		6,                         // len(Name)
		98, 115, 100, 112, 45, 49, // byte-encoding of Name
	}
	actual := b.ToBytes()
	assert.Equal(t, actual, expected, "serialized BootImage should be equal")

	b = BootImage{
		ID: BootImageID{
			IsInstall: false,
			ImageType: BootImageTypeMacOSX,
			Index:     0x1010,
		},
		Name: "bsdp-21",
	}
	expected = []byte{
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                             // len(Name)
		98, 115, 100, 112, 45, 50, 49, // byte-encoding of Name
	}
	actual = b.ToBytes()
	assert.Equal(t, actual, expected, "serialized BootImage should be equal")
}

func TestBootImageFromBytes(t *testing.T) {
	input := []byte{
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                             // len(Name)
		98, 115, 100, 112, 45, 50, 49, // byte-encoding of Name
	}
	b, err := BootImageFromBytes(input)
	assert.Nil(t, err, "error while marshalling BootImage")
	expectedBootImage := BootImage{
		ID: BootImageID{
			IsInstall: false,
			ImageType: BootImageTypeMacOSX,
			Index:     0x1010,
		},
		Name: "bsdp-21",
	}
	assert.Equal(t, *b, expectedBootImage, "invalid marshalling of BootImage")
}

func TestBootImageFromBytesOnlyBootImageID(t *testing.T) {
	// Only a BootImageID, nothing else.
	input := []byte{0x1, 0, 0x10, 0x10}
	b, err := BootImageFromBytes(input)
	assert.Nil(t, b, "short bytestream should return nil BootImageID")
	assert.NotNil(t, err, "short bytestream should return error")
}

func TestBootImageFromBytesShortBootImage(t *testing.T) {
	input := []byte{
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                         // len(Name)
		98, 115, 100, 112, 45, 50, // Name bytes (intentionally off-by-one)
	}
	b, err := BootImageFromBytes(input)
	assert.Nil(t, b, "short bytestream should return nil BootImageID")
	assert.NotNil(t, err, "short bytestream should return error")
}

func TestParseBootImageSingleBootImage(t *testing.T) {
	input := []byte{
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                             // len(Name)
		98, 115, 100, 112, 45, 50, 49, // byte-encoding of Name
	}
	bs, err := ParseBootImagesFromOption(input)
	assert.Nil(t, err, "parsing single boot image should not return error")
	assert.Equal(t, len(bs), 1, "parsing single boot image should return 1")
	b := bs[0]
	expectedBootImage := BootImageID{
		IsInstall: false,
		ImageType: BootImageTypeMacOSX,
		Index:     0x1010,
	}
	assert.Equal(t, b.ID, expectedBootImage, "parsed BootImageIDs should be equal")
	assert.Equal(t, b.Name, "bsdp-21", "BootImage name should be equal")
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
	assert.Nil(t, err, "parsing multiple BootImages should not return error")
	assert.Equal(t, len(bs), 2, "parsing 2 BootImages should return 2")
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
	assert.Equal(t, b1.ID, expectedID1, "first BootImageID should be equal")
	assert.Equal(t, b2.ID, expectedID2, "second BootImageID should be equal")
	assert.Equal(t, b1.Name, "bsdp-21", "first BootImage name should be equal")
	assert.Equal(t, b2.Name, "bsdp-222", "second BootImage name should be equal")
}

func TestParseBootImageFail(t *testing.T) {
	_, err := ParseBootImagesFromOption([]byte{})
	assert.NotNil(t, err, "parseBootImages with empty arg")

	_, err = ParseBootImagesFromOption([]byte{1, 2, 3})
	assert.NotNil(t, err, "parseBootImages with short arg")

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
	assert.NotNil(t, err, "parseBootImages with short arg")
}

/*
 * ParseVendorOptionsFromOptions
 */
func TestParseVendorOptions(t *testing.T) {
	vendorOpts := []dhcpv4.Option{
		dhcpv4.Option{
			Code: OptionMessageType,
			Data: []byte{byte(MessageTypeList)},
		},
		dhcpv4.Option{
			Code: OptionVersion,
			Data: Version1_0,
		},
	}
	recvOpts := []dhcpv4.Option{
		dhcpv4.Option{
			Code: dhcpv4.OptionDHCPMessageType,
			Data: []byte{byte(dhcpv4.MessageTypeAck)},
		},
		dhcpv4.Option{
			Code: dhcpv4.OptionBroadcastAddress,
			Data: []byte{0xff, 0xff, 0xff, 0xff},
		},
		dhcpv4.Option{
			Code: dhcpv4.OptionVendorSpecificInformation,
			Data: dhcpv4.OptionsToBytesWithoutMagicCookie(vendorOpts),
		},
	}
	opts := ParseVendorOptionsFromOptions(recvOpts)
	assert.Equal(t, opts, vendorOpts, "Parsed vendorOpts should be the same")
}

func TestParseVendorOptionsFromOptionsNotPresent(t *testing.T) {
	recvOpts := []dhcpv4.Option{
		dhcpv4.Option{
			Code: dhcpv4.OptionDHCPMessageType,
			Data: []byte{byte(dhcpv4.MessageTypeAck)},
		},
		dhcpv4.Option{
			Code: dhcpv4.OptionBroadcastAddress,
			Data: []byte{0xff, 0xff, 0xff, 0xff},
		},
	}
	opts := ParseVendorOptionsFromOptions(recvOpts)
	assert.Empty(t, opts, "vendor opts should be empty if not present in input")
}

func TestParseVendorOptionsFromOptionsEmpty(t *testing.T) {
	options := ParseVendorOptionsFromOptions([]dhcpv4.Option{})
	assert.Empty(t, options, "vendor opts should be empty if given an empty input")
}

func TestParseVendorOptionsFromOptionsFail(t *testing.T) {
	opts := []dhcpv4.Option{
		dhcpv4.Option{
			Code: dhcpv4.OptionVendorSpecificInformation,
			Data: []byte{
				0x1, 0x1, 0x1, // Option 1: LIST
				0x2, 0x2, 0x01, // Option 2: Version (intentionally left short)
			},
		},
	}
	vendorOpts := ParseVendorOptionsFromOptions(opts)
	assert.Empty(t, vendorOpts, "vendor opts should be empty on parse error")
}

/*
 * ParseBootImageListFromAck
 */
func TestParseBootImageListFromAck(t *testing.T) {
	bootImages := []BootImage{
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
	for _, image := range bootImages {
		bootImageBytes = append(bootImageBytes, image.ToBytes()...)
	}
	ack, _ := dhcpv4.New()
	ack.AddOption(dhcpv4.Option{
		Code: dhcpv4.OptionVendorSpecificInformation,
		Data: dhcpv4.OptionsToBytesWithoutMagicCookie([]dhcpv4.Option{
			dhcpv4.Option{
				Code: OptionBootImageList,
				Data: bootImageBytes,
			},
		}),
	})

	images, err := ParseBootImageListFromAck(*ack)
	assert.Nil(t, err, "error from ParseBootImageListFromAck")
	assert.NotNil(t, images, "parsed boot images from ack")
	assert.Equal(t, images, bootImages, "should get same BootImages")
}

func TestParseBootImageListFromAckNoVendorOption(t *testing.T) {
	ack, _ := dhcpv4.New()
	ack.AddOption(dhcpv4.Option{
		Code: OptionMessageType,
		Data: []byte{byte(dhcpv4.MessageTypeAck)},
	})
	images, err := ParseBootImageListFromAck(*ack)
	assert.Nil(t, err, "no vendor extensions should not return error")
	assert.Empty(t, images, "should not get images from ACK without Vendor extensions")
}

func TestParseBootImageListFromAckFail(t *testing.T) {
	ack, _ := dhcpv4.New()
	ack.AddOption(dhcpv4.Option{
		Code: OptionMessageType,
		Data: []byte{byte(dhcpv4.MessageTypeAck)},
	})
	ack.AddOption(dhcpv4.Option{
		Code: dhcpv4.OptionVendorSpecificInformation,
		Data: dhcpv4.OptionsToBytesWithoutMagicCookie([]dhcpv4.Option{
			dhcpv4.Option{
				Code: OptionBootImageList,
				Data: []byte{
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
		}),
	})

	images, err := ParseBootImageListFromAck(*ack)
	assert.Nil(t, images, "should get nil on parse error")
	assert.NotNil(t, err, "should get error on parse error")
}

/*
 * Private funcs
 */
func TestNeedsReplyPort(t *testing.T) {
	assert.True(t, needsReplyPort(123), "")
	assert.False(t, needsReplyPort(0), "")
	assert.False(t, needsReplyPort(dhcpv4.ClientPort), "")
}
