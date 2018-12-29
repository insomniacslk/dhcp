package dhcpv4

import (
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/uio"
)

// This option implements the parameter request list option
// https://tools.ietf.org/html/rfc2132

// OptParameterRequestList represents the parameter request list option.
type OptParameterRequestList struct {
	RequestedOpts []OptionCode
}

// ParseOptParameterRequestList returns a new OptParameterRequestList from a
// byte stream, or error if any.
func ParseOptParameterRequestList(data []byte) (*OptParameterRequestList, error) {
	buf := uio.NewBigEndianBuffer(data)
	requestedOpts := make([]OptionCode, 0, buf.Len())
	for buf.Len() > 0 {
		requestedOpts = append(requestedOpts, OptionCode(buf.Read8()))
	}
	return &OptParameterRequestList{RequestedOpts: requestedOpts}, buf.Error()
}

// Code returns the option code.
func (o *OptParameterRequestList) Code() OptionCode {
	return OptionParameterRequestList
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptParameterRequestList) ToBytes() []byte {
	ret := []byte{byte(o.Code()), byte(o.Length())}
	for _, req := range o.RequestedOpts {
		ret = append(ret, byte(req))
	}
	return ret
}

// String returns a human-readable string for this option.
func (o *OptParameterRequestList) String() string {
	var optNames []string
	for _, ro := range o.RequestedOpts {
		name := ro.String()
		if name == "Unknown" {
			name += fmt.Sprintf("%s (%v)", name, ro)
		}
		optNames = append(optNames, name)
	}
	return fmt.Sprintf("Parameter Request List -> [%v]", strings.Join(optNames, ", "))
}

// Length returns the length of the data portion (excluding option code and byte
// for length, if any).
func (o *OptParameterRequestList) Length() int {
	return len(o.RequestedOpts)
}
