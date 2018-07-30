package dhcpv6

// This module defines the OptClientId and DUID structures.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

// OptClientId represents a Client ID option
type OptClientId struct {
	Cid Duid
}

func (op *OptClientId) Code() OptionCode {
	return OptionClientID
}

func (op *OptClientId) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionClientID))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	buf = append(buf, op.Cid.ToBytes()...)
	return buf
}

func (op *OptClientId) Length() int {
	return op.Cid.Length()
}

func (op *OptClientId) String() string {
	return fmt.Sprintf("OptClientId{cid=%v}", op.Cid.String())
}

// ParseOptClientId builds an OptClientId structure from a sequence
// of bytes. The input data does not include option code and length
// bytes.
func ParseOptClientId(data []byte) (*OptClientId, error) {
	if len(data) < 2 {
		// at least the DUID type is necessary to continue
		return nil, fmt.Errorf("Invalid OptClientId data: shorter than 2 bytes")
	}
	opt := OptClientId{}
	cid, err := DuidFromBytes(data)
	if err != nil {
		return nil, err
	}
	opt.Cid = *cid
	return &opt, nil
}
