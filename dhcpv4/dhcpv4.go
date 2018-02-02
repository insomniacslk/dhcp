package dhcpv4

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/insomniacslk/dhcp/iana"
	"log"
	"net"
	"strings"
)

const HeaderSize = 236
const MaxMessageSize = 576

type DHCPv4 struct {
	opcode         OpcodeType
	hwType         iana.HwTypeType
	hwAddrLen      uint8
	hopCount       uint8
	transactionID  uint32
	numSeconds     uint16
	flags          uint16
	clientIPAddr   net.IP
	yourIPAddr     net.IP
	serverIPAddr   net.IP
	gatewayIPAddr  net.IP
	clientHwAddr   [16]byte
	serverHostName [64]byte
	bootFileName   [128]byte
	options        []Option
}

// Generate a random 32-bits number suitable as TransactionID
func GenerateTransactionID() (*uint32, error) {
	b := make([]byte, 4)
	n, err := rand.Read(b)
	if n != 4 {
		return nil, errors.New("Invalid random sequence: smaller than 32 bits")
	}
	if err != nil {
		return nil, err
	}
	tid := binary.LittleEndian.Uint32(b)
	return &tid, nil
}

// New creates a new DHCPv4 structure and fill it up with default values. It
// won't be a valid DHCPv4 message so you will need to adjust its fields.
// See also NewDiscovery, NewOffer, NewRequest, NewAcknowledge, NewInform and
// NewRelease .
func New() (*DHCPv4, error) {
	tid, err := GenerateTransactionID()
	if err != nil {
		return nil, err
	}
	d := DHCPv4{
		opcode:        OpcodeBootRequest,
		hwType:        iana.HwTypeEthernet,
		hwAddrLen:     6,
		hopCount:      0,
		transactionID: *tid,
		numSeconds:    0,
		flags:         0,
		clientIPAddr:  net.IPv4zero,
		yourIPAddr:    net.IPv4zero,
		serverIPAddr:  net.IPv4zero,
		gatewayIPAddr: net.IPv4zero,
	}
	copy(d.clientHwAddr[:], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	copy(d.serverHostName[:], []byte{})
	copy(d.bootFileName[:], []byte{})
	options, err := OptionsFromBytes(MagicCookie)
	if err != nil {
		return nil, err
	}
	d.options = options
	return &d, nil
}

// NewDiscoveryForInterface builds a new DHCPv4 Discovery message, with a default
// Ethernet HW type and the hardware address obtained from the specified
// interface.
func NewDiscoveryForInterface(ifname string) (*DHCPv4, error) {
	d, err := New()
	if err != nil {
		return nil, err
	}
	// get hw addr
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	d.SetOpcode(OpcodeBootRequest)
	d.SetHwType(iana.HwTypeEthernet)
	d.SetHwAddrLen(uint8(len(iface.HardwareAddr)))
	d.SetClientHwAddr(iface.HardwareAddr)
	d.SetBroadcast()
	d.AddOption(Option{
		Code: OptionDHCPMessageType,
		Data: []byte{1},
	})
	d.AddOption(Option{
		Code: OptionParameterRequestList,
		Data: []byte{OptionSubnetMask, OptionRouter, OptionDomainName, OptionDomainNameServer},
	})
	// the End option has to be added explicitly
	d.AddOption(Option{Code: OptionEnd})
	return d, nil
}

// RequestFromOffer builds a DHCPv4 request from an offer.
func RequestFromOffer(offer DHCPv4) (*DHCPv4, error) {
	d, err := New()
	if err != nil {
		return nil, err
	}
	d.SetOpcode(OpcodeBootRequest)
	d.SetHwType(offer.HwType())
	d.SetHwAddrLen(offer.HwAddrLen())
	hwaddr := offer.ClientHwAddr()
	d.SetClientHwAddr(hwaddr[:])
	d.SetTransactionID(offer.TransactionID())
	if offer.IsBroadcast() {
		d.SetBroadcast()
	} else {
		d.SetUnicast()
	}
	// find server IP address
	var serverIP []byte
	for _, opt := range offer.options {
		if opt.Code == OptionServerIdentifier {
			serverIP = make([]byte, 4)
			copy(serverIP, opt.Data)
		}
	}
	if serverIP == nil {
		return nil, errors.New("Missing Server IP Address in DHCP Offer")
	}
	d.SetServerIPAddr(serverIP)
	d.AddOption(Option{
		Code: OptionDHCPMessageType,
		Data: []byte{3},
	})
	d.AddOption(Option{
		Code: OptionRequestedIPAddress,
		Data: offer.YourIPAddr(),
	})
	d.AddOption(Option{
		Code: OptionServerIdentifier,
		Data: serverIP,
	})
	// the End option has to be added explicitly
	d.AddOption(Option{Code: OptionEnd})
	return d, nil
}

// FromBytes encodes the DHCPv4 packet into a sequence of bytes, and returns an
// error if the packet is not valid.
func FromBytes(data []byte) (*DHCPv4, error) {
	if len(data) < HeaderSize {
		return nil, fmt.Errorf("Invalid DHCPv4 header: shorter than %v bytes", HeaderSize)
	}
	d := DHCPv4{
		opcode:        OpcodeType(data[0]),
		hwType:        iana.HwTypeType(data[1]),
		hwAddrLen:     data[2],
		hopCount:      data[3],
		transactionID: binary.BigEndian.Uint32(data[4:8]),
		numSeconds:    binary.BigEndian.Uint16(data[8:10]),
		flags:         binary.BigEndian.Uint16(data[10:12]),
		clientIPAddr:  net.IP(data[12:16]),
		yourIPAddr:    net.IP(data[16:20]),
		serverIPAddr:  net.IP(data[20:24]),
		gatewayIPAddr: net.IP(data[24:28]),
	}
	copy(d.clientHwAddr[:], data[28:44])
	copy(d.serverHostName[:], data[44:108])
	copy(d.bootFileName[:], data[108:236])
	options, err := OptionsFromBytes(data[236:])
	if err != nil {
		return nil, err
	}
	d.options = options
	return &d, nil
}

// Opcode returns the OpcodeType for the packet,
func (d *DHCPv4) Opcode() OpcodeType {
	return d.opcode
}

// OpcodeToString returns the mnemonic name for the packet's opcode.
func (d *DHCPv4) OpcodeToString() string {
	opcode := OpcodeToString[d.opcode]
	if opcode == "" {
		opcode = "Invalid"
	}
	return opcode
}

// SetOpcode sets a new opcode for the packet. It prints a warning if the opcode
// is unknown, but does not generate an error.
func (d *DHCPv4) SetOpcode(opcode OpcodeType) {
	if OpcodeToString[opcode] == "" {
		log.Printf("Warning: unknown DHCPv4 opcode: %v", opcode)
	}
	d.opcode = opcode
}

// HwType returns the hardware type as defined by IANA.
func (d *DHCPv4) HwType() iana.HwTypeType {
	return d.hwType
}

// HwTypeToString returns the mnemonic name for the hardware type, e.g.
// "Ethernet". If the type is unknown, it returns "Unknown".
func (d *DHCPv4) HwTypeToString() string {
	hwtype, ok := iana.HwTypeToString[d.hwType]
	if !ok {
		hwtype = "Invalid"
	}
	return hwtype
}

// SetHwType returns the hardware type as defined by IANA.
func (d *DHCPv4) SetHwType(hwType iana.HwTypeType) {
	if _, ok := iana.HwTypeToString[hwType]; !ok {
		log.Printf("Warning: Invalid DHCPv4 hwtype: %v", hwType)
	}
	d.hwType = hwType
}

// HwAddrLen returns the hardware address length. E.g. for Ethernet it would
// return 6.
func (d *DHCPv4) HwAddrLen() uint8 {
	return d.hwAddrLen
}

// SetHwAddrLen sets the hardware address length, limiting it to the maximum
// size 16 that the standard allows.
func (d *DHCPv4) SetHwAddrLen(hwAddrLen uint8) {
	if hwAddrLen > 16 {
		log.Printf("Warning: invalid HwAddrLen: %v > 16, using 16 instead", hwAddrLen)
		hwAddrLen = 16
	}
	d.hwAddrLen = hwAddrLen
}

// HopCount returns the hop count field.
func (d *DHCPv4) HopCount() uint8 {
	return d.hopCount
}

// SetHopCount sets the hop count value.
func (d *DHCPv4) SetHopCount(hopCount uint8) {
	d.hopCount = hopCount
}

// TransactionID returns the transaction ID as 32 bit unsigned integer.
func (d *DHCPv4) TransactionID() uint32 {
	return d.transactionID
}

// SetTransactionID sets the value for the transaction ID.
func (d *DHCPv4) SetTransactionID(transactionID uint32) {
	d.transactionID = transactionID
}

func (d *DHCPv4) NumSeconds() uint16 {
	return d.numSeconds
}

func (d *DHCPv4) SetNumSeconds(numSeconds uint16) {
	d.numSeconds = numSeconds
}

func (d *DHCPv4) Flags() uint16 {
	return d.flags
}

func (d *DHCPv4) SetFlags(flags uint16) {
	d.flags = flags
}

func (d *DHCPv4) FlagsToString() string {
	flags := ""
	if d.IsBroadcast() {
		flags += "Broadcast"
	} else {
		flags += "Unicast"
	}
	if d.flags&0xfe != 0 {
		flags += " (reserved bits not zeroed)"
	}
	return flags
}

func (d *DHCPv4) IsBroadcast() bool {
	return d.flags&0x8000 == 0x8000
}

func (d *DHCPv4) SetBroadcast() {
	d.flags |= 0x8000
}

func (d *DHCPv4) IsUnicast() bool {
	return d.flags&0x8000 == 0
}

func (d *DHCPv4) SetUnicast() {
	d.flags &= ^uint16(0x8000)
}

func (d *DHCPv4) ClientIPAddr() net.IP {
	return d.clientIPAddr
}

func (d *DHCPv4) SetClientIPAddr(clientIPAddr net.IP) {
	d.clientIPAddr = clientIPAddr
}

func (d *DHCPv4) YourIPAddr() net.IP {
	return d.yourIPAddr
}

func (d *DHCPv4) SetYourIPAddr(yourIPAddr net.IP) {
	d.yourIPAddr = yourIPAddr
}

func (d *DHCPv4) ServerIPAddr() net.IP {
	return d.serverIPAddr
}

func (d *DHCPv4) SetServerIPAddr(serverIPAddr net.IP) {
	d.serverIPAddr = serverIPAddr
}

func (d *DHCPv4) GatewayIPAddr() net.IP {
	return d.gatewayIPAddr
}

func (d *DHCPv4) SetGatewayIPAddr(gatewayIPAddr net.IP) {
	d.gatewayIPAddr = gatewayIPAddr
}

func (d *DHCPv4) ClientHwAddr() [16]byte {
	return d.clientHwAddr
}

func (d *DHCPv4) ClientHwAddrToString() string {
	var ret []string
	for _, b := range d.clientHwAddr[:d.hwAddrLen] {
		ret = append(ret, fmt.Sprintf("%02x", b))
	}
	return strings.Join(ret, ":")
}

func (d *DHCPv4) SetClientHwAddr(clientHwAddr []byte) {
	if len(clientHwAddr) > 16 {
		log.Printf("Warning: too long HW Address (%d bytes), truncating to 16 bytes", len(clientHwAddr))
		clientHwAddr = clientHwAddr[:16]
	}
	copy(d.clientHwAddr[:len(clientHwAddr)], clientHwAddr)
	// pad the remaining bytes, if any
	for i := len(clientHwAddr); i < 16; i++ {
		d.clientHwAddr[i] = 0
	}
}

func (d *DHCPv4) ServerHostName() [64]byte {
	return d.serverHostName
}

func (d *DHCPv4) ServerHostNameToString() string {
	return strings.TrimRight(string(d.serverHostName[:]), "\x00")
}

func (d *DHCPv4) SetServerHostName(serverHostName []byte) {
	if len(serverHostName) > 64 {
		serverHostName = serverHostName[:64]
	} else if len(serverHostName) < 64 {
		for i := len(serverHostName) - 1; i < 64; i++ {
			serverHostName = append(serverHostName, 0)
		}
	}
	// need an array, not a slice, so let's copy it
	var newServerHostName [64]byte
	copy(newServerHostName[:], serverHostName)
	d.serverHostName = newServerHostName
}

func (d *DHCPv4) BootFileName() [128]byte {
	return d.bootFileName
}

func (d *DHCPv4) BootFileNameToString() string {
	return strings.TrimRight(string(d.bootFileName[:]), "\x00")
}

func (d *DHCPv4) SetBootFileName(bootFileName []byte) {
	if len(bootFileName) > 128 {
		bootFileName = bootFileName[:128]
	} else if len(bootFileName) < 128 {
		for i := len(bootFileName) - 1; i < 128; i++ {
			bootFileName = append(bootFileName, 0)
		}
	}
	// need an array, not a slice, so let's copy it
	var newBootFileName [128]byte
	copy(newBootFileName[:], bootFileName)
	d.bootFileName = newBootFileName
}

func (d *DHCPv4) Options() []Option {
	return d.options
}

func (d *DHCPv4) StrippedOptions() []Option {
	// differently from Options() this function strips away anything coming
	// after the End option (normally just Pad options).
	strippedOptions := []Option{}
	for _, opt := range d.options {
		strippedOptions = append(strippedOptions, opt)
		if opt.Code == OptionEnd {
			break
		}
	}
	return strippedOptions
}

func (d *DHCPv4) SetOptions(options []Option) {
	d.options = options
}

func (d *DHCPv4) AddOption(option Option) {
	d.options = append(d.options, option)
}

func (d *DHCPv4) String() string {
	return fmt.Sprintf("DHCPv4(opcode=%v hwtype=%v hwaddr=%v)",
		d.OpcodeToString(), d.HwTypeToString(), d.ClientHwAddr())
}

func (d *DHCPv4) Summary() string {
	ret := fmt.Sprintf(
		"DHCPv4\n"+
			"  opcode=%v\n"+
			"  hwtype=%v\n"+
			"  hwaddrlen=%v\n"+
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
		d.OpcodeToString(),
		d.HwTypeToString(),
		d.HwAddrLen(),
		d.HopCount(),
		d.TransactionID(),
		d.NumSeconds(),
		d.FlagsToString(),
		d.Flags(),
		d.ClientIPAddr(),
		d.YourIPAddr(),
		d.ServerIPAddr(),
		d.GatewayIPAddr(),
		d.ClientHwAddrToString(),
		d.ServerHostNameToString(),
		d.BootFileNameToString(),
	)
	ret += "  options=\n"
	for _, opt := range d.options {
		ret += fmt.Sprintf("    %v\n", opt.String())
		if opt.Code == OptionEnd {
			break
		}
	}
	return ret
}

func (d *DHCPv4) ValidateOptions() {
	// TODO find duplicate options
	foundOptionEnd := false
	for _, opt := range d.options {
		if foundOptionEnd {
			if opt.Code == OptionEnd {
				log.Print("Warning: found duplicate End option")
			}
			if opt.Code != OptionEnd && opt.Code != OptionPad {
				name := OptionCodeToString[opt.Code]
				log.Printf("Warning: found option %v (%v) after End option", opt.Code, name)
			}
		}
		if opt.Code == OptionEnd {
			foundOptionEnd = true
		}
	}
	if !foundOptionEnd {
		log.Print("Warning: no End option found")
	}
}

// Convert a DHCPv4 structure into its binary representation, suitable for being
// sent over the network
func (d *DHCPv4) ToBytes() []byte {
	// This won't check if the End option is present, you've been warned
	var ret []byte
	u32 := make([]byte, 4)
	u16 := make([]byte, 2)

	ret = append(ret, byte(d.opcode))
	ret = append(ret, byte(d.hwType))
	ret = append(ret, byte(d.hwAddrLen))
	ret = append(ret, byte(d.hopCount))
	binary.BigEndian.PutUint32(u32, d.transactionID)
	ret = append(ret, u32...)
	binary.BigEndian.PutUint16(u16, d.numSeconds)
	ret = append(ret, u16...)
	binary.BigEndian.PutUint16(u16, d.flags)
	ret = append(ret, u16...)
	ret = append(ret, d.clientIPAddr[:4]...)
	ret = append(ret, d.yourIPAddr[:4]...)
	ret = append(ret, d.serverIPAddr[:4]...)
	ret = append(ret, d.gatewayIPAddr[:4]...)
	ret = append(ret, d.clientHwAddr[:16]...)
	ret = append(ret, d.serverHostName[:64]...)
	ret = append(ret, d.bootFileName[:128]...)
	d.ValidateOptions() // print warnings about broken options, if any
	ret = append(ret, OptionsToBytes(d.options)...)
	return ret
}
