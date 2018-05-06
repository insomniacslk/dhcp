package dhcpv4

import (
	"fmt"
	"net"
)

/*
This is a helper option that provides facilities to create options that
contain an IPv4 address. It cannot be used directly, but the deriving
implementations have to embed it, and implement `Code() OptionCode`,
`ToBytes() []byte`, `String() string` and
`Parse(data []byte) (*YourNewOption, error)`. The latter is a function,
the others are methods of the new option.

Example of a new OptMyIPOption based on OptGenericIP:

type OptMyIPOption struct {
	OptGenericIP
}

func ParseOptMyIPOption(data []byte) (*OptMyIPOption, error) {
	opt, err := ParseOptGenericIP(OptionMyIPOption, data)
	if err != nil {
		return nil, err
	}
	return &Opt
	return &OptMyIPOption{OptGenericIP: *opt}, nil
}

func (o OptMyIPOption) String() string {
	return fmt.Sprintf("OptMyIPOption -> %v", o.Value)
}

func (o OptMyIPOption) Code() OptionCode {
	return OptionMyIPOption
}

func (o OptMyIPOption) ToBytes() []byte {
	return o.OptGenericIP.ToBytes(OptionMyIPOption)
}
*/

// OptGenericIP represents an option encapsulating the server identifier.
type OptGenericIP struct {
	IP net.IP
}

// ParseOptGenericIP returns a new OptGenericIP from a byte stream, or error if
// any.
func ParseOptGenericIP(optCode OptionCode, data []byte) (*OptGenericIP, error) {
	if len(data) < 2 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != optCode {
		return nil, fmt.Errorf("expected code %v, got %v", optCode, code)
	}
	length := int(data[1])
	if length != 4 {
		return nil, fmt.Errorf("unexepcted length: expected 4, got %v", length)
	}
	if len(data) < 6 {
		return nil, ErrShortByteStream
	}
	return &OptGenericIP{IP: net.IP(data[2 : 2+length])}, nil
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptGenericIP) ToBytes(optCode OptionCode) []byte {
	ret := []byte{byte(optCode), byte(o.Length())}
	return append(ret, o.IP.To4()...)
}

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptGenericIP) Length() int {
	return len(o.IP.To4())
}
