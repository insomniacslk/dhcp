package dhcpv4

import "fmt"

// This option implements the host name option
// https://tools.ietf.org/html/rfc2132.txt

// OptHostName represents an option encapsulating the host name.
type OptHostName struct {
	HostName string
}

// ParseOptHostName returns a new OptHostName from a byte stream, or error if
// any.
func ParseOptHostName(data []byte) (*OptHostName, error) {
	if len(data) < 3 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionHostName {
		return nil, fmt.Errorf("expected code %v, got %v", OptionHostName, code)
	}
	length := int(data[1])
	if len(data) < 2+length {
		return nil, ErrShortByteStream
	}
	return &OptHostName{HostName: string(data[2 : 2+length])}, nil
}

// Code returns the option code.
func (o *OptHostName) Code() OptionCode {
	return OptionHostName
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptHostName) ToBytes() []byte {
	return append([]byte{byte(o.Code()), byte(o.Length())}, []byte(o.HostName)...)
}

// String returns a human-readable string.
func (o *OptHostName) String() string {
	return fmt.Sprintf("Host Name -> %v", o.HostName)
}

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptHostName) Length() int {
	return len(o.HostName)
}
