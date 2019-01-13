package dhcpv4

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/u-root/u-root/pkg/uio"
)

const (
	// minPacketLen is the minimum DHCP header length.
	minPacketLen = 236

	// Maximum length of the ClientHWAddr (client hardware address) according to
	// RFC 2131, Section 2. This is the link-layer destination a server
	// must send responses to.
	maxHWAddrLen = 16

	// MaxMessageSize is the maximum size in bytes that a DHCPv4 packet can hold.
	MaxMessageSize = 576
)

// magicCookie is the magic 4-byte value at the beginning of the list of options
// in a DHCPv4 packet.
var magicCookie = [4]byte{99, 130, 83, 99}

// DHCPv4 represents a DHCPv4 packet header and options. See the New* functions
// to build DHCPv4 packets.
type DHCPv4 struct {
	OpCode         OpcodeType
	HWType         iana.HWType
	HopCount       uint8
	TransactionID  TransactionID
	NumSeconds     uint16
	Flags          uint16
	ClientIPAddr   net.IP
	YourIPAddr     net.IP
	ServerIPAddr   net.IP
	GatewayIPAddr  net.IP
	ClientHWAddr   net.HardwareAddr
	ServerHostName string
	BootFileName   string
	Options        Options
}

// Modifier defines the signature for functions that can modify DHCPv4
// structures. This is used to simplify packet manipulation
type Modifier func(d *DHCPv4) *DHCPv4

// IPv4AddrsForInterface obtains the currently-configured, non-loopback IPv4
// addresses for iface.
func IPv4AddrsForInterface(iface *net.Interface) ([]net.IP, error) {
	if iface == nil {
		return nil, errors.New("IPv4AddrsForInterface: iface cannot be nil")
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	return GetExternalIPv4Addrs(addrs)
}

// GetExternalIPv4Addrs obtains the currently-configured, non-loopback IPv4
// addresses from `addrs` coming from a particular interface (e.g.
// net.Interface.Addrs).
func GetExternalIPv4Addrs(addrs []net.Addr) ([]net.IP, error) {
	var v4addrs []net.IP
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPAddr:
			ip = v.IP
		case *net.IPNet:
			ip = v.IP
		}

		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue
		}
		v4addrs = append(v4addrs, ip)
	}
	return v4addrs, nil
}

// GenerateTransactionID generates a random 32-bits number suitable for use as
// TransactionID
func GenerateTransactionID() (TransactionID, error) {
	var xid TransactionID
	n, err := rand.Read(xid[:])
	if n != 4 {
		return xid, errors.New("invalid random sequence for transaction ID: smaller than 32 bits")
	}
	return xid, err
}

// New creates a new DHCPv4 structure and fill it up with default values. It
// won't be a valid DHCPv4 message so you will need to adjust its fields.
// See also NewDiscovery, NewOffer, NewRequest, NewAcknowledge, NewInform and
// NewRelease .
func New() (*DHCPv4, error) {
	xid, err := GenerateTransactionID()
	if err != nil {
		return nil, err
	}
	d := DHCPv4{
		OpCode:        OpcodeBootRequest,
		HWType:        iana.HWTypeEthernet,
		HopCount:      0,
		TransactionID: xid,
		NumSeconds:    0,
		Flags:         0,
		ClientIPAddr:  net.IPv4zero,
		YourIPAddr:    net.IPv4zero,
		ServerIPAddr:  net.IPv4zero,
		GatewayIPAddr: net.IPv4zero,
		Options:       make([]Option, 0, 10),
	}
	return &d, nil
}

// NewDiscoveryForInterface builds a new DHCPv4 Discovery message, with a default
// Ethernet HW type and the hardware address obtained from the specified
// interface.
func NewDiscoveryForInterface(ifname string) (*DHCPv4, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	return NewDiscovery(iface.HardwareAddr)
}

// NewDiscovery builds a new DHCPv4 Discovery message, with a default Ethernet
// HW type and specified hardware address.
func NewDiscovery(hwaddr net.HardwareAddr) (*DHCPv4, error) {
	d, err := New()
	if err != nil {
		return nil, err
	}
	// get hw addr
	d.OpCode = OpcodeBootRequest
	d.HWType = iana.HWTypeEthernet
	d.ClientHWAddr = hwaddr
	d.SetBroadcast()
	d.UpdateOption(&OptMessageType{MessageType: MessageTypeDiscover})
	d.UpdateOption(&OptParameterRequestList{
		RequestedOpts: []OptionCode{
			OptionSubnetMask,
			OptionRouter,
			OptionDomainName,
			OptionDomainNameServer,
		},
	})
	return d, nil
}

// NewInformForInterface builds a new DHCPv4 Informational message with default
// Ethernet HW type and the hardware address obtained from the specified
// interface.
func NewInformForInterface(ifname string, needsBroadcast bool) (*DHCPv4, error) {
	// get hw addr
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}

	// Set Client IP as iface's currently-configured IP.
	localIPs, err := IPv4AddrsForInterface(iface)
	if err != nil || len(localIPs) == 0 {
		return nil, fmt.Errorf("could not get local IPs for iface %s", ifname)
	}
	pkt, err := NewInform(iface.HardwareAddr, localIPs[0])
	if err != nil {
		return nil, err
	}

	if needsBroadcast {
		pkt.SetBroadcast()
	} else {
		pkt.SetUnicast()
	}
	return pkt, nil
}

// NewInform builds a new DHCPv4 Informational message with default Ethernet HW
// type and specified hardware address. It does NOT put a DHCP End option at the
// end.
func NewInform(hwaddr net.HardwareAddr, localIP net.IP) (*DHCPv4, error) {
	d, err := New()
	if err != nil {
		return nil, err
	}

	d.OpCode = OpcodeBootRequest
	d.HWType = iana.HWTypeEthernet
	d.ClientHWAddr = hwaddr
	d.ClientIPAddr = localIP
	d.UpdateOption(&OptMessageType{MessageType: MessageTypeInform})
	return d, nil
}

// NewRequestFromOffer builds a DHCPv4 request from an offer.
func NewRequestFromOffer(offer *DHCPv4, modifiers ...Modifier) (*DHCPv4, error) {
	d, err := New()
	if err != nil {
		return nil, err
	}
	d.OpCode = OpcodeBootRequest
	d.HWType = offer.HWType
	d.ClientHWAddr = offer.ClientHWAddr
	d.TransactionID = offer.TransactionID
	if offer.IsBroadcast() {
		d.SetBroadcast()
	} else {
		d.SetUnicast()
	}
	// find server IP address
	var serverIP []byte
	for _, opt := range offer.Options {
		if opt.Code() == OptionServerIdentifier {
			serverIP = opt.(*OptServerIdentifier).ServerID
		}
	}
	if serverIP == nil {
		return nil, errors.New("Missing Server IP Address in DHCP Offer")
	}
	d.ServerIPAddr = serverIP
	d.UpdateOption(&OptMessageType{MessageType: MessageTypeRequest})
	d.UpdateOption(&OptRequestedIPAddress{RequestedAddr: offer.YourIPAddr})
	d.UpdateOption(&OptServerIdentifier{ServerID: serverIP})
	for _, mod := range modifiers {
		d = mod(d)
	}
	return d, nil
}

// NewReplyFromRequest builds a DHCPv4 reply from a request.
func NewReplyFromRequest(request *DHCPv4, modifiers ...Modifier) (*DHCPv4, error) {
	reply, err := New()
	if err != nil {
		return nil, err
	}
	reply.OpCode = OpcodeBootReply
	reply.HWType = request.HWType
	reply.ClientHWAddr = request.ClientHWAddr
	reply.TransactionID = request.TransactionID
	reply.Flags = request.Flags
	reply.GatewayIPAddr = request.GatewayIPAddr
	for _, mod := range modifiers {
		reply = mod(reply)
	}
	return reply, nil
}

// FromBytes encodes the DHCPv4 packet into a sequence of bytes, and returns an
// error if the packet is not valid.
func FromBytes(q []byte) (*DHCPv4, error) {
	var p DHCPv4
	buf := uio.NewBigEndianBuffer(q)

	p.OpCode = OpcodeType(buf.Read8())
	p.HWType = iana.HWType(buf.Read8())

	hwAddrLen := buf.Read8()

	p.HopCount = buf.Read8()
	buf.ReadBytes(p.TransactionID[:])
	p.NumSeconds = buf.Read16()
	p.Flags = buf.Read16()

	p.ClientIPAddr = net.IP(buf.CopyN(net.IPv4len))
	p.YourIPAddr = net.IP(buf.CopyN(net.IPv4len))
	p.ServerIPAddr = net.IP(buf.CopyN(net.IPv4len))
	p.GatewayIPAddr = net.IP(buf.CopyN(net.IPv4len))

	if hwAddrLen > 16 {
		hwAddrLen = 16
	}
	// Always read 16 bytes, but only use hwaddrlen of them.
	p.ClientHWAddr = make(net.HardwareAddr, 16)
	buf.ReadBytes(p.ClientHWAddr)
	p.ClientHWAddr = p.ClientHWAddr[:hwAddrLen]

	var sname [64]byte
	buf.ReadBytes(sname[:])
	length := strings.Index(string(sname[:]), "\x00")
	if length == -1 {
		length = 64
	}
	p.ServerHostName = string(sname[:length])

	var file [128]byte
	buf.ReadBytes(file[:])
	length = strings.Index(string(file[:]), "\x00")
	if length == -1 {
		length = 128
	}
	p.BootFileName = string(file[:length])

	var cookie [4]byte
	buf.ReadBytes(cookie[:])

	if err := buf.Error(); err != nil {
		return nil, err
	}
	if cookie != magicCookie {
		return nil, fmt.Errorf("malformed DHCP packet: got magic cookie %v, want %v", cookie[:], magicCookie[:])
	}

	opts, err := OptionsFromBytes(buf.Data())
	if err != nil {
		return nil, err
	}
	p.Options = opts
	return &p, nil
}

// FlagsToString returns a human-readable representation of the flags field.
func (d *DHCPv4) FlagsToString() string {
	flags := ""
	if d.IsBroadcast() {
		flags += "Broadcast"
	} else {
		flags += "Unicast"
	}
	if d.Flags&0xfe != 0 {
		flags += " (reserved bits not zeroed)"
	}
	return flags
}

// IsBroadcast indicates whether the packet is a broadcast packet.
func (d *DHCPv4) IsBroadcast() bool {
	return d.Flags&0x8000 == 0x8000
}

// SetBroadcast sets the packet to be a broadcast packet.
func (d *DHCPv4) SetBroadcast() {
	d.Flags |= 0x8000
}

// IsUnicast indicates whether the packet is a unicast packet.
func (d *DHCPv4) IsUnicast() bool {
	return d.Flags&0x8000 == 0
}

// SetUnicast sets the packet to be a unicast packet.
func (d *DHCPv4) SetUnicast() {
	d.Flags &= ^uint16(0x8000)
}

// GetOneOption returns the option that matches the given option code.
//
// If no matching option is found, nil is returned.
func (d *DHCPv4) GetOneOption(code OptionCode) Option {
	return d.Options.GetOne(code)
}

// UpdateOption replaces an existing option with the same option code with the
// given one, adding it if not already present.
func (d *DHCPv4) UpdateOption(option Option) {
	d.Options.Update(option)
}

// MessageType returns the message type, trying to extract it from the
// OptMessageType option. It returns nil if the message type cannot be extracted
func (d *DHCPv4) MessageType() MessageType {
	opt := d.GetOneOption(OptionDHCPMessageType)
	if opt == nil {
		return MessageTypeNone
	}
	return opt.(*OptMessageType).MessageType
}

// HumanXID returns a human-readably integer transaction ID.
func (d *DHCPv4) HumanXID() uint32 {
	return binary.LittleEndian.Uint32(d.TransactionID[:])
}

// String implements fmt.Stringer.
func (d *DHCPv4) String() string {
	return fmt.Sprintf("DHCPv4(opcode=%s xid=%v hwtype=%s hwaddr=%s)",
		d.OpCode.String(), d.HumanXID(), d.HWType, d.ClientHWAddr)
}

// Summary prints detailed information about the packet.
func (d *DHCPv4) Summary() string {
	ret := fmt.Sprintf(
		"DHCPv4\n"+
			"  opcode=%v\n"+
			"  hwtype=%v\n"+
			"  hopcount=%v\n"+
			"  transactionid=0x%08x\n"+
			"  numseconds=%v\n"+
			"  flags=%v (0x%02x)\n"+
			"  clientipaddr=%v\n"+
			"  youripaddr=%v\n"+
			"  serveripaddr=%v\n"+
			"  gatewayipaddr=%v\n"+
			"  clienthwaddr=%v\n"+
			"  serverhostname=%v\n"+
			"  bootfilename=%v\n",
		d.OpCode,
		d.HWType,
		d.HopCount,
		d.HumanXID(),
		d.NumSeconds,
		d.FlagsToString(),
		d.Flags,
		d.ClientIPAddr,
		d.YourIPAddr,
		d.ServerIPAddr,
		d.GatewayIPAddr,
		d.ClientHWAddr,
		d.ServerHostName,
		d.BootFileName,
	)
	ret += "  options=\n"
	for _, opt := range d.Options {
		optString := opt.String()
		// If this option has sub structures, offset them accordingly.
		if strings.Contains(optString, "\n") {
			optString = strings.Replace(optString, "\n  ", "\n      ", -1)
		}
		ret += fmt.Sprintf("    %v\n", optString)
		if opt.Code() == OptionEnd {
			break
		}
	}
	return ret
}

// IsOptionRequested returns true if that option is within the requested
// options of the DHCPv4 message.
func (d *DHCPv4) IsOptionRequested(requested OptionCode) bool {
	optprl := d.GetOneOption(OptionParameterRequestList)
	if optprl == nil {
		return false
	}
	for _, o := range optprl.(*OptParameterRequestList).RequestedOpts {
		if o == requested {
			return true
		}
	}
	return false
}

// In case somebody forgets to set an IP, just write 0s as default values.
func writeIP(b *uio.Lexer, ip net.IP) {
	var zeros [net.IPv4len]byte
	if ip == nil {
		b.WriteBytes(zeros[:])
	} else {
		b.WriteBytes(ip[:net.IPv4len])
	}
}

// ToBytes writes the packet to binary.
func (d *DHCPv4) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(make([]byte, 0, minPacketLen))
	buf.Write8(uint8(d.OpCode))
	buf.Write8(uint8(d.HWType))

	// HwAddrLen
	hlen := uint8(len(d.ClientHWAddr))
	if hlen == 0 && d.HWType == iana.HWTypeEthernet {
		hlen = 6
	}
	buf.Write8(hlen)
	buf.Write8(d.HopCount)
	buf.WriteBytes(d.TransactionID[:])
	buf.Write16(d.NumSeconds)
	buf.Write16(d.Flags)

	writeIP(buf, d.ClientIPAddr)
	writeIP(buf, d.YourIPAddr)
	writeIP(buf, d.ServerIPAddr)
	writeIP(buf, d.GatewayIPAddr)
	copy(buf.WriteN(16), d.ClientHWAddr)

	var sname [64]byte
	copy(sname[:], []byte(d.ServerHostName))
	sname[len(d.ServerHostName)] = 0
	buf.WriteBytes(sname[:])

	var file [128]byte
	copy(file[:], []byte(d.BootFileName))
	file[len(d.BootFileName)] = 0
	buf.WriteBytes(file[:])

	// The magic cookie.
	buf.WriteBytes(magicCookie[:])
	d.Options.Marshal(buf)
	buf.Write8(uint8(OptionEnd))
	return buf.Data()
}
