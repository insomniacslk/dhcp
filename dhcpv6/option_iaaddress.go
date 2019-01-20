package dhcpv6

import (
	"fmt"
	"net"

	"github.com/u-root/u-root/pkg/uio"
)

// OptIAAddress represents an OptionIAAddr.
//
// This module defines the OptIAAddress structure.
// https://www.ietf.org/rfc/rfc3633.txt
type OptIAAddress struct {
	IPv6Addr          net.IP
	PreferredLifetime uint32
	ValidLifetime     uint32
	Options           Options
}

// Code returns the option's code
func (op *OptIAAddress) Code() OptionCode {
	return OptionIAAddr
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptIAAddress) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(make([]byte, 0, 28))
	buf.Write16(uint16(OptionIAAddr))
	buf.Write16(uint16(op.Length()))
	buf.WriteBytes(op.IPv6Addr.To16())
	buf.Write32(op.PreferredLifetime)
	buf.Write32(op.ValidLifetime)
	buf.WriteBytes(op.Options.ToBytes())
	return buf.Data()
}

// Length returns the option length
func (op *OptIAAddress) Length() int {
	opLen := 24
	for _, opt := range op.Options {
		opLen += 4 + opt.Length()
	}
	return opLen
}

func (op *OptIAAddress) String() string {
	return fmt.Sprintf("OptIAAddress{ipv6addr=%v, preferredlifetime=%v, validlifetime=%v, options=%v}",
		op.IPv6Addr, op.PreferredLifetime, op.ValidLifetime, op.Options)
}

// ParseOptIAAddress builds an OptIAAddress structure from a sequence
// of bytes. The input data does not include option code and length
// bytes.
func ParseOptIAAddress(data []byte) (*OptIAAddress, error) {
	var opt OptIAAddress
	buf := uio.NewBigEndianBuffer(data)
	opt.IPv6Addr = net.IP(buf.CopyN(net.IPv6len))
	opt.PreferredLifetime = buf.Read32()
	opt.ValidLifetime = buf.Read32()
	if err := opt.Options.FromBytes(buf.ReadAll()); err != nil {
		return nil, err
	}
	return &opt, buf.FinError()
}
