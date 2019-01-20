package dhcpv6

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptInterfaceId implements the interface-id option as defined by RFC 3315,
// Section 22.18.
//
// This module defines the OptInterfaceId structure.
// https://www.ietf.org/rfc/rfc3315.txt
type OptInterfaceId struct {
	interfaceId []byte
}

func (op *OptInterfaceId) Code() OptionCode {
	return OptionInterfaceID
}

func (op *OptInterfaceId) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(OptionInterfaceID))
	buf.Write16(uint16(len(op.interfaceId)))
	buf.WriteBytes(op.interfaceId)
	return buf.Data()
}

func (op *OptInterfaceId) InterfaceID() []byte {
	return op.interfaceId
}

func (op *OptInterfaceId) SetInterfaceID(interfaceId []byte) {
	op.interfaceId = append([]byte(nil), interfaceId...)
}

func (op *OptInterfaceId) Length() int {
	return len(op.interfaceId)
}

func (op *OptInterfaceId) String() string {
	return fmt.Sprintf("OptInterfaceId{interfaceid=%v}", op.interfaceId)
}

// build an OptInterfaceId structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptInterfaceId(data []byte) (*OptInterfaceId, error) {
	var opt OptInterfaceId
	opt.interfaceId = append([]byte(nil), data...)
	return &opt, nil
}
