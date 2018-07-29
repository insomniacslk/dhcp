package dhcpv4

import (
	"fmt"
	"net"
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
	if len(data) < 2 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionSubnetMask {
		return nil, fmt.Errorf("expected code %v, got %v", OptionSubnetMask, code)
	}
	length := int(data[1])
	if length != 4 {
		return nil, fmt.Errorf("unexepcted length: expected 4, got %v", length)
	}
	if len(data) < 6 {
		return nil, ErrShortByteStream
	}
	return &OptSubnetMask{SubnetMask: net.IPMask(data[2 : 2+length])}, nil
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
	return len(o.SubnetMask)
}
