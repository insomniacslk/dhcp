package dhcpv6

import (
	"fmt"
	"strings"

	"github.com/u-root/uio/uio"
)

// Option is an interface that all DHCPv6 options adhere to.
type Option interface {
	Code() OptionCode
	ToBytes() []byte
	String() string
	FromBytes([]byte) error
}

type OptionGeneric struct {
	OptionCode OptionCode
	OptionData []byte
}

func (og *OptionGeneric) Code() OptionCode {
	return og.OptionCode
}

func (og *OptionGeneric) ToBytes() []byte {
	return og.OptionData
}

func (og *OptionGeneric) String() string {
	if len(og.OptionData) == 0 {
		return og.OptionCode.String()
	}
	return fmt.Sprintf("%s: %v", og.OptionCode, og.OptionData)
}

// FromBytes resets OptionData to p.
func (og *OptionGeneric) FromBytes(p []byte) error {
	og.OptionData = append([]byte(nil), p...)
	return nil
}

// ParseOption parses data according to the given code.
func ParseOption(code OptionCode, optData []byte) (Option, error) {
	// Parse a sequence of bytes as a single DHCPv6 option.
	// Returns the option structure, or an error if any.
	var (
		err error
		opt Option
	)
	switch code {
	case OptionClientID:
		var o optClientID
		err = o.FromBytes(optData)
		opt = &o
	case OptionServerID:
		var o optServerID
		err = o.FromBytes(optData)
		opt = &o
	case OptionIANA:
		var o OptIANA
		err = o.FromBytes(optData)
		opt = &o
	case OptionIATA:
		var o OptIATA
		err = o.FromBytes(optData)
		opt = &o
	case OptionIAAddr:
		var o OptIAAddress
		err = o.FromBytes(optData)
		opt = &o
	case OptionORO:
		var o optRequestedOption
		err = o.FromBytes(optData)
		opt = &o
	case OptionElapsedTime:
		var o optElapsedTime
		err = o.FromBytes(optData)
		opt = &o
	case OptionRelayMsg:
		var o optRelayMsg
		err = o.FromBytes(optData)
		opt = &o
	case OptionStatusCode:
		var o OptStatusCode
		err = o.FromBytes(optData)
		opt = &o
	case OptionUserClass:
		var o OptUserClass
		err = o.FromBytes(optData)
		opt = &o
	case OptionVendorClass:
		var o OptVendorClass
		err = o.FromBytes(optData)
		opt = &o
	case OptionVendorOpts:
		var o OptVendorOpts
		err = o.FromBytes(optData)
		opt = &o
	case OptionInterfaceID:
		var o optInterfaceID
		err = o.FromBytes(optData)
		opt = &o
	case OptionDNSRecursiveNameServer:
		var o optDNS
		err = o.FromBytes(optData)
		opt = &o
	case OptionDomainSearchList:
		var o optDomainSearchList
		err = o.FromBytes(optData)
		opt = &o
	case OptionIAPD:
		var o OptIAPD
		err = o.FromBytes(optData)
		opt = &o
	case OptionIAPrefix:
		var o OptIAPrefix
		err = o.FromBytes(optData)
		opt = &o
	case OptionInformationRefreshTime:
		var o optInformationRefreshTime
		err = o.FromBytes(optData)
		opt = &o
	case OptionRemoteID:
		var o OptRemoteID
		err = o.FromBytes(optData)
		opt = &o
	case OptionFQDN:
		var o OptFQDN
		err = o.FromBytes(optData)
		opt = &o
	case OptionNTPServer:
		var o OptNTPServer
		err = o.FromBytes(optData)
		opt = &o
	case OptionBootfileURL:
		var o optBootFileURL
		err = o.FromBytes(optData)
		opt = &o
	case OptionBootfileParam:
		var o optBootFileParam
		err = o.FromBytes(optData)
		opt = &o
	case OptionClientArchType:
		var o optClientArchType
		err = o.FromBytes(optData)
		opt = &o
	case OptionNII:
		var o OptNetworkInterfaceID
		err = o.FromBytes(optData)
		opt = &o
	case OptionClientLinkLayerAddr:
		var o optClientLinkLayerAddress
		err = o.FromBytes(optData)
		opt = &o
	case OptionDHCPv4Msg:
		var o OptDHCPv4Msg
		err = o.FromBytes(optData)
		opt = &o
	case OptionDHCP4oDHCP6Server:
		var o OptDHCP4oDHCP6Server
		err = o.FromBytes(optData)
		opt = &o
	case Option4RD:
		var o Opt4RD
		err = o.FromBytes(optData)
		opt = &o
	case Option4RDMapRule:
		var o Opt4RDMapRule
		err = o.FromBytes(optData)
		opt = &o
	case Option4RDNonMapRule:
		var o Opt4RDNonMapRule
		err = o.FromBytes(optData)
		opt = &o
	case OptionRelayPort:
		var o optRelayPort
		err = o.FromBytes(optData)
		opt = &o
	default:
		opt = &OptionGeneric{OptionCode: code, OptionData: optData}
	}
	if err != nil {
		return nil, err
	}
	return opt, nil
}

type longStringer interface {
	LongString(spaceIndent int) string
}

// Options is a collection of options.
type Options []Option

// LongString prints options with indentation of at least spaceIndent spaces.
func (o Options) LongString(spaceIndent int) string {
	indent := strings.Repeat(" ", spaceIndent)
	var s strings.Builder
	if len(o) == 0 {
		s.WriteString("[]")
	} else {
		s.WriteString("[\n")
		for _, opt := range o {
			s.WriteString(indent)
			s.WriteString("  ")
			if ls, ok := opt.(longStringer); ok {
				s.WriteString(ls.LongString(spaceIndent + 2))
			} else {
				s.WriteString(opt.String())
			}
			s.WriteString("\n")
		}
		s.WriteString(indent)
		s.WriteString("]")
	}
	return s.String()
}

// Get returns all options matching the option code.
func (o Options) Get(code OptionCode) []Option {
	var ret []Option
	for _, opt := range o {
		if opt.Code() == code {
			ret = append(ret, opt)
		}
	}
	return ret
}

// GetOne returns the first option matching the option code.
func (o Options) GetOne(code OptionCode) Option {
	for _, opt := range o {
		if opt.Code() == code {
			return opt
		}
	}
	return nil
}

// Add appends one option.
func (o *Options) Add(option Option) {
	*o = append(*o, option)
}

// Del deletes all options matching the option code.
func (o *Options) Del(code OptionCode) {
	newOpts := make(Options, 0, len(*o))
	for _, opt := range *o {
		if opt.Code() != code {
			newOpts = append(newOpts, opt)
		}
	}
	*o = newOpts
}

// Update replaces the first option of the same type as the specified one.
func (o *Options) Update(option Option) {
	for idx, opt := range *o {
		if opt.Code() == option.Code() {
			(*o)[idx] = option
			// don't look further
			return
		}
	}
	// if not found, add it
	o.Add(option)
}

// ToBytes marshals all options to bytes.
func (o Options) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, opt := range o {
		buf.Write16(uint16(opt.Code()))

		val := opt.ToBytes()
		buf.Write16(uint16(len(val)))
		buf.WriteBytes(val)
	}
	return buf.Data()
}

// FromBytes reads data into o and returns an error if the options are not a
// valid serialized representation of DHCPv6 options per RFC 3315.
func (o *Options) FromBytes(data []byte) error {
	return o.FromBytesWithParser(data, ParseOption)
}

// OptionParser is a function signature for option parsing
type OptionParser func(code OptionCode, data []byte) (Option, error)

// FromBytesWithParser parses Options from byte sequences using the parsing
// function that is passed in as a paremeter
func (o *Options) FromBytesWithParser(data []byte, parser OptionParser) error {
	*o = make(Options, 0, 10)
	if len(data) == 0 {
		// no options, no party
		return nil
	}

	buf := uio.NewBigEndianBuffer(data)
	for buf.Has(4) {
		code := OptionCode(buf.Read16())
		length := int(buf.Read16())

		// Consume, but do not Copy. Each parser will make a copy of
		// pertinent data.
		optData := buf.Consume(length)

		opt, err := parser(code, optData)
		if err != nil {
			return err
		}
		*o = append(*o, opt)
	}
	return buf.FinError()
}
