package dhcpv4

import (
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/uio"
)

// OptionCodeList is a list of DHCP option codes.
type OptionCodeList []OptionCode

// Has returns whether c is in the list.
func (ol OptionCodeList) Has(c OptionCode) bool {
	for _, code := range ol {
		if code == c {
			return true
		}
	}
	return false
}

// Add adds option codes in cs to ol.
func (ol *OptionCodeList) Add(cs ...OptionCode) {
	for _, c := range cs {
		if !ol.Has(c) {
			*ol = append(*ol, c)
		}
	}
}

// String returns a human-readable string for the option names.
func (ol OptionCodeList) String() string {
	var names []string
	for _, code := range ol {
		names = append(names, code.String())
	}
	return strings.Join(names, ", ")
}

// OptParameterRequestList implements the parameter request list option
// described by RFC 2132, Section 9.8.
type OptParameterRequestList struct {
	RequestedOpts OptionCodeList
}

// ParseOptParameterRequestList returns a new OptParameterRequestList from a
// byte stream, or error if any.
func ParseOptParameterRequestList(data []byte) (*OptParameterRequestList, error) {
	buf := uio.NewBigEndianBuffer(data)
	requestedOpts := make(OptionCodeList, 0, buf.Len())
	for buf.Has(1) {
		requestedOpts = append(requestedOpts, optionCode(buf.Read8()))
	}
	return &OptParameterRequestList{RequestedOpts: requestedOpts}, buf.Error()
}

// Code returns the option code.
func (o *OptParameterRequestList) Code() OptionCode {
	return OptionParameterRequestList
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptParameterRequestList) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, req := range o.RequestedOpts {
		buf.Write8(req.Code())
	}
	return buf.Data()
}

// String returns a human-readable string for this option.
func (o *OptParameterRequestList) String() string {
	return fmt.Sprintf("Parameter Request List -> %s", o.RequestedOpts)
}
