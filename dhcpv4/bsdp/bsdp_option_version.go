// +build darwin

package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Implements the BSDP option version. Can be one of 1.0 or 1.1

// Specific versions.
var (
	Version1_0 = []byte{1, 0}
	Version1_1 = []byte{1, 1}
)

// OptVersion represents a BSDP protocol version.
type OptVersion struct {
	Version []byte
}

// ParseOptVersion constructs an OptVersion struct from a sequence of
// bytes and returns it, or an error.
func ParseOptVersion(data []byte) (*OptVersion, error) {
	if len(data) < 4 {
		return nil, dhcpv4.ErrShortByteStream
	}
	code := dhcpv4.OptionCode(data[0])
	if code != OptionVersion {
		return nil, fmt.Errorf("expected option %v, got %v instead", OptionVersion, code)
	}
	length := int(data[1])
	if length != 2 {
		return nil, fmt.Errorf("expected length 2, got %d instead", length)
	}
	return &OptVersion{data[2:4]}, nil
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
