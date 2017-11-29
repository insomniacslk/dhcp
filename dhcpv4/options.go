package dhcpv4

import (
	"bytes"
	"errors"
	"fmt"
)

type OptionCode byte

var MagicCookie = []byte{99, 130, 83, 99}

type Option struct {
	Code OptionCode
	Data []byte
}

func ParseOption(dataStart []byte) (*Option, error) {
	// Parse a sequence of bytes as a single DHCPv4 option.
	// Returns the option code, its data, and an error if any.
	if len(dataStart) == 0 {
		return nil, errors.New("Invalid zero-length DHCPv4 option")
	}
	opt := OptionCode(dataStart[0])
	switch opt {
	case OptionPad, OptionEnd:
		return &Option{Code: opt, Data: []byte{}}, nil
	default:
		length := int(dataStart[1])
		if len(dataStart) < length+2 {
			return nil, errors.New(
				fmt.Sprintf("Invalid data length. Declared %v, actual %v",
					length, len(dataStart),
				))
		}
		data := dataStart[2 : 2+length]
		return &Option{Code: opt, Data: data}, nil
	}
}

func OptionsFromBytes(data []byte) ([]Option, error) {
	// Parse a sequence of bytes until the end and build a list of options from
	// it. The sequence must contain the Magic Cookie.
	// Returns an error if any invalid option or length is found.
	if len(data) < 4 {
		return nil, errors.New("Invalid options: shorter than 4 bytes")
	}
	if !bytes.Equal(data[:4], MagicCookie) {
		return nil, errors.New(fmt.Sprintf("Invalid Magic Cookie: %v", data[:4]))
	}
	options := make([]Option, 0, 10)
	idx := 4
	for {
		if idx == len(data) {
			break
		}
		if idx > len(data) {
			// this should never happen
			return nil, errors.New("Error: Reading past the end of options")
		}
		opt, err := ParseOption(data[idx:])
		idx++
		if err != nil {
			return nil, err
		}
		options = append(options, *opt)
		if len(opt.Data) > 0 {
			// options with zero length have no length byte, so here we handle the ones with
			// nonzero length
			idx++
		}
		idx += len(opt.Data)
	}
	return options, nil
}

func OptionsToBytes(options []Option) []byte {
	// Convert a list of options to a wire-format representation. This will
	// include the Magic Cookie
	ret := MagicCookie
	for _, opt := range options {
		ret = append(ret, opt.ToBytes()...)
	}
	return ret
}

func (o *Option) String() string {
	code, ok := OptionCodeToString[o.Code]
	if !ok {
		code = "Unknown"
	}
	return fmt.Sprintf("%v -> %v", code, o.Data)
}

func (o *Option) ToBytes() []byte {
	// Convert a single option to its wire-format representation
	ret := []byte{byte(o.Code)}
	if o.Code != OptionPad && o.Code != OptionEnd {
		ret = append(ret, byte(len(o.Data)))
	}
	return append(ret, o.Data...)
}
