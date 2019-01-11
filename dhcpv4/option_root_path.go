package dhcpv4

import (
	"fmt"
)

// OptRootPath implements the root path option described by RFC 2132, Section
// 3.19.
type OptRootPath struct {
	Path string
}

// ParseOptRootPath constructs an OptRootPath struct from a sequence of  bytes
// and returns it, or an error.
func ParseOptRootPath(data []byte) (*OptRootPath, error) {
	return &OptRootPath{Path: string(data)}, nil
}

// Code returns the option code.
func (o *OptRootPath) Code() OptionCode {
	return OptionRootPath
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptRootPath) ToBytes() []byte {
	return []byte(o.Path)
}

// String returns a human-readable string for this option.
func (o *OptRootPath) String() string {
	return fmt.Sprintf("Root Path -> %v", o.Path)
}
