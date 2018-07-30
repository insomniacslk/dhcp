package dhcpv4

import (
	"encoding/binary"
	"fmt"
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
	// Should at least have code, length, and lease time.
	if len(data) < 6 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionIPAddressLeaseTime {
		return nil, fmt.Errorf("expected option %v, got %v instead", OptionIPAddressLeaseTime, code)
	}
	length := int(data[1])
	if length != 4 {
		return nil, fmt.Errorf("expected length 4, got %v instead", length)
	}
	leaseTime := binary.BigEndian.Uint32(data[2:6])
	return &OptIPAddressLeaseTime{LeaseTime: leaseTime}, nil
}

// Code returns the option code.
func (o *OptIPAddressLeaseTime) Code() OptionCode {
	return OptionIPAddressLeaseTime
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptIPAddressLeaseTime) ToBytes() []byte {
	serializedTime := make([]byte, 4)
	binary.BigEndian.PutUint32(serializedTime, o.LeaseTime)
	serializedOpt := []byte{byte(o.Code()), byte(o.Length())}
	return append(serializedOpt, serializedTime...)
}

// String returns a human-readable string for this option.
func (o *OptIPAddressLeaseTime) String() string {
	return fmt.Sprintf("IP Addresses Lease Time -> %v", o.LeaseTime)
}

// Length returns the length of the data portion (excluding option code and byte
// for length, if any).
func (o *OptIPAddressLeaseTime) Length() int {
	return 4
}
