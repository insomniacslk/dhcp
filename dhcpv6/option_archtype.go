package dhcpv6

// This module defines the OptClientArchType structure.
// https://www.ietf.org/rfc/rfc5970.txt

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/insomniacslk/dhcp/iana"
)

// OptClientArchType represents an option CLIENT_ARCH_TYPE
type OptClientArchType struct {
	ArchTypes []iana.ArchType
}

func (op *OptClientArchType) Code() OptionCode {
	return OptionClientArchType
}

func (op *OptClientArchType) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionClientArchType))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	u16 := make([]byte, 2)
	for _, at := range op.ArchTypes {
		binary.BigEndian.PutUint16(u16, uint16(at))
		buf = append(buf, u16...)
	}
	return buf
}

func (op *OptClientArchType) Length() int {
	return 2*len(op.ArchTypes)
}

func (op *OptClientArchType) String() string {
	atStrings := make([]string, 0)
	for _, at := range op.ArchTypes {
		name := iana.ArchTypeToString(at)
		atStrings = append(atStrings, name)
	}
	return fmt.Sprintf("OptClientArchType{archtype=%v}", strings.Join(atStrings, ", "))
}

// ParseOptClientArchType builds an OptClientArchType structure from
// a sequence of bytes The input data does not include option code and
// length bytes.
func ParseOptClientArchType(data []byte) (*OptClientArchType, error) {
	opt := OptClientArchType{}
	if len(data) == 0 || len(data)%2 != 0 {
		return nil, fmt.Errorf("Invalid arch type data length. Expected multiple of 2 larger than 2, got %v", len(data))
	}
	for idx := 0; idx < len(data); idx += 2 {
		b := data[idx : idx+2]
		opt.ArchTypes = append(opt.ArchTypes, iana.ArchType(binary.BigEndian.Uint16(b)))
	}
	return &opt, nil
}
