package dhcpv6

// This module defines the OptRemoteId structure.
// https://www.ietf.org/rfc/rfc4649.txt

import (
	"encoding/binary"
	"fmt"
)

type OptRemoteId struct {
	enterpriseNumber uint32
	remoteId         []byte
}

func (op *OptRemoteId) Code() OptionCode {
	return OptionRemoteID
}

func (op *OptRemoteId) ToBytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OptionRemoteID))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	binary.BigEndian.PutUint32(buf[4:8], uint32(op.enterpriseNumber))
	buf = append(buf, op.remoteId...)
	return buf
}

func (op *OptRemoteId) EnterpriseNumber() uint32 {
	return op.enterpriseNumber
}

func (op *OptRemoteId) SetEnterpriseNumber(enterpriseNumber uint32) {
	op.enterpriseNumber = enterpriseNumber
}

func (op *OptRemoteId) RemoteID() []byte {
	return op.remoteId
}

func (op *OptRemoteId) SetRemoteID(remoteId []byte) {
	op.remoteId = append([]byte(nil), remoteId...)
}

func (op *OptRemoteId) Length() int {
	return 4 + len(op.remoteId)
}

func (op *OptRemoteId) String() string {
	return fmt.Sprintf("OptRemoteId{enterprisenum=%v, remoteid=%v}",
		op.enterpriseNumber, op.remoteId,
	)
}

// build an OptRemoteId structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptRemoteId(data []byte) (*OptRemoteId, error) {
	opt := OptRemoteId{}
	if len(data) < 4 {
		return nil, fmt.Errorf("Invalid remote id data length. Expected at least 4 bytes, got %v", len(data))
	}
	opt.enterpriseNumber = binary.BigEndian.Uint32(data[:4])
	opt.remoteId = append([]byte(nil), data[4:]...)
	return &opt, nil
}
