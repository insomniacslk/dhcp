package dhcpv4

import (
	"bytes"
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// This option implements the Vendor-Identifying Vendor Class Option
// https://tools.ietf.org/html/rfc3925

// VIVCIdentifier represents one Vendor-Identifying vendor class option.
type VIVCIdentifier struct {
	EntID uint32
	Data  []byte
}

// OptVIVC represents the DHCP message type option.
type OptVIVC struct {
	Identifiers []VIVCIdentifier
}

// ParseOptVIVC contructs an OptVIVC tsruct from a sequence of bytes and returns
// it, or an error.
func ParseOptVIVC(data []byte) (*OptVIVC, error) {
	buf := uio.NewBigEndianBuffer(data)

	var ids []VIVCIdentifier
	for buf.Has(5) {
		entID := buf.Read32()
		idLen := int(buf.Read8())
		ids = append(ids, VIVCIdentifier{EntID: entID, Data: buf.CopyN(idLen)})
	}

	return &OptVIVC{Identifiers: ids}, buf.FinError()
}

// Code returns the option code.
func (o *OptVIVC) Code() OptionCode {
	return OptionVendorIdentifyingVendorClass
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptVIVC) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, id := range o.Identifiers {
		buf.Write32(id.EntID)
		buf.Write8(uint8(len(id.Data)))
		buf.WriteBytes(id.Data)
	}
	return buf.Data()
}

// String returns a human-readable string for this option.
func (o *OptVIVC) String() string {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "Vendor-Identifying Vendor Class ->")

	for _, id := range o.Identifiers {
		fmt.Fprintf(&buf, " %d:'%s',", id.EntID, id.Data)
	}

	return buf.String()[:buf.Len()-1]
}
