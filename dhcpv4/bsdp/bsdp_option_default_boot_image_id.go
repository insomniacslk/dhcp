package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// OptDefaultBootImageID contains the selected boot image ID.
//
// Implements the BSDP option default boot image ID, which tells the client
// which image is the default boot image if one is not selected.
type OptDefaultBootImageID struct {
	ID BootImageID
}

// ParseOptDefaultBootImageID constructs an OptDefaultBootImageID struct from a sequence of
// bytes and returns it, or an error.
func ParseOptDefaultBootImageID(data []byte) (*OptDefaultBootImageID, error) {
	var o OptDefaultBootImageID
	buf := uio.NewBigEndianBuffer(data)
	if err := o.ID.Unmarshal(buf); err != nil {
		return nil, err
	}
	return &o, buf.FinError()
}

// Code returns the option code.
func (o *OptDefaultBootImageID) Code() dhcpv4.OptionCode {
	return OptionDefaultBootImageID
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptDefaultBootImageID) ToBytes() []byte {
	serializedID := o.ID.ToBytes()
	return append([]byte{byte(o.Code()), byte(len(serializedID))}, serializedID...)
}

// String returns a human-readable string for this option.
func (o *OptDefaultBootImageID) String() string {
	return fmt.Sprintf("BSDP Default Boot Image ID -> %s", o.ID.String())
}

// Length returns the length of the data portion of this option.
func (o *OptDefaultBootImageID) Length() int {
	return 4
}
