package dhcpv6

// This module defines the OptInterfaceId structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"encoding/binary"
	"fmt"
)

type OptInterfaceId struct {
	interfaceId []byte
}

func (op *OptInterfaceId) Code() OptionCode {
	return OPTION_INTERFACE_ID
}

func (op *OptInterfaceId) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_INTERFACE_ID))
	binary.BigEndian.PutUint16(buf[2:4], uint16(len(op.interfaceId)))
	buf = append(buf, op.interfaceId...)
	return buf
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
	opt := OptInterfaceId{}
	opt.interfaceId = append([]byte(nil), data...)
	return &opt, nil
}
