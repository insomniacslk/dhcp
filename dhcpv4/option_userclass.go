package dhcpv4

import (
	"errors"
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/uio"
)

// This option implements the User Class option
// https://tools.ietf.org/html/rfc3004

// OptUserClass represents an option encapsulating User Classes.
type OptUserClass struct {
	UserClasses [][]byte
	Rfc3004     bool
}

// Code returns the option code
func (op *OptUserClass) Code() OptionCode {
	return OptionUserClassInformation
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptUserClass) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	if !op.Rfc3004 {
		buf.WriteBytes(op.UserClasses[0])
	} else {
		for _, uc := range op.UserClasses {
			buf.Write8(uint8(len(uc)))
			buf.WriteBytes(uc)
		}
	}
	return buf.Data()
}

func (op *OptUserClass) String() string {
	ucStrings := make([]string, 0, len(op.UserClasses))
	if !op.Rfc3004 {
		ucStrings = append(ucStrings, string(op.UserClasses[0]))
	} else {
		for _, uc := range op.UserClasses {
			ucStrings = append(ucStrings, string(uc))
		}
	}
	return fmt.Sprintf("User Class Information -> %v", strings.Join(ucStrings, ", "))
}

// ParseOptUserClass returns a new OptUserClass from a byte stream or
// error if any
func ParseOptUserClass(data []byte) (*OptUserClass, error) {
	opt := OptUserClass{}
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
		opt.UserClasses = append(opt.UserClasses, data)
		return &opt, nil
	}
	opt.Rfc3004 = true
	for buf.Has(1) {
		ucLen := buf.Read8()
		if ucLen == 0 {
			return nil, fmt.Errorf("DHCP user class must have length greater than 0")
		}
		opt.UserClasses = append(opt.UserClasses, buf.CopyN(int(ucLen)))
	}
	if len(opt.UserClasses) == 0 {
		return nil, errors.New("ParseOptUserClass: at least one user class is required")
	}
	return &opt, buf.FinError()
}
