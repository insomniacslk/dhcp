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

// OptionsFromBytesWithMagicCookie parses a sequence of bytes until the end and
// builds a list of options from it. The sequence must contain the Magic Cookie.
// Returns an error if any invalid option or length is found.
func OptionsFromBytesWithMagicCookie(data []byte) ([]Option, error) {
	if len(data) < len(MagicCookie) {
		return nil, errors.New("Invalid options: shorter than 4 bytes")
	}
	if !bytes.Equal(data[:len(MagicCookie)], MagicCookie) {
		return nil, fmt.Errorf("Invalid Magic Cookie: %v", data[:len(MagicCookie)])
	}
	opts, err := OptionsFromBytes(data[len(MagicCookie):])
	if err != nil {
		return nil, err
	}
	return opts, nil
}

// OptionsFromBytes parses a sequence of bytes until the end and builds a list
// of options from it. The sequence should not contain the DHCP magic cookie.
// Returns an error if any invalid option or length is found.
func OptionsFromBytes(data []byte) ([]Option, error) {
	options := make([]Option, 0, 10)
	idx := 0
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

		// Options with zero length have no length byte, so here we handle the
		// ones with nonzero length
		if len(opt.Data) > 0 {
			idx++
		}
		idx += len(opt.Data)
	}
	return options, nil
}

// OptionsToBytesWithMagicCookie converts a list of options to a wire-format
// representation with the DHCP magic cookie prepended.
func OptionsToBytesWithMagicCookie(options []Option) []byte {
	ret := MagicCookie
	return append(ret, OptionsToBytes(options)...)
}

// OptionsToBytes converts a list of options to a wire-format representation.
func OptionsToBytes(options []Option) []byte {
	ret := []byte{}
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

// BSDPString converts a BSDP-specific option embedded in
// vendor-specific information to a human-readable string.
func (o *Option) BSDPString() string {
	code, ok := BSDPOptionCodeToString[o.Code]
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
