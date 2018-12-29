package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// Specific versions.
var (
	Version1_0 = []byte{1, 0}
	Version1_1 = []byte{1, 1}
)

// OptVersion represents a BSDP protocol version.
//
// Implements the BSDP option version. Can be one of 1.0 or 1.1
type OptVersion struct {
	Version []byte
}

// ParseOptVersion constructs an OptVersion struct from a sequence of
// bytes and returns it, or an error.
func ParseOptVersion(data []byte) (*OptVersion, error) {
	buf := uio.NewBigEndianBuffer(data)
	return &OptVersion{buf.CopyN(2)}, buf.FinError()
}

// Code returns the option code.
func (o *OptVersion) Code() dhcpv4.OptionCode {
	return OptionVersion
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptVersion) ToBytes() []byte {
	return append([]byte{byte(o.Code()), 2}, o.Version...)
}

// String returns a human-readable string for this option.
func (o *OptVersion) String() string {
	return fmt.Sprintf("BSDP Version -> %v.%v", o.Version[0], o.Version[1])
}

// Length returns the length of the data portion of this option.
func (o *OptVersion) Length() int {
	return 2
}
