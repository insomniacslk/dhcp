package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// OptSelectedBootImageID contains the selected boot image ID.
//
// Implements the BSDP option selected boot image ID, which tells the server
// which boot image has been selected by the client.
type OptSelectedBootImageID struct {
	ID BootImageID
}

// ParseOptSelectedBootImageID constructs an OptSelectedBootImageID struct from a sequence of
// bytes and returns it, or an error.
func ParseOptSelectedBootImageID(data []byte) (*OptSelectedBootImageID, error) {
	var o OptSelectedBootImageID
	buf := uio.NewBigEndianBuffer(data)
	if err := o.ID.Unmarshal(buf); err != nil {
		return nil, err
	}
	return &o, buf.FinError()
}

// Code returns the option code.
func (o *OptSelectedBootImageID) Code() dhcpv4.OptionCode {
	return OptionSelectedBootImageID
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptSelectedBootImageID) ToBytes() []byte {
	serializedID := o.ID.ToBytes()
	return append([]byte{byte(o.Code()), byte(len(serializedID))}, serializedID...)
}

// String returns a human-readable string for this option.
func (o *OptSelectedBootImageID) String() string {
	return fmt.Sprintf("BSDP Selected Boot Image ID -> %s", o.ID.String())
}

// Length returns the length of the data portion of this option.
func (o *OptSelectedBootImageID) Length() int {
	return 4
}
