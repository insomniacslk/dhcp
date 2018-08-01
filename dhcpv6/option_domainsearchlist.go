package dhcpv6

// This module defines the OptDomainSearchList structure.
// https://www.ietf.org/rfc/rfc3646.txt

import (
	"encoding/binary"
	"fmt"

	"github.com/insomniacslk/dhcp/rfc1035label"
)

// OptDomainSearchList list implements a DOMAIN_SEARCH_LIST option
type OptDomainSearchList struct {
	DomainSearchList []string
}

func (op *OptDomainSearchList) Code() OptionCode {
	return DOMAIN_SEARCH_LIST
}

func (op *OptDomainSearchList) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(DOMAIN_SEARCH_LIST))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	buf = append(buf, rfc1035label.LabelsToBytes(op.DomainSearchList)...)
	return buf
}

func (op *OptDomainSearchList) Length() int {
	var length int
	for _, label := range op.DomainSearchList {
		length += len(label) + 2 // add the first and the last length bytes
	}
	return length
}

func (op *OptDomainSearchList) String() string {
	return fmt.Sprintf("OptDomainSearchList{searchlist=%v}", op.DomainSearchList)
}

// build an OptDomainSearchList structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptDomainSearchList(data []byte) (*OptDomainSearchList, error) {
	opt := OptDomainSearchList{}
	var err error
	opt.DomainSearchList, err = rfc1035label.LabelsFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
