package dhcpv6

// This module defines the OptDomainSearchList structure.
// https://www.ietf.org/rfc/rfc3646.txt

import (
	"encoding/binary"
	"fmt"

	"github.com/insomniacslk/dhcp/rfc1035label"
)

// OptDomainSearchList list implements a OptionDomainSearchList option
type OptDomainSearchList struct {
	DomainSearchList *rfc1035label.Labels
}

func (op *OptDomainSearchList) Code() OptionCode {
	return OptionDomainSearchList
}

func (op *OptDomainSearchList) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionDomainSearchList))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	buf = append(buf, op.DomainSearchList.ToBytes()...)
	return buf
}

func (op *OptDomainSearchList) Length() int {
	var length int
	for _, label := range op.DomainSearchList.Labels {
		length += len(label) + 2 // add the first and the last length bytes
	}
	return length
}

func (op *OptDomainSearchList) String() string {
	return fmt.Sprintf("OptDomainSearchList{searchlist=%v}", op.DomainSearchList.Labels)
}

// ParseOptDomainSearchList builds an OptDomainSearchList structure from a sequence
// of bytes. The input data does not include option code and length bytes.
func ParseOptDomainSearchList(data []byte) (*OptDomainSearchList, error) {
	opt := OptDomainSearchList{}
	labels, err := rfc1035label.FromBytes(data)
	if err != nil {
		return nil, err
	}
	opt.DomainSearchList = labels
	return &opt, nil
}
