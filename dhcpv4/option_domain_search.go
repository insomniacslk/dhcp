package dhcpv4

// This module defines the OptDomainSearch structure.
// https://tools.ietf.org/html/rfc3397

import (
	"fmt"

	"github.com/insomniacslk/dhcp/rfc1035label"
)

// OptDomainSearch represents an option encapsulating a domain search list.
type OptDomainSearch struct {
	DomainSearch []string
}

// Code returns the option code.
func (op *OptDomainSearch) Code() OptionCode {
	return OptionDNSDomainSearchList
}

// ToBytes returns a serialized stream of bytes for this option.
func (op *OptDomainSearch) ToBytes() []byte {
	buf := []byte{byte(op.Code()), byte(op.Length())}
	buf = append(buf, rfc1035label.LabelsToBytes(op.DomainSearch)...)
	return buf
}

// Length returns the length of the data portion (excluding option code an byte
// length). 
func (op *OptDomainSearch) Length() int {
	var length int
	for _, label := range op.DomainSearch {
		length += len(label) + 2 // add the first and the last length bytes
	}
	return length
}

// String returns a human-readable string.
func (op *OptDomainSearch) String() string {
	return fmt.Sprintf("DNS Domain Search List -> %v", op.DomainSearch)
}

// ParseOptDomainSearch returns a new OptDomainSearch from a byte stream, or
// error if any.
func ParseOptDomainSearch(data []byte) (*OptDomainSearch, error) {
	if len(data) < 2 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionDNSDomainSearchList {
		return nil, fmt.Errorf("expected code %v, got %v", OptionDNSDomainSearchList, code)
	}
	length := int(data[1])
	if len(data) < 2+length {
		return nil, ErrShortByteStream
	}
	domainSearch, err := rfc1035label.LabelsFromBytes(data[2:length+2])
	if err != nil {
		return nil, err
	}
	return &OptDomainSearch{DomainSearch: domainSearch}, nil
}
