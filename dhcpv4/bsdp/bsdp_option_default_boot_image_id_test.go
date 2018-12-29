package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/u-root/u-root/pkg/uio"
)

func TestOptDefaultBootImageIDInterfaceMethods(t *testing.T) {
	b := BootImageID{IsInstall: true, ImageType: BootImageTypeMacOSX, Index: 1001}
	o := OptDefaultBootImageID{b}
	require.Equal(t, OptionDefaultBootImageID, o.Code(), "Code")
	require.Equal(t, 4, o.Length(), "Length")
	require.Equal(t, uio.ToBigEndian(b), o.ToBytes(), "ToBytes")
}

func TestParseOptDefaultBootImageID(t *testing.T) {
	b := BootImageID{IsInstall: true, ImageType: BootImageTypeMacOSX, Index: 1001}
	o, err := ParseOptDefaultBootImageID(uio.ToBigEndian(b))
	require.NoError(t, err)
	require.Equal(t, &OptDefaultBootImageID{b}, o)

	// Short byte stream
	data := []byte{}
	_, err = ParseOptDefaultBootImageID(data)
	require.Error(t, err, "should get error from short byte stream")

	// Bad length
	data = []byte{1, 0, 0, 0, 0}
	_, err = ParseOptDefaultBootImageID(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptDefaultBootImageIDString(t *testing.T) {
	b := BootImageID{IsInstall: true, ImageType: BootImageTypeMacOSX, Index: 1001}
	o := OptDefaultBootImageID{b}
	require.Equal(t, "BSDP Default Boot Image ID -> [1001] installable macOS image", o.String())

	b = BootImageID{IsInstall: false, ImageType: BootImageTypeMacOS9, Index: 1001}
	o = OptDefaultBootImageID{b}
	require.Equal(t, "BSDP Default Boot Image ID -> [1001] uninstallable macOS 9 image", o.String())

	b = BootImageID{IsInstall: false, ImageType: BootImageType(99), Index: 1001}
	o = OptDefaultBootImageID{b}
	require.Equal(t, "BSDP Default Boot Image ID -> [1001] uninstallable unknown image", o.String())
}
