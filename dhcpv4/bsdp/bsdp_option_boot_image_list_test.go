// +build darwin

package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptBootImageListInterfaceMethods(t *testing.T) {
	bs := []BootImage{
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
	}
	o := OptBootImageList{bs}
	require.Equal(t, OptionBootImageList, o.Code(), "Code")
	require.Equal(t, 22, o.Length(), "Length")
	expectedBytes := []byte{
		9,  // code
		22, // length
		// boot image 1
		0x1, 0x0, 0x03, 0xe9, // ID
		6, // name length
		'b', 's', 'd', 'p', '-', '1',
		// boot image 1
		0x80, 0x0, 0x23, 0x31, // ID
		6, // name length
		'b', 's', 'd', 'p', '-', '2',
	}
	require.Equal(t, expectedBytes, o.ToBytes(), "ToBytes")
}

func TestParseOptBootImageList(t *testing.T) {
	data := []byte{
		9,  // code
		22, // length
		// boot image 1
		0x1, 0x0, 0x03, 0xe9, // ID
		6, // name length
		'b', 's', 'd', 'p', '-', '1',
		// boot image 1
		0x80, 0x0, 0x23, 0x31, // ID
		6, // name length
		'b', 's', 'd', 'p', '-', '2',
	}
	o, err := ParseOptBootImageList(data)
	require.NoError(t, err)
	expectedBootImages := []BootImage{
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
	}
	require.Equal(t, &OptBootImageList{expectedBootImages}, o)

	// Short byte stream
	data = []byte{9}
	_, err = ParseOptBootImageList(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{54, 1, 1}
	_, err = ParseOptBootImageList(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{9, 10, 1, 1, 1}
	_, err = ParseOptBootImageList(data)
	require.Error(t, err, "should get error from bad length")

	// Error parsing boot image (malformed)
	data = []byte{
		9,  // code
		22, // length
		// boot image 1
		0x1, 0x0, 0x03, 0xe9, // ID
		4, // name length
		'b', 's', 'd', 'p', '-', '1',
		// boot image 1
		0x80, 0x0, 0x23, 0x31, // ID
		6, // name length
		'b', 's', 'd', 'p', '-', '2',
	}
	_, err = ParseOptBootImageList(data)
	require.Error(t, err, "should get error from bad boot image")
}

func TestOptBootImageListString(t *testing.T) {
	bs := []BootImage{
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
	}
	o := OptBootImageList{bs}
	expectedString := "BSDP Boot Image List ->\n  bsdp-1 [1001] uninstallable macOS image\n  bsdp-2 [9009] installable macOS 9 image"
	require.Equal(t, expectedString, o.String())
}
