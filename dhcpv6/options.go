package dhcpv6

import (
	"encoding/binary"
	"fmt"
)

type OptionCode uint16

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
		return nil, fmt.Errorf("Invalid option length. Declared %v, actual %v",
			length, len(dataStart)-4,
		)
	}
	var (
		err error
		opt Option
	)
	optData := dataStart[4 : 4+length]
	switch code {
	case OPTION_CLIENTID:
		opt, err = ParseOptClientId(optData)
	case OPTION_SERVERID:
		opt, err = ParseOptServerId(optData)
	case OPTION_ELAPSED_TIME:
		opt, err = ParseOptElapsedTime(optData)
	case OPTION_ORO:
		opt, err = ParseOptRequestedOption(optData)
	case DNS_RECURSIVE_NAME_SERVER:
		opt, err = ParseOptDNSRecursiveNameServer(optData)
	case DOMAIN_SEARCH_LIST:
		opt, err = ParseOptDomainSearchList(optData)
	case OPTION_IA_NA:
		opt, err = ParseOptIANA(optData)
	case OPTION_IA_PD:
		opt, err = ParseOptIAForPrefixDelegation(optData)
	case OPTION_IAADDR:
		opt, err = ParseOptIAAddress(optData)
	case OPTION_IAPREFIX:
		opt, err = ParseOptIAPrefix(optData)
	case OPTION_STATUS_CODE:
		opt, err = ParseOptStatusCode(optData)
	case OPTION_RELAY_MSG:
		opt, err = ParseOptRelayMsg(optData)
	default:
		opt = &OptionGeneric{OptionCode: code, OptionData: optData}
	}
	if err != nil {
		return nil, err
	}
	return opt, nil
}

func OptionsFromBytes(data []byte) ([]Option, error) {
	// Parse a sequence of bytes until the end and build a list of options from
	// it. Returns an error if any invalid option or length is found.
	if len(data) < 4 {
		// cannot be shorter than option code (2 bytes) + length (2 bytes)
		return nil, fmt.Errorf("Invalid options: shorter than 4 bytes")
	}
	options := make([]Option, 0, 10)
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
