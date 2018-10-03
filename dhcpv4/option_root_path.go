package dhcpv4

import (
	"fmt"
)

// This option implements the root path option
// https://tools.ietf.org/html/rfc2132

// OptRootPath represents the path to the client's root disk.
type OptRootPath struct {
	Path string
}

// ParseOptRootPath constructs an OptRootPath struct from a sequence of  bytes
// and returns it, or an error.
func ParseOptRootPath(data []byte) (*OptRootPath, error) {
	// Should at least have code and length
	if len(data) < 2 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionRootPath {
		return nil, fmt.Errorf("expected option %v, got %v instead", OptionRootPath, code)
	}
	length := int(data[1])
	if len(data) < 2+length {
		return nil, ErrShortByteStream
	}
	return &OptRootPath{Path: string(data[2 : 2+length])}, nil
}

// Code returns the option code.
func (o *OptRootPath) Code() OptionCode {
	return OptionRootPath
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptRootPath) ToBytes() []byte {
	return append([]byte{byte(o.Code()), byte(o.Length())}, []byte(o.Path)...)
}

// String returns a human-readable string for this option.
func (o *OptRootPath) String() string {
	return fmt.Sprintf("Root Path -> %v", o.Path)
}

// Length returns the length of the data portion (excluding option code and byte
// for length, if any).
func (o *OptRootPath) Length() int {
	return len(o.Path)
}
