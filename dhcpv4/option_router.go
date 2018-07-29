package dhcpv4

import (
	"fmt"
	"net"
)

// This option implements the router option
// https://tools.ietf.org/html/rfc2132

// OptDomainRouter represents an option encapsulating the routers.
type OptRouter struct {
	Routers []net.IP
}

// ParseOptRouter returns a new OptRouter from a byte  stream, or error if any.
func ParseOptRouter(data []byte) (*OptRouter, error) {
	if len(data) < 2 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionRouter {
		return nil, fmt.Errorf("expected code %v, got %v", OptionRouter, code)
	}
	length := int(data[1])
	if length == 0 || length%4 != 0 {
		return nil, fmt.Errorf("Invalid length: expected multiple of 4 larger than 4, got %v", length)
	}
	if len(data) < 2+length {
		return nil, ErrShortByteStream
	}
	routers := make([]net.IP, 0, length%4)
	for idx := 0; idx < length; idx += 4 {
		b := data[2+idx : 2+idx+4]
		routers = append(routers, net.IPv4(b[0], b[1], b[2], b[3]))
	}
	return &OptRouter{Routers: routers}, nil
}

// Code returns the option code.
func (o *OptRouter) Code() OptionCode {
	return OptionRouter
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptRouter) ToBytes() []byte {
	ret := []byte{byte(o.Code()), byte(o.Length())}
	for _, router := range o.Routers {
		ret = append(ret, router.To4()...)
	}
	return ret
}

// String returns a human-readable string.
func (o *OptRouter) String() string {
	var routers string
	for idx, router := range o.Routers {
		routers += router.String()
		if idx < len(o.Routers)-1 {
			routers += ", "
		}
	}
	return fmt.Sprintf("Routers -> %v", routers)
}

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptRouter) Length() int {
	return len(o.Routers) * 4
}
