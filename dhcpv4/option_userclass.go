package dhcpv4

import (
	"errors"
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/uio"
)

// UserClass implements the user class option described by RFC 3004.
type UserClass struct {
	UserClasses [][]byte
	RFC3004     bool
}

// GetUserClass returns the user class in o if present.
//
// The user class information option is defined by RFC 3004.
func GetUserClass(o Options) *UserClass {
	v := o.Get(OptionUserClassInformation)
	if v == nil {
		return nil
	}
	var uc UserClass
	if err := uc.FromBytes(v); err != nil {
		return nil
	}
	return &uc
}

// OptUserClass returns a new user class option.
func OptUserClass(v []byte) Option {
	return Option{
		Code: OptionUserClassInformation,
		Value: &UserClass{
			UserClasses: [][]byte{v},
			RFC3004:     false,
		},
	}
}

// OptRFC3004UserClass returns a new user class option according to RFC 3004.
func OptRFC3004UserClass(v [][]byte) Option {
	return Option{
		Code: OptionUserClassInformation,
		Value: &UserClass{
			UserClasses: v,
			RFC3004:     true,
		},
	}
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *UserClass) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	if !op.RFC3004 {
		buf.WriteBytes(op.UserClasses[0])
	} else {
		for _, uc := range op.UserClasses {
			buf.Write8(uint8(len(uc)))
			buf.WriteBytes(uc)
		}
	}
	return buf.Data()
}

// String returns a human-readable user class.
func (op *UserClass) String() string {
	ucStrings := make([]string, 0, len(op.UserClasses))
	if !op.RFC3004 {
		ucStrings = append(ucStrings, string(op.UserClasses[0]))
	} else {
		for _, uc := range op.UserClasses {
			ucStrings = append(ucStrings, string(uc))
		}
	}
	return strings.Join(ucStrings, ", ")
}

// FromBytes parses data into op.
func (op *UserClass) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)

	// Check if option is Microsoft style instead of RFC compliant, issue #113

	// User-class options are, according to RFC3004, supposed to contain a set
	// of strings each with length UC_Len_i. Here we check that this is so,
	// by seeing if all the UC_Len_i lengths are consistent with the overall
	// option length. If the lengths don't add up, we assume that the option
	// is a single string and non RFC3004 compliant
	var counting int
	for counting < buf.Len() {
		// UC_Len_i does not include itself so add 1
		counting += int(data[counting]) + 1
	}
	if counting != buf.Len() {
		op.UserClasses = append(op.UserClasses, data)
		return nil
	}
	op.RFC3004 = true
	for buf.Has(1) {
		ucLen := buf.Read8()
		if ucLen == 0 {
			return fmt.Errorf("DHCP user class must have length greater than 0")
		}
		op.UserClasses = append(op.UserClasses, buf.CopyN(int(ucLen)))
	}
	if len(op.UserClasses) == 0 {
		return errors.New("ParseOptUserClass: at least one user class is required")
	}
	return buf.FinError()
}
