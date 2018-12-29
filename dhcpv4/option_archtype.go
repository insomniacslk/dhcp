package dhcpv4

// This option implements the Client System Architecture Type option
// https://tools.ietf.org/html/rfc4578

import (
	"encoding/binary"
	"fmt"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/u-root/u-root/pkg/uio"
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
	return 2 * len(o.ArchTypes)
}

// String returns a human-readable string.
func (o *OptClientArchType) String() string {
	var archTypes string
	for idx, at := range o.ArchTypes {
		archTypes += at.String()
		if idx < len(o.ArchTypes)-1 {
			archTypes += ", "
		}
	}
	return fmt.Sprintf("Client System Architecture Type -> %v", archTypes)
}

// ParseOptClientArchType returns a new OptClientArchType from a byte stream,
// or error if any.
func ParseOptClientArchType(data []byte) (*OptClientArchType, error) {
	buf := uio.NewBigEndianBuffer(data)
	if buf.Len() == 0 {
		return nil, fmt.Errorf("must have at least one archtype if option is present")
	}

	archTypes := make([]iana.ArchType, 0, buf.Len()/2)
	for buf.Has(2) {
		archTypes = append(archTypes, iana.ArchType(buf.Read16()))
	}
	return &OptClientArchType{ArchTypes: archTypes}, buf.FinError()
}
