package dhcpv6

// This module defines the OptServerId and DUID structures.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

// OptServerId represents a Client ID option
type OptServerId struct {
	Sid Duid
}

func (op *OptServerId) Code() OptionCode {
	return OPTION_SERVERID
}

func (op *OptServerId) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_SERVERID))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	buf = append(buf, op.Sid.ToBytes()...)
	return buf
}

func (op *OptServerId) Length() int {
	return op.Sid.Length()
}

func (op *OptServerId) String() string {
	return fmt.Sprintf("OptServerId{sid=%v}", op.Sid.String())
}

// ParseOptServerId builds an OptServerId structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptServerId(data []byte) (*OptServerId, error) {
	if len(data) < 2 {
		// at least the DUID type is necessary to continue
		return nil, fmt.Errorf("Invalid OptServerId data: shorter than 2 bytes")
	}
	opt := OptServerId{}
	sid, err := DuidFromBytes(data)
	if err != nil {
		return nil, err
	}
	opt.Sid = *sid
	return &opt, nil
}
