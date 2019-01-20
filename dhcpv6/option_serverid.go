package dhcpv6

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptServerId represents a Server ID option
//
// This module defines the OptServerId and DUID structures.
// https://www.ietf.org/rfc/rfc3315.txt
type OptServerId struct {
	Sid Duid
}

func (op *OptServerId) Code() OptionCode {
	return OptionServerID
}

// ToBytes serializes this option.
func (op *OptServerId) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(OptionServerID))
	buf.Write16(uint16(op.Length()))
	buf.WriteBytes(op.Sid.ToBytes())
	return buf.Data()
}

func (op *OptServerId) Length() int {
	return op.Sid.Length()
}

func (op *OptServerId) String() string {
	return fmt.Sprintf("OptServerId{sid=%v}", op.Sid.String())
}

// ParseOptServerId builds an OptServerId structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptServerId(data []byte) (*OptServerId, error) {
	var opt OptServerId
	sid, err := DuidFromBytes(data)
	if err != nil {
		return nil, err
	}
	opt.Sid = *sid
	return &opt, nil
}
