package bsdp

import (
	"encoding/binary"
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Implements the BSDP option reply port. This is used when BSDP responses
// should be sent to a reply port other than the DHCP default. The macOS GUI
// "Startup Disk Select" sends this option since it's operating in an
// unprivileged context.

// OptReplyPort represents a BSDP protocol version.
type OptReplyPort struct {
	Port uint16
}

// ParseOptReplyPort constructs an OptReplyPort struct from a sequence of
// bytes and returns it, or an error.
func ParseOptReplyPort(data []byte) (*OptReplyPort, error) {
	if len(data) < 4 {
		return nil, dhcpv4.ErrShortByteStream
	}
	code := dhcpv4.OptionCode(data[0])
	if code != OptionReplyPort {
		return nil, fmt.Errorf("expected option %v, got %v instead", OptionReplyPort, code)
	}
	length := int(data[1])
	if length != 2 {
		return nil, fmt.Errorf("expected length 2, got %d instead", length)
	}
	port := binary.BigEndian.Uint16(data[2:4])
	return &OptReplyPort{port}, nil
}

// Code returns the option code.
func (o *OptReplyPort) Code() dhcpv4.OptionCode {
	return OptionReplyPort
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptReplyPort) ToBytes() []byte {
	serialized := make([]byte, 2)
	binary.BigEndian.PutUint16(serialized, o.Port)
	return append([]byte{byte(o.Code()), 2}, serialized...)
}

// String returns a human-readable string for this option.
func (o *OptReplyPort) String() string {
	return fmt.Sprintf("BSDP Reply Port -> %v", o.Port)
}

// Length returns the length of the data portion of this option.
func (o *OptReplyPort) Length() int {
	return 2
}
