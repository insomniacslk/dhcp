package dhcpv4

import (
	"fmt"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/u-root/u-root/pkg/uio"
)

// OptClientArchType represents an option encapsulating the Client System
// Architecture Type option definition. See RFC 4578.
type OptClientArchType struct {
	ArchTypes []iana.Arch
}

// Code returns the option code.
func (o *OptClientArchType) Code() OptionCode {
	return OptionClientSystemArchitectureType
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptClientArchType) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, at := range o.ArchTypes {
		buf.Write16(uint16(at))
	}
	return buf.Data()
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

	archTypes := make([]iana.Arch, 0, buf.Len()/2)
	for buf.Has(2) {
		archTypes = append(archTypes, iana.Arch(buf.Read16()))
	}
	return &OptClientArchType{ArchTypes: archTypes}, buf.FinError()
}
