package dhcpv6

// This module defines the OptStatusCode structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

// OptStatusCode represents a DHCPv6 Status Code option
type OptStatusCode struct {
	statusCode    uint16
	statusMessage []byte
}

// Code returns the option code
func (op *OptStatusCode) Code() OptionCode {
	return OPTION_STATUS_CODE
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptStatusCode) ToBytes() []byte {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_STATUS_CODE))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint16(buf[4:6], op.statusCode)
	buf = append(buf, op.statusMessage...)
	return buf
}

// StatusCode returns the status code
func (op *OptStatusCode) StatusCode() uint16 {
	return op.statusCode
}

// SetStatusCode sets the status code
func (op *OptStatusCode) SetStatusCode(code uint16) {
	op.statusCode = code
}

// StatusMessage returns the status message
func (op *OptStatusCode) StatusMessage() []byte {
	return op.statusMessage
}

// SetStatusMessage sets the status message
func (op *OptStatusCode) SetStatusMessage(message []byte) {
	op.statusMessage = message
}

// Length returns the option length
func (op *OptStatusCode) Length() int {
	return 2 + len(op.statusMessage)
}

func (op *OptStatusCode) String() string {
	return fmt.Sprintf("OptStatusCode{code=%v, message=%v}", op.statusCode, string(op.statusMessage))
}

// ParseOptStatusCode builds an OptStatusCode structure from a sequence of
// bytes. The input data does not include option code and length bytes.
func ParseOptStatusCode(data []byte) (*OptStatusCode, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("Invalid OptStatusCode data: length is shorter than 2")
	}
	opt := OptStatusCode{}
	opt.statusCode = binary.BigEndian.Uint16(data[0:2])
	opt.statusMessage = append(opt.statusMessage, data[2:]...)
	return &opt, nil
}
