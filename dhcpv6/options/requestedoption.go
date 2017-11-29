package options

// This module defines the OptRequestedOption structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

type OptRequestedOption struct {
	requestedOptions []OptionCode
}

func (op *OptRequestedOption) Code() OptionCode {
	return OPTION_ORO
}

func (op *OptRequestedOption) ToBytes() []byte {
	buf := make([]byte, 4)
	roBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_ORO))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	for _, ro := range op.requestedOptions {
		binary.BigEndian.PutUint16(roBytes, uint16(ro))
		buf = append(buf, roBytes...)
	}
	return buf
}

func (op *OptRequestedOption) RequestedOptions() []OptionCode {
	return op.requestedOptions
}

func (op *OptRequestedOption) SetRequestedOptions(opts []OptionCode) {
	op.requestedOptions = opts
}

func (op *OptRequestedOption) Length() int {
	return len(op.requestedOptions) * 2
}

func (op *OptRequestedOption) String() string {
	roString := "["
	for idx, code := range op.requestedOptions {
		if name, ok := OptionCodeToString[OptionCode(code)]; ok {
			roString += name
		} else {
			roString += "Unknown"
		}
		if idx < len(op.requestedOptions)-1 {
			roString += ", "
		}
	}
	roString += "]"
	return fmt.Sprintf("OptRequestedOption{options=%v}", roString)
}

// build an OptRequestedOption structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptRequestedOption(data []byte) (*OptRequestedOption, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("Invalid OptRequestedOption data: length is not a multiple of 2")
	}
	opt := OptRequestedOption{}
	var rOpts []OptionCode
	for i := 0; i < len(data); i += 2 {
		rOpts = append(rOpts, OptionCode(binary.BigEndian.Uint16(data[i:i+2])))
	}
	opt.requestedOptions = rOpts
	return &opt, nil
}
