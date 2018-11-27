package dhcpv4

import (
	"fmt"
)

// This option implements the TFTP server name option.
// https://tools.ietf.org/html/rfc2132

// OptTFTPServerName implements the TFTP server name option.
type OptTFTPServerName struct {
	TFTPServerName []byte
}

// Code returns the option code
func (op *OptTFTPServerName) Code() OptionCode {
	return OptionTFTPServerName
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptTFTPServerName) ToBytes() []byte {
	return append([]byte{byte(op.Code()), byte(op.Length())}, op.TFTPServerName...)
}

// Length returns the option length in bytes
func (op *OptTFTPServerName) Length() int {
	return len(op.TFTPServerName)
}

func (op *OptTFTPServerName) String() string {
	return fmt.Sprintf("TFTP Server Name -> %s", op.TFTPServerName)
}

// ParseOptTFTPServerName returns a new OptTFTPServerName from a byte stream or error if any
func ParseOptTFTPServerName(data []byte) (*OptTFTPServerName, error) {
	if len(data) < 3 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionTFTPServerName {
		return nil, fmt.Errorf("ParseOptTFTPServerName: invalid code: %v; want %v",
			code, OptionTFTPServerName)
	}
	length := int(data[1])
	if length < 1 {
		return nil, fmt.Errorf("TFTP server name has invalid length of %d", length)
	}
	TFTPServerNameData := data[2:]
	if len(TFTPServerNameData) < length {
		return nil, fmt.Errorf("ParseOptTFTPServerName: short data: %d bytes; want %d",
			len(TFTPServerNameData), length)
	}
	return &OptTFTPServerName{TFTPServerName: TFTPServerNameData[:length]}, nil
}
