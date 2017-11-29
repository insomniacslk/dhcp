package options

// This module defines the OptServerId and DUID structures.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

type OptServerId struct {
	sid Duid
}

func (op *OptServerId) Code() OptionCode {
	return OPTION_SERVERID
}

func (op *OptServerId) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_SERVERID))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	buf = append(buf, op.sid.ToBytes()...)
	return buf
}

func (op *OptServerId) ServerID() Duid {
	return op.sid
}

func (op *OptServerId) SetServerID(sid Duid) {
	op.sid = sid
}

func (op *OptServerId) Length() int {
	return op.sid.Length()
}

func (op *OptServerId) String() string {
	return fmt.Sprintf("OptServerId{sid=%v}", op.sid.String())
}

// build an OptServerId structure from a sequence of bytes.
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
	opt.sid = *sid
	return &opt, nil
}
