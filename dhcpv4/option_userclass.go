package dhcpv4

import (
	"errors"
	"fmt"
	"strings"
)

// This option implements the User Class option
// https://tools.ietf.org/html/rfc3004

// OptUserClass represents an option encapsulating User Classes.
type OptUserClass struct {
	UserClasses [][]byte
	Rfc3004 bool
}

// Code returns the option code
func (op *OptUserClass) Code() OptionCode {
	return OptionUserClassInformation
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptUserClass) ToBytes() []byte {
	buf := []byte{byte(op.Code()), byte(op.Length())}
	if !op.Rfc3004 && len(op.UserClasses) == 1 {
		return append(buf, op.UserClasses[0]...)
	}
	for _, uc := range op.UserClasses {
		buf = append(buf, byte(len(uc)))
		buf = append(buf, uc...)
	}
	return buf
}

// Length returns the option length
func (op *OptUserClass) Length() int {
	ret := 0
	if !op.Rfc3004 && len(op.UserClasses) == 1 {
		return len(op.UserClasses[0])
	}
	for _, uc := range op.UserClasses {
		ret += 1 + len(uc)
	}
	return ret
}

func (op *OptUserClass) String() string {
	ucStrings := make([]string, 0, len(op.UserClasses))
	for _, uc := range op.UserClasses {
		ucStrings = append(ucStrings, string(uc))
	}
	return fmt.Sprintf("OptUserClass{userclass=[%s]}", strings.Join(ucStrings, ", "))
}

// ParseOptUserClass returns a new OptUserClass from a byte stream or
// error if any
func ParseOptUserClass(data []byte) (*OptUserClass, error) {
	opt := OptUserClass{}

	if len(data) < 3 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionUserClassInformation {
		return nil, fmt.Errorf("expected code %v, got %v", OptionUserClassInformation, code)
	}

	totalLength := int(data[1])
	data = data[2:]
	if len(data) < totalLength {
		return nil, fmt.Errorf("ParseOptUserClass: short data: length is %d but got %d bytes",
			totalLength, len(data))
	}

	// Check if option is Microsoft style instead of RFC compliant, issue #113

	// User-class options are, according to RFC3004, supposed to contain a set
	// of strings each with length UC_Len_i. Here we check that this is so,
	// by seeing if all the UC_Len_i lengths are consistent with the overall
	// option length. If the lengths don't add up, we assume that the option
	// is a single string and non RFC3004 compliant
	var counting int
	for counting < totalLength {
		// UC_Len_i does not include itself so add 1
		counting += int(data[counting]) + 1
	}
	if counting != totalLength {
		opt.UserClasses = append(opt.UserClasses, data[:totalLength])
		return &opt, nil
	}
	opt.Rfc3004 = true
	for i := 0; i < totalLength; {
		ucLen := int(data[i])
		if ucLen == 0 {
			return nil, errors.New("User Class value has invalid length of 0")
		}
		base := i + 1
		if len(data) < base+ucLen {
			return nil, fmt.Errorf("ParseOptUserClass: short data: %d bytes; want: %d", len(data), base+ucLen)
		}
		opt.UserClasses = append(opt.UserClasses, data[base:base+ucLen])
		i += base + ucLen
	}
	if len(opt.UserClasses) < 1 {
		return nil, errors.New("ParseOptUserClass: at least one user class is required")
	}
	return &opt, nil
}
