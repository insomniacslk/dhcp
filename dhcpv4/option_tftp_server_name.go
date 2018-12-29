package dhcpv4

import (
	"fmt"
)

// This option implements the TFTP server name option.
// https://tools.ietf.org/html/rfc2132

// OptTFTPServerName implements the TFTP server name option.
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
