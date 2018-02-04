package dhcpv4

import (
	"bytes"
	"testing"
)

/*
 * BootImageID
 */
func TestBootImageIDToBytes(t *testing.T) {
	b := BootImageID{
		isInstall: true,
		imageKind: BSDPBootImageMacOSX,
		index:     0x1000,
	}
	actual := b.toBytes()
	expected := []byte{0x81, 0, 0x10, 0}
	if !bytes.Equal(actual, expected) {
		t.Fatalf("Invalid bytes conversion: expected %v, got %v", expected, actual)
	}

	b.isInstall = false
	actual = b.toBytes()
	expected = []byte{0x01, 0, 0x10, 0}
	if !bytes.Equal(actual, expected) {
		t.Fatalf("Invalid bytes conversion: expected %v, got %v", expected, actual)
	}
}

func TestBootImageIDFromBytes(t *testing.T) {
	b := BootImageID{
		isInstall: false,
		imageKind: BSDPBootImageMacOSX,
		index:     0x1000,
	}
	newBootImage := bootImageIDFromBytes(b.toBytes())
	if b != newBootImage {
		t.Fatalf("Difference in BootImageIDs: expected %v, got %v", b, newBootImage)
	}

	b = BootImageID{
		isInstall: true,
		imageKind: BSDPBootImageMacOSX,
		index:     0x1011,
	}
	newBootImage = bootImageIDFromBytes(b.toBytes())
	if b != newBootImage {
		t.Fatalf("Difference in BootImageIDs: expected %v, got %v", b, newBootImage)
	}
}

/*
 * BootImage
 */
func TestBootImageToBytes(t *testing.T) {
	b := BootImage{
		ID: BootImageID{
			isInstall: true,
			imageKind: BSDPBootImageMacOSX,
			index:     0x1000,
		},
		Name: "bsdp-1",
	}
	expected := []byte{
		0x81, 0, 0x10, 0, // boot image ID
		6,                         // len(Name)
		98, 115, 100, 112, 45, 49, // byte-encoding of Name
	}
	actual := b.toBytes()
	if !bytes.Equal(expected, actual) {
		t.Fatalf("Invalid bytes conversion: expected %v, got %v", expected, actual)
	}

	b = BootImage{
		ID: BootImageID{
			isInstall: false,
			imageKind: BSDPBootImageMacOSX,
			index:     0x1010,
		},
		Name: "bsdp-21",
	}
	expected = []byte{
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                             // len(Name)
		98, 115, 100, 112, 45, 50, 49, // byte-encoding of Name
	}
	actual = b.toBytes()
	if !bytes.Equal(expected, actual) {
		t.Fatalf("Invalid bytes conversion: expected %v, got %v", expected, actual)
	}
}

func TestBootImageFromBytes(t *testing.T) {
	input := []byte{
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                             // len(Name)
		98, 115, 100, 112, 45, 50, 49, // byte-encoding of Name
	}
	b, read, err := bootImageFromBytes(input)
	AssertNil(t, err, "error while marshalling BootImage")
	AssertEqual(t, read, len(input), "number of bytes from input")
	expectedBootImage := BootImage{
		ID: BootImageID{
			isInstall: false,
			imageKind: BSDPBootImageMacOSX,
			index:     0x1010,
		},
		Name: "bsdp-21",
	}
	AssertEqual(t, *b, expectedBootImage, "invalid marshalling of BootImage")
}

func TestBootImageFromBytesOnlyBootImageID(t *testing.T) {
	// Only a BootImageID, nothing else.
	input := []byte{0x1, 0, 0x10, 0x10}
	_, _, err := bootImageFromBytes(input)
	AssertNotNil(t, err, "short bytestream")
}

func TestBootImageFromBytesShortBootImage(t *testing.T) {
	input := []byte{
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                         // len(Name)
		98, 115, 100, 112, 45, 50, // Name bytes (intentially off-by-one)
	}
	_, _, err := bootImageFromBytes(input)
	AssertNotNil(t, err, "short bytestream")
}

func TestParseBootImageSingleBootImage(t *testing.T) {
	input := []byte{
		0x1, 0, 0x10, 0x10, // boot image ID
		7,                             // len(Name)
		98, 115, 100, 112, 45, 50, 49, // byte-encoding of Name
	}
	bs, err := parseBootImagesFromBSDPOption(input)
	AssertNil(t, err, "parsing boot image")
	AssertEqual(t, len(bs), 1, "length of boot images")
	b := bs[0]
	expectedBootImage := BootImageID{
		isInstall: false,
		imageKind: BSDPBootImageMacOSX,
		index:     0x1010,
	}
	AssertEqual(t, b.ID, expectedBootImage, "boot image ID")
	AssertEqual(t, b.Name, "bsdp-21", "boot image name")
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
	bs, err := parseBootImagesFromBSDPOption(input)
	AssertNil(t, err, "parsing boot image")
	AssertEqual(t, len(bs), 2, "length of boot images")
	b1 := bs[0]
	b2 := bs[1]
	expectedID1 := BootImageID{
		isInstall: false,
		imageKind: BSDPBootImageMacOSX,
		index:     0x1010,
	}
	expectedID2 := BootImageID{
		isInstall: true,
		imageKind: BSDPBootImageMacOSXServer,
		index:     0x1122,
	}
	AssertEqual(t, b1.ID, expectedID1, "boot image ID 1")
	AssertEqual(t, b2.ID, expectedID2, "boot image ID 2")
	AssertEqual(t, b1.Name, "bsdp-21", "boot image 1 name")
	AssertEqual(t, b2.Name, "bsdp-222", "boot image 1 name")
}

func TestParseBootImageFail(t *testing.T) {
	_, err := parseBootImagesFromBSDPOption([]byte{})
	AssertNotNil(t, err, "parseBootImages with empty arg")

	_, err = parseBootImagesFromBSDPOption([]byte{1, 2, 3})
	AssertNotNil(t, err, "parseBootImages with short arg")
}

/*
 * parseVendorOptionsFromOptions
 */
func TestParseVendorOptions(t *testing.T) {
	recvOpts := []Option{
		Option{
			Code: OptionDHCPMessageType,
			Data: []byte{MessageTypeAck},
		},
		Option{
			Code: OptionBroadcastAddress,
			Data: []byte{0xff, 0xff, 0xff, 0xff},
		},
		Option{
			Code: OptionVendorSpecificInformation,
			Data: OptionsToBytes([]Option{
				Option{
					Code: BSDPOptionMessageType,
					Data: []byte{BSDPMessageTypeList},
				},
				Option{
					Code: BSDPOptionVersion,
					Data: BSDPVersion1_0,
				},
			}),
		},
	}
	opts := parseVendorOptionsFromOptions(recvOpts)
	AssertEqual(t, len(opts), 2, "len of vendor opts")
}

func TestParseVendorOptionsFromOptionsNotPresent(t *testing.T) {
	recvOpts := []Option{
		Option{
			Code: OptionDHCPMessageType,
			Data: []byte{MessageTypeAck},
		},
		Option{
			Code: OptionBroadcastAddress,
			Data: []byte{0xff, 0xff, 0xff, 0xff},
		},
	}
	opts := parseVendorOptionsFromOptions(recvOpts)
	AssertEqual(t, len(opts), 0, "len of vendor opts")
}

func TestParseVendorOptionsFromOptionsEmpty(t *testing.T) {
	options := parseVendorOptionsFromOptions([]Option{})
	AssertEqual(t, len(options), 0, "size of options")
}

/*
 * ParseBootImageListFromAck
 */
func TestParseBootImageListFromAck(t *testing.T) {
	bootImages := []BootImage{
		BootImage{
			ID: BootImageID{
				isInstall: true,
				imageKind: BSDPBootImageMacOSX,
				index:     0x1010,
			},
			Name: "bsdp-1",
		},
		BootImage{
			ID: BootImageID{
				isInstall: false,
				imageKind: BSDPBootImageMacOS9,
				index:     0x1111,
			},
			Name: "bsdp-2",
		},
	}
	var bootImageBytes []byte
	for _, image := range bootImages {
		bootImageBytes = append(bootImageBytes, image.toBytes()...)
	}
	ack := DHCPv4{
		options: []Option{
			Option{
				Code: OptionVendorSpecificInformation,
				Data: OptionsToBytes([]Option{
					Option{
						Code: BSDPOptionBootImageList,
						Data: bootImageBytes,
					},
				}),
			},
		},
	}

	images, err := ParseBootImageListFromAck(ack)
	AssertNil(t, err, "error from ParseBootImageListFromAck")
	AssertNotNil(t, images, "parsed boot images from ack")
	if len(images) != len(bootImages) {
		t.Fatalf("Expected same number of BootImages, got %d instead", len(images))
	}
	for i := range images {
		if images[i] != bootImages[i] {
			t.Fatalf("Expected boot images to be same. %v != %v", images[i], bootImages[i])
		}
	}
}

/*
 * NewInformListForInterface
 */
