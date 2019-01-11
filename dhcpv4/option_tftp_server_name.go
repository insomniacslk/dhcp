package dhcpv4

import (
	"fmt"
)

// OptTFTPServerName implements the TFTP server name option described by RFC
// 2132, Section 9.4.
type OptTFTPServerName struct {
	TFTPServerName string
}

// Code returns the option code
func (op *OptTFTPServerName) Code() OptionCode {
	return OptionTFTPServerName
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptTFTPServerName) ToBytes() []byte {
	return []byte(op.TFTPServerName)
}

func (op *OptTFTPServerName) String() string {
	return fmt.Sprintf("TFTP Server Name -> %s", op.TFTPServerName)
}

// ParseOptTFTPServerName returns a new OptTFTPServerName from a byte stream or error if any
func ParseOptTFTPServerName(data []byte) (*OptTFTPServerName, error) {
	return &OptTFTPServerName{TFTPServerName: string(data)}, nil
}
