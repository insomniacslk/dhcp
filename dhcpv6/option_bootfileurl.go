package dhcpv6

import (
	"fmt"
)

// OptBootFileURL returns a OptionBootfileURL as defined by RFC 5970.
func OptBootFileURL(url string) Option {
	return &optBootFileURL{url}
}

type String string

// ToBytes serializes the option and returns it as a sequence of bytes
func (s String) ToBytes() []byte {
	return []byte(s)
}

func (s String) String() string {
	return string(s)
}

// FromBytes builds an String structure from a sequence of bytes. The input
// data does not include option code and length bytes.
func (s *String) FromBytes(data []byte) error {
	*s = String(string(data))
	return nil
}

type optBootFileURL struct {
	url string
}

// Code returns the option code
func (op optBootFileURL) Code() OptionCode {
	return OptionBootfileURL
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op optBootFileURL) ToBytes() []byte {
	return []byte(op.url)
}

func (op optBootFileURL) String() string {
	return fmt.Sprintf("%s: %s", op.Code(), op.url)
}

// FromBytes builds an optBootFileURL structure from a sequence
// of bytes. The input data does not include option code and length bytes.
func (op *optBootFileURL) FromBytes(data []byte) error {
	op.url = string(data)
	return nil
}
