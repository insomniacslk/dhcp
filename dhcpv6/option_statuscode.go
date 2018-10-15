package dhcpv6

// This module defines the OptStatusCode structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"

	"github.com/insomniacslk/dhcp/iana"
)

// OptStatusCode represents a DHCPv6 Status Code option
type OptStatusCode struct {
	StatusCode    iana.StatusCode
	StatusMessage []byte
}

// Code returns the option code
func (op *OptStatusCode) Code() OptionCode {
	return OptionStatusCode
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptStatusCode) ToBytes() []byte {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionStatusCode))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint16(buf[4:6], uint16(op.StatusCode))
	buf = append(buf, op.StatusMessage...)
	return buf
}

// Length returns the option length
func (op *OptStatusCode) Length() int {
	return 2 + len(op.StatusMessage)
}

func (op *OptStatusCode) String() string {
	return fmt.Sprintf("OptStatusCode{code=%s (%d), message=%v}",
		op.StatusCode.String(), op.StatusCode,
		string(op.StatusMessage))
}

// ParseOptStatusCode builds an OptStatusCode structure from a sequence of
// bytes. The input data does not include option code and length bytes.
func ParseOptStatusCode(data []byte) (*OptStatusCode, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("Invalid OptStatusCode data: length is shorter than 2")
	}
	opt := OptStatusCode{}
	opt.StatusCode = iana.StatusCode(binary.BigEndian.Uint16(data[0:2]))
	opt.StatusMessage = append(opt.StatusMessage, data[2:]...)
	return &opt, nil
}
