package dhcpv4

import (
	"fmt"
	"net"

	"github.com/u-root/u-root/pkg/uio"
)

// This option implements the subnet mask option
// https://tools.ietf.org/html/rfc2132

// OptSubnetMask represents an option encapsulating the subnet mask.
type OptSubnetMask struct {
	SubnetMask net.IPMask
}

// ParseOptSubnetMask returns a new OptSubnetMask from a byte
// stream, or error if any.
func ParseOptSubnetMask(data []byte) (*OptSubnetMask, error) {
	buf := uio.NewBigEndianBuffer(data)
	return &OptSubnetMask{SubnetMask: net.IPMask(buf.CopyN(net.IPv4len))}, buf.FinError()
}

// Code returns the option code.
func (o *OptSubnetMask) Code() OptionCode {
	return OptionSubnetMask
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptSubnetMask) ToBytes() []byte {
	ret := []byte{byte(o.Code()), byte(o.Length())}
	return append(ret, o.SubnetMask[:4]...)
}

// String returns a human-readable string.
func (o *OptSubnetMask) String() string {
	return fmt.Sprintf("Subnet Mask -> %v", o.SubnetMask.String())
}

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptSubnetMask) Length() int {
	return 4
}
