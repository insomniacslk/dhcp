package options

// This module defines the OptClientId and DUID structures.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

type OptClientId struct {
	cid Duid
}

func (op *OptClientId) Code() OptionCode {
	return OPTION_CLIENTID
}

func (op *OptClientId) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_CLIENTID))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	buf = append(buf, op.cid.ToBytes()...)
	return buf
}

func (op *OptClientId) ClientID() Duid {
	return op.cid
}

func (op *OptClientId) SetClientID(cid Duid) {
	op.cid = cid
}

func (op *OptClientId) Length() int {
	return op.cid.Length()
}

func (op *OptClientId) String() string {
	return fmt.Sprintf("OptClientId{cid=%v}", op.cid.String())
}

// build an OptClientId structure from a sequence of bytes.
// The input data does not include option code and length bytes.
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
	opt.cid = *cid
	return &opt, nil
}
