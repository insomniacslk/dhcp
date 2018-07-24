package dhcpv4

import (
	"errors"
	"fmt"
	"strings"
)

// OptUserClass represents a DHCPv4 User Class option
type OptUserClass struct {
	UserClasses [][]byte
}

// Code returns the option code
func (op *OptUserClass) Code() OptionCode {
	return OptionUserClassInformation
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptUserClass) ToBytes() []byte {
	buf := []byte{byte(op.Code()), byte(op.Length())}
	for _, uc := range op.UserClasses {
		buf = append(buf, byte(len(uc)))
		buf = append(buf, uc...)
	}
	return buf
}

// Length returns the option length
func (op *OptUserClass) Length() int {
	ret := 0
	for _, uc := range op.UserClasses {
		ret += 1 + len(uc)
	}
	return ret
}

func (op *OptUserClass) String() string {
	ucStrings := make([]string, 0)
	for _, uc := range op.UserClasses {
		ucStrings = append(ucStrings, string(uc))
	}
	return fmt.Sprintf("OptUserClass{userclass=[%s]}", strings.Join(ucStrings, ", "))
}

// ParseOptUserClass returns a new OptUserClass from a byte stream or
// error if any
func ParseOptUserClass(data []byte) (*OptUserClass, error) {
	opt := OptUserClass{}

	if len(data) < 4 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionUserClassInformation {
		return nil, fmt.Errorf("expected code %v, got %v", OptionUserClassInformation, code)
	}
	data = data[2:]

	for {
		if len(data) == 0 {
			break
		}
		if len(data) < 1 {
			return nil, errors.New("ParseOptUserClass: short data: missing length field")
		}
		ucLen := int(data[0])
		if len(data) < ucLen+1 {
			return nil, fmt.Errorf("ParseOptUserClass: short data: less than %d bytes", ucLen+2)
		}
		opt.UserClasses = append(opt.UserClasses, data[1:ucLen+1])
		data = data[1+ucLen:]
	}
	if len(opt.UserClasses) < 1 {
		return nil, errors.New("ParseOptUserClass: at least one user class is required")
	}
	return &opt, nil
}
