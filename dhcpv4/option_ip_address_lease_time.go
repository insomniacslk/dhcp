package dhcpv4

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// This option implements the IP Address Lease Time option
// https://tools.ietf.org/html/rfc2132

// OptIPAddressLeaseTime represents the IP Address Lease Time option.
type OptIPAddressLeaseTime struct {
	LeaseTime uint32
}

// ParseOptIPAddressLeaseTime constructs an OptIPAddressLeaseTime struct from a
// sequence of bytes and returns it, or an error.
func ParseOptIPAddressLeaseTime(data []byte) (*OptIPAddressLeaseTime, error) {
	buf := uio.NewBigEndianBuffer(data)
	leaseTime := buf.Read32()
	return &OptIPAddressLeaseTime{LeaseTime: leaseTime}, buf.FinError()
}

// Code returns the option code.
func (o *OptIPAddressLeaseTime) Code() OptionCode {
	return OptionIPAddressLeaseTime
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptIPAddressLeaseTime) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write32(o.LeaseTime)
	return buf.Data()
}

// String returns a human-readable string for this option.
func (o *OptIPAddressLeaseTime) String() string {
	return fmt.Sprintf("IP Addresses Lease Time -> %v", o.LeaseTime)
}
