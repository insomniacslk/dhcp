package dhcpv6

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptClientId represents a Client ID option
//
// This module defines the OptClientId and DUID structures.
// https://www.ietf.org/rfc/rfc3315.txt
type OptClientId struct {
	Cid Duid
}

func (op *OptClientId) Code() OptionCode {
	return OptionClientID
}

// ToBytes marshals the Client ID option as defined by RFC 3315, Section 22.2.
func (op *OptClientId) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(OptionClientID))
	buf.Write16(uint16(op.Length()))
	buf.WriteBytes(op.Cid.ToBytes())
	return buf.Data()
}

func (op *OptClientId) Length() int {
	return op.Cid.Length()
}

func (op *OptClientId) String() string {
	return fmt.Sprintf("OptClientId{cid=%v}", op.Cid.String())
}

// ParseOptClientId builds an OptClientId structure from a sequence
// of bytes. The input data does not include option code and length
// bytes.
func ParseOptClientId(data []byte) (*OptClientId, error) {
	var opt OptClientId
	cid, err := DuidFromBytes(data)
	if err != nil {
		return nil, err
	}
	opt.Cid = *cid
	return &opt, nil
}
