package dhcpv6

import (
	"errors"
	"fmt"
)

var errBufferNotLengthOne = errors.New("pref-value must be exactly 1 byte")

// OptPreference represents a Preference option as defined by RFC 9915
// Section 21.8.
func OptPreference(prefValue uint8) Option {
	return &optPreference{prefValue}
}

type optPreference struct {
	prefValue uint8
}

func (*optPreference) Code() OptionCode {
	return OptionPreference
}

func (op *optPreference) String() string {
	return fmt.Sprintf("%s: %v", op.Code(), op.prefValue)
}

// FromBytes builds an optPreference structure from a sequence of bytes. The
// input data does not include option code and length bytes.
func (op *optPreference) FromBytes(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("buffer is length %d, %w", len(data), errBufferNotLengthOne)
	}
	op.prefValue = data[0]
	return nil
}

func (op *optPreference) ToBytes() []byte {
	return []byte{op.prefValue}
}
