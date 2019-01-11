package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// Version is the BSDP protocol version. Can be one of 1.0 or 1.1.
type Version [2]byte

// Specific versions.
var (
	Version1_0 = Version{1, 0}
	Version1_1 = Version{1, 1}
)

// ParseOptVersion constructs an OptVersion struct from a sequence of
// bytes and returns it, or an error.
func ParseOptVersion(data []byte) (Version, error) {
	buf := uio.NewBigEndianBuffer(data)
	var v Version
	buf.ReadBytes(v[:])
	return v, buf.FinError()
}

// Code returns the option code.
func (o Version) Code() dhcpv4.OptionCode {
	return OptionVersion
}

// ToBytes returns a serialized stream of bytes for this option.
func (o Version) ToBytes() []byte {
	return o[:]
}

// String returns a human-readable string for this option.
func (o Version) String() string {
	return fmt.Sprintf("BSDP Version -> %d.%d", o[0], o[1])
}
