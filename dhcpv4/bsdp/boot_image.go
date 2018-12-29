package bsdp

import (
	"encoding/binary"
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
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

// Unmarshal reads b's binary representation from buf.
func (b *BootImageID) Unmarshal(buf *uio.Lexer) error {
	byte0 := buf.Read8()
	_ = buf.Read8()
	b.IsInstall = byte0&0x80 != 0
	b.ImageType = BootImageType(byte0 & 0x7f)
	b.Index = buf.Read16()
	return buf.Error()
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
func (b BootImage) String() string {
	return fmt.Sprintf("%v %v", b.Name, b.ID.String())
}

// Unmarshal reads data from buf into b.
func (b *BootImage) Unmarshal(buf *uio.Lexer) error {
	if err := (&b.ID).Unmarshal(buf); err != nil {
		return err
	}
	nameLength := buf.Read8()
	b.Name = string(buf.Consume(int(nameLength)))
	return buf.Error()
}
