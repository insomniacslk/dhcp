package bsdp

import (
	"encoding/binary"
	"fmt"
)

// BootImageType represents the different BSDP boot image types.
type BootImageType byte

// Different types of BootImages - e.g. for different flavors of macOS.
const (
	BootImageTypeMacOS9              BootImageType = 0
	BootImageTypeMacOSX              BootImageType = 1
	BootImageTypeMacOSXServer        BootImageType = 2
	BootImageTypeHardwareDiagnostics BootImageType = 3
	// 4 - 127 are reserved for future use.
)

// BootImageTypeToString maps the different BootImageTypes to human-readable
// representations.
var BootImageTypeToString = map[BootImageType]string{
	BootImageTypeMacOS9:              "macOS 9",
	BootImageTypeMacOSX:              "macOS",
	BootImageTypeMacOSXServer:        "macOS Server",
	BootImageTypeHardwareDiagnostics: "Hardware Diagnostic",
}

// BootImageID describes a boot image ID - whether it's an install image and
// what kind of boot image (e.g. OS 9, macOS, hardware diagnostics)
type BootImageID struct {
	IsInstall bool
	ImageType BootImageType
	Index     uint16
}

// ToBytes serializes a BootImageID to network-order bytes.
func (b BootImageID) ToBytes() []byte {
	bytes := make([]byte, 4)
	if b.IsInstall {
		bytes[0] |= 0x80
	}
	bytes[0] |= byte(b.ImageType)
	binary.BigEndian.PutUint16(bytes[2:], b.Index)
	return bytes
}

// String converts a BootImageID to a human-readable representation.
func (b BootImageID) String() string {
	s := fmt.Sprintf("[%d]", b.Index)
	if b.IsInstall {
		s += " installable"
	} else {
		s += " uninstallable"
	}
	t, ok := BootImageTypeToString[b.ImageType]
	if !ok {
		t = "unknown"
	}
	return s + " " + t + " image"
}

// BootImageIDFromBytes deserializes a collection of 4 bytes to a BootImageID.
func BootImageIDFromBytes(bytes []byte) (*BootImageID, error) {
	if len(bytes) < 4 {
		return nil, fmt.Errorf("not enough bytes to serialize BootImageID")
	}
	return &BootImageID{
		IsInstall: bytes[0]&0x80 != 0,
		ImageType: BootImageType(bytes[0] & 0x7f),
		Index:     binary.BigEndian.Uint16(bytes[2:]),
	}, nil
}

// BootImage describes a boot image - contains the boot image ID and the name.
type BootImage struct {
	ID   BootImageID
	Name string
}

// ToBytes converts a BootImage to a slice of bytes.
func (b *BootImage) ToBytes() []byte {
	bytes := b.ID.ToBytes()
	bytes = append(bytes, byte(len(b.Name)))
	bytes = append(bytes, []byte(b.Name)...)
	return bytes
}

// String converts a BootImage to a human-readable representation.
func (b *BootImage) String() string {
	return fmt.Sprintf("%v %v", b.Name, b.ID.String())
}

// BootImageFromBytes returns a deserialized BootImage struct from bytes.
func BootImageFromBytes(bytes []byte) (*BootImage, error) {
	// Should at least contain 4 bytes of BootImageID + byte for length of
	// boot image name.
	if len(bytes) < 5 {
		return nil, fmt.Errorf("not enough bytes to serialize BootImage")
	}
	imageID, err := BootImageIDFromBytes(bytes[:4])
	if err != nil {
		return nil, err
	}
	nameLength := int(bytes[4])
	if 5+nameLength > len(bytes) {
		return nil, fmt.Errorf("not enough bytes for BootImage")
	}
	name := string(bytes[5 : 5+nameLength])
	return &BootImage{ID: *imageID, Name: name}, nil
}
