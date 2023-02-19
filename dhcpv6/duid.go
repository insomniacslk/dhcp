package dhcpv6

import (
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/u-root/uio/uio"
)

// DUID is the interface that all DUIDs adhere to.
type DUID interface {
	fmt.Stringer

	ToBytes() []byte
	FromBytes(p []byte) error
	DUIDType() DUIDType
}

// DUIDLLT is a DUID based on link-layer address plus time (RFC 8415 Section 11.2).
type DUIDLLT struct {
	HWType        iana.HWType
	Time          uint32
	LinkLayerAddr net.HardwareAddr
}

// String pretty-prints DUIDLLT information.
func (d DUIDLLT) String() string {
	return fmt.Sprintf("DUID-LLT{HWType=%s HWAddr=%s Time=%d}", d.HWType, d.LinkLayerAddr, d.Time)
}

// DUIDType returns the DUID_LLT type.
func (d DUIDLLT) DUIDType() DUIDType {
	return DUID_LLT
}

// ToBytes serializes the option out to bytes.
func (d DUIDLLT) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(d.DUIDType()))
	buf.Write16(uint16(d.HWType))
	buf.Write32(d.Time)
	buf.WriteBytes(d.LinkLayerAddr)
	return buf.Data()
}

// FromBytes reads the option.
func (d *DUIDLLT) FromBytes(p []byte) error {
	buf := uio.NewBigEndianBuffer(p)
	d.HWType = iana.HWType(buf.Read16())
	d.Time = buf.Read32()
	d.LinkLayerAddr = buf.ReadAll()
	return buf.FinError()
}

// DUIDLL is a DUID based on link-layer (RFC 8415 Section 11.4).
type DUIDLL struct {
	HWType        iana.HWType
	LinkLayerAddr net.HardwareAddr
}

// String pretty-prints DUIDLL information.
func (d DUIDLL) String() string {
	return fmt.Sprintf("DUID-LL{HWType=%s HWAddr=%s}", d.HWType, d.LinkLayerAddr)
}

// DUIDType returns the DUID_LL type.
func (d DUIDLL) DUIDType() DUIDType {
	return DUID_LL
}

// ToBytes serializes the option out to bytes.
func (d DUIDLL) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(d.DUIDType()))
	buf.Write16(uint16(d.HWType))
	buf.WriteBytes(d.LinkLayerAddr)
	return buf.Data()
}

// FromBytes reads the option.
func (d *DUIDLL) FromBytes(p []byte) error {
	buf := uio.NewBigEndianBuffer(p)
	d.HWType = iana.HWType(buf.Read16())
	d.LinkLayerAddr = buf.ReadAll()
	return buf.FinError()
}

// DUIDEN is a DUID based on enterprise number (RFC 8415 Section 11.3).
type DUIDEN struct {
	EnterpriseNumber     uint32
	EnterpriseIdentifier []byte
}

// String pretty-prints DUIDEN information.
func (d DUIDEN) String() string {
	return fmt.Sprintf("DUID-EN{EnterpriseNumber=%d EnterpriseIdentifier=%s}", d.EnterpriseNumber, d.EnterpriseIdentifier)
}

// DUIDType returns the DUID_EN type.
func (d DUIDEN) DUIDType() DUIDType {
	return DUID_EN
}

// ToBytes serializes the option out to bytes.
func (d DUIDEN) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(d.DUIDType()))
	buf.Write32(d.EnterpriseNumber)
	buf.WriteBytes(d.EnterpriseIdentifier)
	return buf.Data()
}

// FromBytes reads the option.
func (d *DUIDEN) FromBytes(p []byte) error {
	buf := uio.NewBigEndianBuffer(p)
	d.EnterpriseNumber = buf.Read32()
	d.EnterpriseIdentifier = buf.ReadAll()
	return buf.FinError()
}

// DUIDUUID is a DUID based on UUID (RFC 8415 Section 11.5).
type DUIDUUID struct {
	// Defined by RFC 6355.
	UUID [16]byte
}

// String pretty-prints DUIDUUID information.
func (d DUIDUUID) String() string {
	return fmt.Sprintf("DUID-UUID{%#x}", d.UUID[:])
}

// DUIDType returns the DUID_UUID type.
func (d DUIDUUID) DUIDType() DUIDType {
	return DUID_UUID
}

// ToBytes serializes the option out to bytes.
func (d DUIDUUID) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(d.DUIDType()))
	buf.WriteData(d.UUID[:])
	return buf.Data()
}

// FromBytes reads the option.
func (d *DUIDUUID) FromBytes(p []byte) error {
	if len(p) != 16 {
		return fmt.Errorf("buffer is length %d, DUID-UUID must be exactly 16 bytes", len(p))
	}
	copy(d.UUID[:], p)
	return nil
}

// Equal returns true if e is a DUID-UUID with the same values as d.
func (d *DUIDUUID) Equal(e DUID) bool {
	euuid, ok := e.(*DUIDUUID)
	if !ok {
		return false
	}
	return d.UUID == euuid.UUID
}

// DUIDOpaque is a DUID of unknown type.
type DUIDOpaque struct {
	Type DUIDType
	Data []byte
}

// String pretty-prints opaque DUID information.
func (d DUIDOpaque) String() string {
	return fmt.Sprintf("DUID-Opaque{Type=%d Data=%#x}", d.Type, d.Data)
}

// DUIDType returns the opaque DUID type.
func (d DUIDOpaque) DUIDType() DUIDType {
	return d.Type
}

// ToBytes serializes the option out to bytes.
func (d DUIDOpaque) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(d.Type))
	buf.WriteData(d.Data)
	return buf.Data()
}

// FromBytes reads the option.
func (d *DUIDOpaque) FromBytes(p []byte) error {
	d.Data = append([]byte(nil), p...)
	return nil
}

// DUIDType is the DUID type as defined in RFC 3315.
type DUIDType uint16

// DUID types
const (
	DUID_LLT  DUIDType = 1
	DUID_EN   DUIDType = 2
	DUID_LL   DUIDType = 3
	DUID_UUID DUIDType = 4
)

// duidTypeToString maps a DUIDType to a name.
var duidTypeToString = map[DUIDType]string{
	DUID_LL:   "DUID-LL",
	DUID_LLT:  "DUID-LLT",
	DUID_EN:   "DUID-EN",
	DUID_UUID: "DUID-UUID",
}

func (d DUIDType) String() string {
	if dtype, ok := duidTypeToString[d]; ok {
		return dtype
	}
	return "unknown"
}

// DUIDFromBytes parses a DUID from a byte slice.
func DUIDFromBytes(data []byte) (DUID, error) {
	buf := uio.NewBigEndianBuffer(data)
	if !buf.Has(2) {
		return nil, fmt.Errorf("buffer too short: have %d bytes, want 2 bytes", buf.Len())
	}

	typ := DUIDType(buf.Read16())
	var d DUID
	switch typ {
	case DUID_LLT:
		d = &DUIDLLT{}
	case DUID_LL:
		d = &DUIDLL{}
	case DUID_EN:
		d = &DUIDEN{}
	case DUID_UUID:
		d = &DUIDUUID{}
	default:
		d = &DUIDOpaque{Type: typ}
	}
	return d, d.FromBytes(buf.Data())
}
