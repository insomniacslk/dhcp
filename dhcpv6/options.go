package dhcpv6

import (
	"errors"
	"fmt"
	"strings"

	"github.com/u-root/uio/uio"
)

// Optioner is an interface that all DHCPv6 options adhere to.
type Optioner interface {
	ToBytes() []byte
	//FromBytes([]byte) error
	String() string
}

// Option is an interface that all DHCPv6 options adhere to.
type Option interface {
	Code() OptionCode
	FromBytes([]byte) error
	Optioner
}

type OptionGeneric struct {
	OptionCode OptionCode
	OptionData []byte
}

func (og OptionGeneric) Code() OptionCode {
	return og.OptionCode
}

func (og OptionGeneric) ToBytes() []byte {
	return og.OptionData
}

func (og OptionGeneric) String() string {
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
//
// Parse a sequence of bytes as a single DHCPv6 option.
// Returns the option structure, or an error if any.
func ParseOption(code OptionCode, optData []byte) (Option, error) {
	var opt Option
	switch code {
	case OptionClientID:
		opt = &optClientID{}
	case OptionServerID:
		opt = &optServerID{}
	case OptionIANA:
		opt = &OptIANA{}
	case OptionIATA:
		opt = &OptIATA{}
	case OptionIAAddr:
		opt = &OptIAAddress{}
	case OptionORO:
		opt = &optRequestedOption{}
	case OptionElapsedTime:
		opt = &optElapsedTime{}
	case OptionRelayMsg:
		opt = &optRelayMsg{}
	case OptionStatusCode:
		opt = &OptStatusCode{}
	case OptionUserClass:
		opt = &OptUserClass{}
	case OptionVendorClass:
		opt = &OptVendorClass{}
	case OptionVendorOpts:
		opt = &OptVendorOpts{}
	case OptionInterfaceID:
		opt = &optInterfaceID{}
	case OptionDNSRecursiveNameServer:
		opt = &optDNS{}
	case OptionDomainSearchList:
		opt = &optDomainSearchList{}
	case OptionIAPD:
		opt = &OptIAPD{}
	case OptionIAPrefix:
		opt = &OptIAPrefix{}
	case OptionInformationRefreshTime:
		opt = &optInformationRefreshTime{}
	case OptionRemoteID:
		opt = &OptRemoteID{}
	case OptionFQDN:
		opt = &OptFQDN{}
	case OptionNTPServer:
		opt = &OptNTPServer{}
	case OptionBootfileURL:
		opt = &optBootFileURL{}
	case OptionBootfileParam:
		opt = &optBootFileParam{}
	case OptionClientArchType:
		opt = &optClientArchType{}
	case OptionNII:
		opt = &OptNetworkInterfaceID{}
	case OptionClientLinkLayerAddr:
		opt = &optClientLinkLayerAddress{}
	case OptionDHCPv4Msg:
		opt = &OptDHCPv4Msg{}
	case OptionDHCP4oDHCP6Server:
		opt = &OptDHCP4oDHCP6Server{}
	case Option4RD:
		opt = &Opt4RD{}
	case Option4RDMapRule:
		opt = &Opt4RDMapRule{}
	case Option4RDNonMapRule:
		opt = &Opt4RDNonMapRule{}
	case OptionRelayPort:
		opt = &optRelayPort{}
	default:
		opt = &OptionGeneric{OptionCode: code}
	}
	return opt, opt.FromBytes(optData)
}

type longStringer interface {
	LongString(spaceIndent int) string
}

// Options is a collection of options.
type Options map[OptionCode][][]byte

func OptionsFrom(list ...Option) Options {
	o := make(Options)
	for _, opt := range list {
		o.Add(opt)
	}
	return o
}

// LongString prints options with indentation of at least spaceIndent spaces.
func (o Options) LongString(spaceIndent int) string {
	indent := strings.Repeat(" ", spaceIndent)
	var s strings.Builder
	if len(o) == 0 {
		s.WriteString("[]")
	} else {
		s.WriteString("[\n")
		/* TODO
		* for _, opt := range o {
			s.WriteString(indent)
			s.WriteString("  ")
			if ls, ok := opt.(longStringer); ok {
				s.WriteString(ls.LongString(spaceIndent + 2))
			} else {
				s.WriteString(opt.String())
			}
			s.WriteString("\n")
		}*/
		s.WriteString(indent)
		s.WriteString("]")
	}
	return s.String()
}

func (o Options) Get(code OptionCode) []Option {
	opts, err := GetOptioner[OptionGeneric, *OptionGeneric](code, o)
	if err != nil {
		return nil
	}
	var os []Option
	for _, opt := range opts {
		os = append(os, &opt)
	}
	return os
}

// Get returns all options matching the option code.
func (o Options) GetRaw(code OptionCode) [][]byte {
	return o[code]
}

var ErrOptionNotFound = errors.New("option not found")

// Da musste erstmal drauf kommen.
// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#pointer-method-example
type Decoder[O any] interface {
	FromBytes([]byte) error
	*O
}

func GetOptioner[T any, PT Decoder[T]](code OptionCode, o Options) ([]T, error) {
	data, ok := o[code]
	if !ok {
		return nil, ErrOptionNotFound
	}

	var ret []T
	for i, p := range data {
		var opt T
		if err := PT(&opt).FromBytes(p); err != nil {
			return nil, fmt.Errorf("option #%d could not be parsed: %w", i+1, err)
		}
		ret = append(ret, opt)
	}
	return ret, nil
}

func MustGetOptioner[T any, PT Decoder[T]](code OptionCode, o Options) []T {
	vals, err := GetOptioner[T, PT](code, o)
	if err != nil {
		return nil
	}
	return vals
}

func GetPtrOptioner[T any, PT Decoder[T]](code OptionCode, o Options) ([]*T, error) {
	data, ok := o[code]
	if !ok {
		return nil, ErrOptionNotFound
	}

	var ret []*T
	for i, p := range data {
		var opt T
		if err := PT(&opt).FromBytes(p); err != nil {
			return nil, fmt.Errorf("option #%d could not be parsed: %w", i+1, err)
		}
		ret = append(ret, &opt)
	}
	return ret, nil
}

func MustGetPtrOptioner[T any, PT Decoder[T]](code OptionCode, o Options) []*T {
	vals, err := GetPtrOptioner[T, PT](code, o)
	if err != nil {
		return nil
	}
	return vals
}

// GetOne returns the first option matching the option code.
func (o Options) GetOneRaw(code OptionCode) []byte {
	data, ok := o[code]
	if !ok || len(data) == 0 {
		return nil
	}
	return data[0]
}

func (o Options) GetOne(code OptionCode) Option {
	opt, err := GetOneOptioner[OptionGeneric, *OptionGeneric](code, o)
	if err != nil {
		return nil
	}
	return &opt
}

func MustGetOneOptioner[T any, PT Decoder[T]](code OptionCode, o Options) T {
	var zerovalue T
	t, err := GetOneOptioner[T, PT](code, o)
	if err != nil {
		return zerovalue
	}
	return t
}

func GetOneOptioner[T any, PT Decoder[T]](code OptionCode, o Options) (T, error) {
	var opt T
	data, ok := o[code]
	if !ok || len(data) == 0 {
		return opt, ErrOptionNotFound
	}
	if err := PT(&opt).FromBytes(data[0]); err != nil {
		return opt, err
	}
	return opt, nil
}

func MustGetOnePtrOptioner[T any, PT Decoder[T]](code OptionCode, o Options) *T {
	t, err := GetOnePtrOptioner[T, PT](code, o)
	if err != nil {
		return nil
	}
	return t
}

func GetOnePtrOptioner[T any, PT Decoder[T]](code OptionCode, o Options) (*T, error) {
	var opt T
	data, ok := o[code]
	if !ok || len(data) == 0 {
		return nil, ErrOptionNotFound
	}
	if err := PT(&opt).FromBytes(data[0]); err != nil {
		return nil, err
	}
	return &opt, nil
}

type DecoderFunc[T interface{}] func([]byte) (T, error)

func GetOneInfOptioner[T interface{}](code OptionCode, o Options, fromBytes DecoderFunc[T]) (T, error) {
	var opt T
	data, ok := o[code]
	if !ok || len(data) == 0 {
		return opt, ErrOptionNotFound
	}
	opt, err := fromBytes(data[0])
	if err != nil {
		return opt, err
	}
	return opt, nil
}

func MustGetOneInfOptioner[T interface{}](code OptionCode, o Options, fromBytes DecoderFunc[T]) T {
	var zerovalue T
	t, err := GetOneInfOptioner[T](code, o, fromBytes)
	if err != nil {
		return zerovalue
	}
	return t
}

// AddRaw appends one option.
func (o Options) AddRaw(code OptionCode, p []byte) {
	o[code] = append(o[code], p)
}

// Add appends one option.
func (o *Options) Add(option Option) {
	if *o == nil {
		*o = make(map[OptionCode][][]byte)
	}
	(*o)[option.Code()] = append((*o)[option.Code()], option.ToBytes())
}

// Del deletes all options matching the option code.
func (o Options) Del(code OptionCode) {
	delete(o, code)
}

// Update replaces the first option of the same type as the specified one.
func (o *Options) Update(option Option) {
	data, ok := (*o)[option.Code()]
	if !ok || len(data) == 0 {
		o.Add(option)
	}
	(*o)[option.Code()][0] = option.ToBytes()
}

// ToBytes marshals all options to bytes.
func (o Options) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for code, opts := range o {
		for _, opt := range opts {
			buf.Write16(uint16(code))
			buf.Write16(uint16(len(opt)))
			buf.WriteBytes(opt)
		}
	}
	return buf.Data()
}

// FromBytes reads option data into o. Options are not deserialized, but the
// overall option structure (type, length, value) has to match or this function
// will return an error.
func (o *Options) FromBytes(data []byte) error {
	*o = make(map[OptionCode][][]byte)
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
		if optData == nil {
			// Buffer did not have `length` bytes left. Malformed
			// packet.
			return fmt.Errorf("error collecting options: %v", buf.Error())
		}

		// TODO: make copy?
		(*o)[code] = append((*o)[code], optData)
	}
	return buf.FinError()
}
