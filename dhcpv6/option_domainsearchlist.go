package dhcpv6

import (
	"fmt"

	"github.com/insomniacslk/dhcp/rfc1035label"
	"github.com/u-root/u-root/pkg/uio"
)

// OptDomainSearchList list implements a OptionDomainSearchList option
//
// This module defines the OptDomainSearchList structure.
// https://www.ietf.org/rfc/rfc3646.txt
type OptDomainSearchList struct {
	DomainSearchList *rfc1035label.Labels
}

func (op *OptDomainSearchList) Code() OptionCode {
	return OptionDomainSearchList
}

// ToBytes marshals this option to bytes.
func (op *OptDomainSearchList) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(OptionDomainSearchList))
	buf.Write16(uint16(op.Length()))
	buf.WriteBytes(op.DomainSearchList.ToBytes())
	return buf.Data()
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
	var opt OptDomainSearchList
	var err error
	opt.DomainSearchList, err = rfc1035label.FromBytes(data)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
