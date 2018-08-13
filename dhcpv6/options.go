package dhcpv6

import (
	"encoding/binary"
	"fmt"
)

// OptionCode is a single byte representing the code for a given Option.
type OptionCode uint16

// Option is an interface that all DHCPv6 options adhere to.
type Option interface {
	Code() OptionCode
	ToBytes() []byte
	Length() int
	String() string
}

type OptionGeneric struct {
	OptionCode OptionCode
	OptionData []byte
}

func (og *OptionGeneric) Code() OptionCode {
	return og.OptionCode
}

func (og *OptionGeneric) ToBytes() []byte {
	var ret []byte
	codeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(codeBytes, uint16(og.OptionCode))
	ret = append(ret, codeBytes...)
	lengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthBytes, uint16(len(og.OptionData)))
	ret = append(ret, lengthBytes...)
	ret = append(ret, og.OptionData...)
	return ret
}

func (og *OptionGeneric) String() string {
	code, ok := OptionCodeToString[og.OptionCode]
	if !ok {
		code = "UnknownOption"
	}
	return fmt.Sprintf("%v -> %v", code, og.OptionData)
}

func (og *OptionGeneric) Length() int {
	return len(og.OptionData)
}

func ParseOption(dataStart []byte) (Option, error) {
	// Parse a sequence of bytes as a single DHCPv6 option.
	// Returns the option structure, or an error if any.
	if len(dataStart) < 4 {
		return nil, fmt.Errorf("Invalid DHCPv6 option: less than 4 bytes")
	}
	code := OptionCode(binary.BigEndian.Uint16(dataStart[:2]))
	length := int(binary.BigEndian.Uint16(dataStart[2:4]))
	if len(dataStart) < length+4 {
		return nil, fmt.Errorf("Invalid option length for option %v. Declared %v, actual %v",
			code, length, len(dataStart)-4,
		)
	}
	var (
		err error
		opt Option
	)
	optData := dataStart[4 : 4+length]
	switch code {
	case OptionClientID:
		opt, err = ParseOptClientId(optData)
	case OptionServerID:
		opt, err = ParseOptServerId(optData)
	case OptionIANA:
		opt, err = ParseOptIANA(optData)
	case OptionIAAddr:
		opt, err = ParseOptIAAddress(optData)
	case OptionORO:
		opt, err = ParseOptRequestedOption(optData)
	case OptionElapsedTime:
		opt, err = ParseOptElapsedTime(optData)
	case OptionRelayMsg:
		opt, err = ParseOptRelayMsg(optData)
	case OptionStatusCode:
		opt, err = ParseOptStatusCode(optData)
	case OptionUserClass:
		opt, err = ParseOptUserClass(optData)
	case OptionVendorClass:
		opt, err = ParseOptVendorClass(optData)
	case OptionInterfaceID:
		opt, err = ParseOptInterfaceId(optData)
	case OptionDNSRecursiveNameServer:
		opt, err = ParseOptDNSRecursiveNameServer(optData)
	case OptionDomainSearchList:
		opt, err = ParseOptDomainSearchList(optData)
	case OptionIAPD:
		opt, err = ParseOptIAForPrefixDelegation(optData)
	case OptionIAPrefix:
		opt, err = ParseOptIAPrefix(optData)
	case OptionRemoteID:
		opt, err = ParseOptRemoteId(optData)
	case OptionBootfileURL:
		opt, err = ParseOptBootFileURL(optData)
	case OptionClientArchType:
		opt, err = ParseOptClientArchType(optData)
	case OptionNII:
		opt, err = ParseOptNetworkInterfaceId(optData)
	case OptionVendorOpts:
		opt, err = ParseOptVendorOpts(optData)
	default:
		opt = &OptionGeneric{OptionCode: code, OptionData: optData}
	}
	if err != nil {
		return nil, err
	}
	if length != opt.Length() {
		return nil, fmt.Errorf("Error: declared length is different from actual length for option %d: %d != %d",
			code, opt.Length(), length)
	}
	return opt, nil
}

func OptionsFromBytes(data []byte) ([]Option, error) {
	// Parse a sequence of bytes until the end and build a list of options from
	// it. Returns an error if any invalid option or length is found.
	options := make([]Option, 0, 10)
	if len(data) == 0 {
		// no options, no party
		return options, nil
	}
	if len(data) < 4 {
		// cannot be shorter than option code (2 bytes) + length (2 bytes)
		return nil, fmt.Errorf("Invalid options: shorter than 4 bytes")
	}
	idx := 0
	for {
		if idx == len(data) {
			break
		}
		if idx > len(data) {
			// this should never happen
			return nil, fmt.Errorf("Error: reading past the end of options")
		}
		opt, err := ParseOption(data[idx:])
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
		idx += opt.Length() + 4 // 4 bytes for type + length
	}
	return options, nil
}
