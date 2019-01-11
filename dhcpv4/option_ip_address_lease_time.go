package dhcpv4

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptIPAddressLeaseTime implements the IP address lease time option described
// by RFC 2132, Section 9.2.
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
