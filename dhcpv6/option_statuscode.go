package dhcpv6

// This module defines the OptStatusCode structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

type OptStatusCode struct {
	statusCode    uint16
	statusMessage []byte
}

func (op *OptStatusCode) Code() OptionCode {
	return OPTION_STATUS_CODE
}

func (op *OptStatusCode) ToBytes() []byte {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_STATUS_CODE))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint16(buf[4:6], op.statusCode)
	buf = append(buf, op.statusMessage...)
	return buf
}

func (op *OptStatusCode) StatusCode() uint16 {
	return op.statusCode
}

func (op *OptStatusCode) SetStatusCode(code uint16) {
	op.statusCode = code
}

func (op *OptStatusCode) StatusMessage() uint16 {
	return op.statusCode
}

func (op *OptStatusCode) SetStatusMessage(message []byte) {
	op.statusMessage = message
}

func (op *OptStatusCode) Length() int {
	return 2 + len(op.statusMessage)
}

func (op *OptStatusCode) String() string {
	return fmt.Sprintf("OptStatusCode{code=%v, message=%v}", op.statusCode, string(op.statusMessage))
}

// build an OptStatusCode structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptStatusCode(data []byte) (*OptStatusCode, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("Invalid OptStatusCode data: length is shorter than 2")
	}
	opt := OptStatusCode{}
	opt.statusCode = binary.BigEndian.Uint16(data[0:2])
	opt.statusMessage = append(opt.statusMessage, data[2:]...)
	return &opt, nil
}
