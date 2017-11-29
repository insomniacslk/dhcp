package dhcpv6

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/insomniacslk/dhcp/dhcpv6/options"
	"github.com/insomniacslk/dhcp/iana"
	"log"
	"net"
	"time"
)

const HeaderSize = 4

type DHCPv6 struct {
	message       MessageType
	transactionID uint32 // only 24 bits are used though
	options       []options.Option
}

func BytesToTransactionID(data []byte) (*uint32, error) {
	// return a uint32 from a  sequence of bytes, representing a transaction ID.
	// Transaction IDs are three-bytes long. If the provided data is shorter than
	// 3 bytes, it return an error. If longer, will use the first three bytes
	// only.
	if len(data) < 3 {
		return nil, fmt.Errorf("Invalid transaction ID: less than 3 bytes")
	}
	buf := make([]byte, 4)
	copy(buf[1:4], data[:3])
	tid := binary.BigEndian.Uint32(buf)
	return &tid, nil
}

func GenerateTransactionID() (*uint32, error) {
	var tid *uint32
	for {
		tidBytes := make([]byte, 4)
		n, err := rand.Read(tidBytes)
		if n != 4 {
			return nil, fmt.Errorf("Invalid random sequence: shorter than 4 bytes")
		}
		tid, err = BytesToTransactionID(tidBytes)
		if err != nil {
			return nil, err
		}
		if tid == nil {
			return nil, fmt.Errorf("Error: got a nil Transaction ID")
		}
		// retry until != 0
		// TODO add retry limit
		if *tid != 0 {
			break
		}
	}
	return tid, nil
}

func FromBytes(data []byte) (*DHCPv6, error) {
	if len(data) < HeaderSize {
		return nil, fmt.Errorf("Invalid DHCPv6 header: shorter than %v bytes", HeaderSize)
	}
	tid, err := BytesToTransactionID(data[1:4])
	if err != nil {
		return nil, err
	}
	d := DHCPv6{
		message:       MessageType(data[0]),
		transactionID: *tid,
	}
	options, err := options.FromBytes(data[4:])
	if err != nil {
		return nil, err
	}
	d.options = options
	return &d, nil
}

func New() (*DHCPv6, error) {
	tid, err := GenerateTransactionID()
	if err != nil {
		return nil, err
	}
	d := DHCPv6{
		message:       SOLICIT,
		transactionID: *tid,
	}
	return &d, nil
}

// Return a time integer suitable for DUID-LLT, i.e. the current time counted in
// seconds since January 1st, 2000, midnight UTC, modulo 2^32
func GetTime() uint32 {
	now := time.Since(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	return uint32((now.Nanoseconds() / 1000000000) % 0xffffffff)
}

// Create a new SOLICIT message with DUID-LLT, using the given network
// interface's hardware address and current time
func NewSolicitForInterface(ifname string) (*DHCPv6, error) {
	d, err := New()
	if err != nil {
		return nil, err
	}
	d.SetMessage(SOLICIT)
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	cid := options.OptClientId{}
	cid.SetClientID(options.Duid{
		Type:          options.DUID_LLT,
		HwType:        iana.HwTypeEthernet,
		Time:          GetTime(),
		LinkLayerAddr: iface.HardwareAddr,
	})

	d.AddOption(&cid)
	oro := options.OptRequestedOption{}
	oro.SetRequestedOptions([]options.OptionCode{
		options.DNS_RECURSIVE_NAME_SERVER,
		options.DOMAIN_SEARCH_LIST,
	})
	d.AddOption(&oro)
	d.AddOption(&options.OptElapsedTime{})
	// FIXME use real values for IA_NA
	iaNa := options.OptIANA{}
	iaNa.SetIAID([4]byte{0x27, 0xfe, 0x8f, 0x95})
	iaNa.SetT1(0xe10)
	iaNa.SetT2(0x1518)
	d.AddOption(&iaNa)
	return d, nil
}

func (d *DHCPv6) Message() MessageType {
	return d.message
}

func (d *DHCPv6) SetMessage(message MessageType) {
	if MessageToString[message] == "" {
		log.Printf("Warning: unknown DHCPv6 message: %v", message)
	}
	d.message = message
}

func (d *DHCPv6) MessageToString() string {
	if m := MessageToString[d.message]; m != "" {
		return m
	}
	return "Invalid"
}

func (d *DHCPv6) TransactionID() uint32 {
	return d.transactionID
}

func (d *DHCPv6) SetTransactionID(tid uint32) {
	ttid := tid & 0x00ffffff
	if ttid != tid {
		log.Printf("Warning: truncating transaction ID that is longer than 24 bits: %v", tid)
	}
	d.transactionID = ttid
}

func (d *DHCPv6) Options() []options.Option {
	return d.options
}

func (d *DHCPv6) SetOptions(options []options.Option) {
	d.options = options
}

func (d *DHCPv6) AddOption(option options.Option) {
	d.options = append(d.options, option)
}

func (d *DHCPv6) String() string {
	return fmt.Sprintf("DHCPv6(message=%v transactionID=0x%06x, %d options)",
		d.MessageToString(), d.TransactionID(), len(d.options),
	)
}

func (d *DHCPv6) Summary() string {
	ret := fmt.Sprintf(
		"DHCPv6\n"+
			"  message=%v\n"+
			"  transactionid=0x%06x\n",
		d.MessageToString(),
		d.TransactionID(),
	)
	ret += "  options=["
	if len(d.options) > 0 {
		ret += "\n"
	}
	for _, opt := range d.options {
		ret += fmt.Sprintf("    %v\n", opt.String())
	}
	ret += "  ]\n"
	return ret
}

// Convert a DHCPv6 structure into its binary representation, suitable for being
// sent over the network
func (d *DHCPv6) ToBytes() []byte {
	var ret []byte
	ret = append(ret, byte(d.message))
	tidBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tidBytes, d.transactionID)
	ret = append(ret, tidBytes[1:4]...) // discard the first byte
	for _, opt := range d.options {
		ret = append(ret, opt.ToBytes()...)
	}
	return ret
}
