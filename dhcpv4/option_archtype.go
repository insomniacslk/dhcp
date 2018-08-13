package dhcpv4

// This option implements the Client System Architecture Type option
// https://tools.ietf.org/html/rfc4578

import (
	"encoding/binary"
	"fmt"

	"github.com/insomniacslk/dhcp/iana"
)

// OptClientArchType represents an option encapsulating the Client System
// Architecture Type option Definition.
type OptClientArchType struct {
	ArchTypes []iana.ArchType
}

// Code returns the option code.
func (o *OptClientArchType) Code() OptionCode {
	return OptionClientSystemArchitectureType
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptClientArchType) ToBytes() []byte {
	ret := []byte{byte(o.Code()), byte(o.Length())}
	for _, at := range o.ArchTypes {
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf[0:2], uint16(at))
		ret = append(ret, buf...)
	}
	return ret
}

// Length returns the length of the data portion (excluding option code an byte
// length).
func (o *OptClientArchType) Length() int {
	return 2*len(o.ArchTypes)
}

// String returns a human-readable string.
func (o *OptClientArchType) String() string {
	var archTypes string
	for idx, at := range o.ArchTypes {
		name := iana.ArchTypeToString(at)
		archTypes += name
		if idx < len(o.ArchTypes)-1 {
			archTypes += ", "
		}
	}
	return fmt.Sprintf("Client System Architecture Type -> %v", archTypes)
}

// ParseOptClientArchType returns a new OptClientArchType from a byte stream,
// or error if any.
func ParseOptClientArchType(data []byte) (*OptClientArchType, error) {
	if len(data) < 2 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionClientSystemArchitectureType {
		return nil, fmt.Errorf("expected code %v, got %v", OptionClientSystemArchitectureType, code)
	}
	length := int(data[1])
	if length == 0 || length%2 != 0 {
		return nil, fmt.Errorf("Invalid length: expected multiple of 2 larger than 2, got %v", length)
	}
	if len(data) < 2+length {
		return nil, ErrShortByteStream
	}
	archTypes := make([]iana.ArchType, 0, length%2)
	for idx := 0; idx < length; idx += 2 {
		b := data[2+idx : 2+idx+2]
		archTypes = append(archTypes, iana.ArchType(binary.BigEndian.Uint16(b)))
	}
	return &OptClientArchType{ArchTypes: archTypes}, nil
}
