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

// Option is an interface that all DHCP v4 options adhere to.
type Option interface {
	Code() OptionCode
	ToBytes() []byte
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

// Update replaces an existing option with the same option code with the given
// one, adding it if not already present.
//
// Per RFC 2131, Section 4.1, "options may appear only once."
//
// An End option is ignored.
func (o *Options) Update(option Option) {
	if option.Code() == OptionEnd {
		return
	}

	for idx, opt := range *o {
		if opt.Code() == option.Code() {
			(*o)[idx] = option
			// Don't look further.
			return
		}
	}
	// If not found, add it.
	*o = append(*o, option)
}

// OptionsFromBytes parses a sequence of bytes until the end and builds a list
// of options from it.
//
// The sequence should not contain the DHCP magic cookie.
//
// Returns an error if any invalid option or length is found.
func OptionsFromBytes(data []byte) (Options, error) {
	return OptionsFromBytesWithParser(data, codeGetter, ParseOption, true)
}

// OptionParser is a function signature for option parsing.
type OptionParser func(code OptionCode, data []byte) (Option, error)

// OptionCodeGetter parses a code into an OptionCode.
type OptionCodeGetter func(code uint8) OptionCode

// codeGetter is an OptionCodeGetter for DHCP optionCodes.
func codeGetter(c uint8) OptionCode {
	return optionCode(c)
}

// OptionsFromBytesWithParser parses Options from byte sequences using the
// parsing function that is passed in as a paremeter
func OptionsFromBytesWithParser(data []byte, coder OptionCodeGetter, parser OptionParser, checkEndOption bool) (Options, error) {
	if len(data) == 0 {
		return nil, nil
	}
	buf := uio.NewBigEndianBuffer(data)
	options := make(map[OptionCode][]byte, 10)
	var order []OptionCode

	// Due to RFC 2131, 3396 allowing an option to be specified multiple
	// times, we have to collect all option data first, and then parse it.
	var end bool
	for buf.Len() >= 1 {
		// 1 byte: option code
		// 1 byte: option length n
		// n bytes: data
		code := buf.Read8()

		if code == OptionPad.Code() {
			continue
		} else if code == OptionEnd.Code() {
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

		// Get the OptionCode for this guy.
		c := coder(code)
		if _, ok := options[c]; !ok {
			order = append(order, c)
		}

		// RFC 2131, Section 4.1 "Options may appear only once, [...].
		// The client concatenates the values of multiple instances of
		// the same option into a single parameter list for
		// configuration."
		//
		// See also RFC 3396 for concatenation order and options longer
		// than 255 bytes.
		options[c] = append(options[c], data...)
	}

	// If we never read the End option, the sender of this packet screwed
	// up.
	if !end && checkEndOption {
		return nil, io.ErrUnexpectedEOF
	}

	// Any bytes left must be padding.
	for buf.Len() >= 1 {
		if buf.Read8() != OptionPad.Code() {
			return nil, ErrInvalidOptions
		}
	}

	opts := make(Options, 0, len(options))
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
		code := opt.Code().Code()

		// Even if the End option is in there, don't marshal it until
		// the end.
		if code == OptionEnd.Code() {
			continue
		} else if code == OptionPad.Code() {
			// Some DHCPv4 options have fixed length and do not put
			// length on the wire.
			b.Write8(code)
			continue
		}

		data := opt.ToBytes()

		// RFC 3396: If more than 256 bytes of data are given, the
		// option is simply listed multiple times.
		for len(data) > 0 {
			// 1 byte: option code
			b.Write8(code)

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
