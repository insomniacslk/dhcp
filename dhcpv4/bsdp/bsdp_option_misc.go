package bsdp

import (
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// OptReplyPort returns a new BSDP reply port option.
//
// Implements the BSDP option reply port. This is used when BSDP responses
// should be sent to a reply port other than the DHCP default. The macOS GUI
// "Startup Disk Select" sends this option since it's operating in an
// unprivileged context.
func OptReplyPort(port uint16) dhcpv4.Option {
	return dhcpv4.Option{Code: OptionReplyPort, Value: dhcpv4.Uint16(port)}
}

// GetReplyPort returns the BSDP reply port in o, if present.
func GetReplyPort(o dhcpv4.Options) (uint16, error) {
	return dhcpv4.GetUint16(OptionReplyPort, o)
}

// OptServerPriority returns a new BSDP server priority option.
func OptServerPriority(prio uint16) dhcpv4.Option {
	return dhcpv4.Option{Code: OptionServerPriority, Value: dhcpv4.Uint16(prio)}
}

// GetServerPriority returns the BSDP server priority in o if present.
func GetServerPriority(o dhcpv4.Options) (uint16, error) {
	return dhcpv4.GetUint16(OptionServerPriority, o)
}

// OptMachineName returns a BSDP Machine Name option.
func OptMachineName(name string) dhcpv4.Option {
	return dhcpv4.Option{Code: OptionMachineName, Value: dhcpv4.String(name)}
}

// GetMachineName finds and parses the BSDP Machine Name option from o.
func GetMachineName(o dhcpv4.Options) string {
	return dhcpv4.GetString(OptionMachineName, o)
}

// Version is the BSDP protocol version. Can be one of 1.0 or 1.1.
type Version [2]byte

// Specific versions.
var (
	Version1_0 = Version{1, 0}
	Version1_1 = Version{1, 1}
)

// ToBytes returns a serialized stream of bytes for this option.
func (o Version) ToBytes() []byte {
	return o[:]
}

// String returns a human-readable string for this option.
func (o Version) String() string {
	return fmt.Sprintf("%d.%d", o[0], o[1])
}

// FromBytes constructs a Version struct from a sequence of
// bytes and returns it, or an error.
func (o *Version) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	buf.ReadBytes(o[:])
	return buf.FinError()
}

// OptVersion returns a new BSDP version option.
func OptVersion(version Version) dhcpv4.Option {
	return dhcpv4.Option{Code: OptionVersion, Value: version}
}

// GetVersion returns the BSDP version in o if present.
func GetVersion(o dhcpv4.Options) (Version, error) {
	v := o.Get(OptionVersion)
	if v == nil {
		return Version{0, 0}, fmt.Errorf("version not found")
	}
	var ver Version
	if err := ver.FromBytes(v); err != nil {
		return Version{0, 0}, err
	}
	return ver, nil
}

// GetServerIdentifier returns the BSDP Server Identifier value in o.
func GetServerIdentifier(o dhcpv4.Options) net.IP {
	return dhcpv4.GetIP(OptionServerIdentifier, o)
}

// OptServerIdentifier returns a new BSDP Server Identifier option.
func OptServerIdentifier(ip net.IP) dhcpv4.Option {
	return dhcpv4.Option{Code: OptionServerIdentifier, Value: dhcpv4.IP(ip)}
}
