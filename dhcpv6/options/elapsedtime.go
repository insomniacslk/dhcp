package options

// This module defines the OptElapsedTime structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

type OptElapsedTime struct {
	elapsedTime uint16
}

func (op *OptElapsedTime) Code() OptionCode {
	return OPTION_ELAPSED_TIME
}

func (op *OptElapsedTime) ToBytes() []byte {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_ELAPSED_TIME))
	binary.BigEndian.PutUint16(buf[2:4], 2)
	binary.BigEndian.PutUint16(buf[4:6], uint16(op.elapsedTime))
	return buf
}

func (op *OptElapsedTime) ElapsedTime() uint16 {
	return op.elapsedTime
}

func (op *OptElapsedTime) SetElapsedTime(elapsedTime uint16) {
	op.elapsedTime = elapsedTime
}

func (op *OptElapsedTime) Length() int {
	return 2
}

func (op *OptElapsedTime) String() string {
	return fmt.Sprintf("OptElapsedTime{elapsedtime=%v}", op.elapsedTime)
}

// build an OptElapsedTime structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptElapsedTime(data []byte) (*OptElapsedTime, error) {
	opt := OptElapsedTime{}
	if len(data) != 2 {
		return nil, fmt.Errorf("Invalid elapsed time data length. Expected 2 bytes, got %v", len(data))
	}
	opt.elapsedTime = binary.BigEndian.Uint16(data)
	return &opt, nil
}
