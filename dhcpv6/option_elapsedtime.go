package dhcpv6

// This module defines the OptElapsedTime structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

type OptElapsedTime struct {
	ElapsedTime uint16
}

func (op *OptElapsedTime) Code() OptionCode {
	return OPTION_ELAPSED_TIME
}

func (op *OptElapsedTime) ToBytes() []byte {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_ELAPSED_TIME))
	binary.BigEndian.PutUint16(buf[2:4], 2)
	binary.BigEndian.PutUint16(buf[4:6], uint16(op.ElapsedTime))
	return buf
}

func (op *OptElapsedTime) Length() int {
	return 2
}

func (op *OptElapsedTime) String() string {
	return fmt.Sprintf("OptElapsedTime{elapsedtime=%v}", op.ElapsedTime)
}

// build an OptElapsedTime structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptElapsedTime(data []byte) (*OptElapsedTime, error) {
	opt := OptElapsedTime{}
	if len(data) != 2 {
		return nil, fmt.Errorf("Invalid elapsed time data length. Expected 2 bytes, got %v", len(data))
	}
	opt.ElapsedTime = binary.BigEndian.Uint16(data)
	return &opt, nil
}
