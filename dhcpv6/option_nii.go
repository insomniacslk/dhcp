package dhcpv6

// This module defines the OptNetworkInterfaceId structure.
// https://www.ietf.org/rfc/rfc5970.txt

import (
	"encoding/binary"
	"fmt"
)

// see rfc4578
const (
	NII_LANDESK_NOPXE = iota
	NII_PXE_GEN_I
	NII_PXE_GEN_II
	NII_UNDI_NOEFI
	NII_UNDI_EFI_GEN_I
	NII_UNDI_EFI_GEN_II
)

var NIIToStringMap = map[uint8]string{
	NII_LANDESK_NOPXE:   "LANDesk service agent boot ROMs. No PXE",
	NII_PXE_GEN_I:       "First gen. PXE boot ROMs",
	NII_PXE_GEN_II:      "Second gen. PXE boot ROMs",
	NII_UNDI_NOEFI:      "UNDI 32/64 bit. UEFI drivers, no UEFI runtime",
	NII_UNDI_EFI_GEN_I:  "UNDI 32/64 bit. UEFI runtime 1st gen",
	NII_UNDI_EFI_GEN_II: "UNDI 32/64 bit. UEFI runtime 2nd gen",
}

type OptNetworkInterfaceId struct {
	type_        uint8
	major, minor uint8 // revision number
}

func (op *OptNetworkInterfaceId) Code() OptionCode {
	return OPTION_NII
}

func (op *OptNetworkInterfaceId) ToBytes() []byte {
	buf := make([]byte, 7)
	binary.BigEndian.PutUint16(buf[0:2], uint16(OPTION_NII))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	buf[4] = op.type_
	buf[5] = op.major
	buf[6] = op.minor
	return buf
}

func (op *OptNetworkInterfaceId) Type() uint8 {
	return op.type_
}

func (op *OptNetworkInterfaceId) SetType(type_ uint8) {
	op.type_ = type_
}

func (op *OptNetworkInterfaceId) Major() uint8 {
	return op.major
}

func (op *OptNetworkInterfaceId) SetMajor(major uint8) {
	op.major = major
}

func (op *OptNetworkInterfaceId) Minor() uint8 {
	return op.minor
}

func (op *OptNetworkInterfaceId) SetMinor(minor uint8) {
	op.minor = minor
}

func (op *OptNetworkInterfaceId) Length() int {
	return 3
}

func (op *OptNetworkInterfaceId) String() string {
	typeName, ok := NIIToStringMap[op.type_]
	if !ok {
		typeName = "Unknown"
	}
	return fmt.Sprintf("OptNetworkInterfaceId{type=%v, revision=%v.%v}",
		typeName, op.major, op.minor,
	)
}

// build an OptNetworkInterfaceId structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptNetworkInterfaceId(data []byte) (*OptNetworkInterfaceId, error) {
	opt := OptNetworkInterfaceId{}
	if len(data) != 3 {
		return nil, fmt.Errorf("Invalid arch type data length. Expected 3 bytes, got %v", len(data))
	}
	opt.type_ = data[0]
	opt.major = data[1]
	opt.minor = data[2]
	return &opt, nil
}
