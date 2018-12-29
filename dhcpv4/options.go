package dhcpv4

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/u-root/u-root/pkg/uio"
)

var (
	// ErrShortByteStream is an error that is thrown any time a short byte stream is
	// detected during option parsing.
	ErrShortByteStream = errors.New("short byte stream")

	// ErrZeroLengthByteStream is an error that is thrown any time a zero-length
	// byte stream is encountered.
	ErrZeroLengthByteStream = errors.New("zero-length byte stream")

	// ErrInvalidOptions is returned when invalid options data is
	// encountered during parsing. The data could report an incorrect
	// length or have trailing bytes which are not part of the option.
	ErrInvalidOptions = errors.New("invalid options data")
)

// OptionCode is a single byte representing the code for a given Option.
type OptionCode byte

// Option is an interface that all DHCP v4 options adhere to.
type Option interface {
	Code() OptionCode
	ToBytes() []byte
	Length() int
	String() string
}

// ParseOption parses a sequence of bytes as a single DHCPv4 option, returning
// the specific option structure or error, if any.
func ParseOption(code OptionCode, data []byte) (Option, error) {
	var (
		opt Option
		err error
	)
	switch code {
	case OptionSubnetMask:
		opt, err = ParseOptSubnetMask(data)
	case OptionRouter:
		opt, err = ParseOptRouter(data)
	case OptionDomainNameServer:
		opt, err = ParseOptDomainNameServer(data)
	case OptionHostName:
		opt, err = ParseOptHostName(data)
	case OptionDomainName:
		opt, err = ParseOptDomainName(data)
	case OptionRootPath:
		opt, err = ParseOptRootPath(data)
	case OptionBroadcastAddress:
		opt, err = ParseOptBroadcastAddress(data)
	case OptionNTPServers:
		opt, err = ParseOptNTPServers(data)
	case OptionRequestedIPAddress:
		opt, err = ParseOptRequestedIPAddress(data)
	case OptionIPAddressLeaseTime:
		opt, err = ParseOptIPAddressLeaseTime(data)
	case OptionDHCPMessageType:
		opt, err = ParseOptMessageType(data)
	case OptionServerIdentifier:
		opt, err = ParseOptServerIdentifier(data)
	case OptionParameterRequestList:
		opt, err = ParseOptParameterRequestList(data)
	case OptionMaximumDHCPMessageSize:
		opt, err = ParseOptMaximumDHCPMessageSize(data)
	case OptionClassIdentifier:
		opt, err = ParseOptClassIdentifier(data)
	case OptionTFTPServerName:
		opt, err = ParseOptTFTPServerName(data)
	case OptionBootfileName:
		opt, err = ParseOptBootfileName(data)
	case OptionUserClassInformation:
		opt, err = ParseOptUserClass(data)
	case OptionRelayAgentInformation:
		opt, err = ParseOptRelayAgentInformation(data)
	case OptionClientSystemArchitectureType:
		opt, err = ParseOptClientArchType(data)
	case OptionDNSDomainSearchList:
		opt, err = ParseOptDomainSearch(data)
	case OptionVendorIdentifyingVendorClass:
		opt, err = ParseOptVIVC(data)
	default:
		opt, err = ParseOptionGeneric(code, data)
	}
	if err != nil {
		return nil, err
	}
	return opt, nil
}

// Options is a collection of options.
type Options []Option

// Get will attempt to get all options that match a DHCPv4 option from its
// OptionCode.  If the option was not found it will return an empty list.
//
// According to RFC 3396, options that are specified more than once are
// concatenated, and hence this should always just return one option.
func (o Options) Get(code OptionCode) []Option {
	opts := []Option{}
	for _, opt := range o {
		if opt.Code() == code {
			opts = append(opts, opt)
		}
	}
	return opts
}

// GetOne will attempt to get an  option that match a Option code.  If there
// are multiple options with the same OptionCode it will only return the first
// one found.  If no matching option is found nil will be returned.
func (o Options) GetOne(code OptionCode) Option {
	for _, opt := range o {
		if opt.Code() == code {
			return opt
		}
	}
	return nil
}

// Has checks whether o has the given `opcode` Option.
func (o Options) Has(code OptionCode) bool {
	return o.GetOne(code) != nil
}

// OptionsFromBytes parses a sequence of bytes until the end and builds a list
// of options from it.
//
// The sequence should not contain the DHCP magic cookie.
//
// Returns an error if any invalid option or length is found.
func OptionsFromBytes(data []byte) (Options, error) {
	return OptionsFromBytesWithParser(data, ParseOption, true)
}

// OptionParser is a function signature for option parsing
type OptionParser func(code OptionCode, data []byte) (Option, error)

// OptionsFromBytesWithParser parses Options from byte sequences using the
// parsing function that is passed in as a paremeter
func OptionsFromBytesWithParser(data []byte, parser OptionParser, checkEndOption bool) (Options, error) {
	if len(data) == 0 {
		return nil, nil
	}
	buf := uio.NewBigEndianBuffer(data)
	options := make(map[OptionCode][]byte, 10)
	var order []OptionCode

	// Due to RFC 3396 allowing an option to be specified multiple times,
	// we have to collect all option data first, and then parse it.
	var end bool
	for buf.Len() >= 1 {
		// 1 byte: option code
		// 1 byte: option length n
		// n bytes: data
		code := OptionCode(buf.Read8())

		if code == OptionPad {
			continue
		} else if code == OptionEnd {
			end = true
			break
		}
		length := int(buf.Read8())

		// N bytes: option data
		data := buf.Consume(length)
		if data == nil {
			return nil, fmt.Errorf("error collecting options: %v", buf.Error())
		}
		data = data[:length:length]

		if _, ok := options[code]; !ok {
			order = append(order, code)
		}
		// RFC 3396: Just concatenate the data if the option code was
		// specified multiple times.
		options[code] = append(options[code], data...)
	}

	// If we never read the End option, the sender of this packet screwed
	// up.
	if !end && checkEndOption {
		return nil, io.ErrUnexpectedEOF
	}

	// Any bytes left must be padding.
	for buf.Len() >= 1 {
		if OptionCode(buf.Read8()) != OptionPad {
			return nil, ErrInvalidOptions
		}
	}

	opts := make(Options, 0, 10)
	for _, code := range order {
		parsedOpt, err := parser(code, options[code])
		if err != nil {
			return nil, fmt.Errorf("error parsing option code %s: %v", code, err)
		}
		opts = append(opts, parsedOpt)
	}
	return opts, nil
}

// Marshal writes options binary representations to b.
func (o Options) Marshal(b *uio.Lexer) {
	for _, opt := range o {
		code := opt.Code()

		// Even if the End option is in there, don't marshal it until
		// the end.
		if code == OptionEnd {
			continue
		} else if code == OptionPad {
			// Some DHCPv4 options have fixed length and do not put
			// length on the wire.
			b.Write8(uint8(code))
			continue
		}

		data := opt.ToBytes()

		// RFC 3396: If more than 256 bytes of data are given, the
		// option is simply listed multiple times.
		for len(data) > 0 {
			// 1 byte: option code
			b.Write8(uint8(code))

			n := len(data)
			if n > math.MaxUint8 {
				n = math.MaxUint8
			}

			// 1 byte: option length
			b.Write8(uint8(n))

			// N bytes: option data
			b.WriteBytes(data[:n])
			data = data[n:]
		}
	}
}
