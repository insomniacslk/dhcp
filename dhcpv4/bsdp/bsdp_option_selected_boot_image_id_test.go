package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/u-root/u-root/pkg/uio"
)

func TestOptSelectedBootImageIDInterfaceMethods(t *testing.T) {
	b := BootImageID{IsInstall: true, ImageType: BootImageTypeMacOSX, Index: 1001}
	o := OptSelectedBootImageID{b}
	require.Equal(t, OptionSelectedBootImageID, o.Code(), "Code")
	require.Equal(t, 4, o.Length(), "Length")
	require.Equal(t, uio.ToBigEndian(b), o.ToBytes(), "ToBytes")
}

func TestParseOptSelectedBootImageID(t *testing.T) {
	b := BootImageID{IsInstall: true, ImageType: BootImageTypeMacOSX, Index: 1001}
	o, err := ParseOptSelectedBootImageID(uio.ToBigEndian(b))
	require.NoError(t, err)
	require.Equal(t, &OptSelectedBootImageID{b}, o)

	// Short byte stream
	data := []byte{}
	_, err = ParseOptSelectedBootImageID(data)
	require.Error(t, err, "should get error from short byte stream")

	// Bad length
	data = []byte{1, 0, 0, 0, 0}
	_, err = ParseOptSelectedBootImageID(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptSelectedBootImageIDString(t *testing.T) {
	b := BootImageID{IsInstall: true, ImageType: BootImageTypeMacOSX, Index: 1001}
	o := OptSelectedBootImageID{b}
	require.Equal(t, "BSDP Selected Boot Image ID -> [1001] installable macOS image", o.String())

	b = BootImageID{IsInstall: false, ImageType: BootImageTypeMacOS9, Index: 1001}
	o = OptSelectedBootImageID{b}
	require.Equal(t, "BSDP Selected Boot Image ID -> [1001] uninstallable macOS 9 image", o.String())

	b = BootImageID{IsInstall: false, ImageType: BootImageType(99), Index: 1001}
	o = OptSelectedBootImageID{b}
	require.Equal(t, "BSDP Selected Boot Image ID -> [1001] uninstallable unknown image", o.String())
}
