package dhcpv6

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/u-root/u-root/pkg/uio"
)

const MessageHeaderSize = 4

// Message represents a DHCPv6 Message as defined by RFC 3315 Section 6.
type Message struct {
	messageType   MessageType
	transactionID TransactionID
	options       Options
}

var randomRead = rand.Read

// GenerateTransactionID generates a random 3-byte transaction ID.
func GenerateTransactionID() (TransactionID, error) {
	var tid TransactionID
	n, err := randomRead(tid[:])
	if err != nil {
		return tid, err
	}
	if n != len(tid) {
		return tid, fmt.Errorf("invalid random sequence: shorter than 3 bytes")
	}
	return tid, nil
}

// GetTime returns a time integer suitable for DUID-LLT, i.e. the current time counted
// in seconds since January 1st, 2000, midnight UTC, modulo 2^32
func GetTime() uint32 {
	now := time.Since(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	return uint32((now.Nanoseconds() / 1000000000) % 0xffffffff)
}

// NewSolicitWithCID creates a new SOLICIT message with CID.
func NewSolicitWithCID(duid Duid, modifiers ...Modifier) (DHCPv6, error) {
	d, err := NewMessage()
	if err != nil {
		return nil, err
	}
	d.(*Message).SetMessage(MessageTypeSolicit)
	d.AddOption(&OptClientId{Cid: duid})
	oro := new(OptRequestedOption)
	oro.SetRequestedOptions([]OptionCode{
		OptionDNSRecursiveNameServer,
		OptionDomainSearchList,
	})
	d.AddOption(oro)
	d.AddOption(&OptElapsedTime{})
	// FIXME use real values for IA_NA
	iaNa := &OptIANA{}
	iaNa.IaId = [4]byte{0xfa, 0xce, 0xb0, 0x0c}
	iaNa.T1 = 0xe10
	iaNa.T2 = 0x1518
	d.AddOption(iaNa)
	// Apply modifiers
	for _, mod := range modifiers {
		d = mod(d)
	}
	return d, nil
}

// NewSolicitForInterface creates a new SOLICIT message with DUID-LLT, using the
// given network interface's hardware address and current time
func NewSolicitForInterface(ifname string, modifiers ...Modifier) (DHCPv6, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	duid := Duid{
		Type:          DUID_LLT,
		HwType:        iana.HWTypeEthernet,
		Time:          GetTime(),
		LinkLayerAddr: iface.HardwareAddr,
	}
	return NewSolicitWithCID(duid, modifiers...)
}

// NewAdvertiseFromSolicit creates a new ADVERTISE packet based on an SOLICIT packet.
func NewAdvertiseFromSolicit(solicit DHCPv6, modifiers ...Modifier) (DHCPv6, error) {
	if solicit == nil {
		return nil, errors.New("SOLICIT cannot be nil")
	}
	if solicit.Type() != MessageTypeSolicit {
		return nil, errors.New("The passed SOLICIT must have SOLICIT type set")
	}
	sol, ok := solicit.(*Message)
	if !ok {
		return nil, errors.New("The passed SOLICIT must be of Message type")
	}
	// build ADVERTISE from SOLICIT
	adv := Message{}
	adv.SetMessage(MessageTypeAdvertise)
	adv.SetTransactionID(sol.TransactionID())
	// add Client ID
	cid := sol.GetOneOption(OptionClientID)
	if cid == nil {
		return nil, errors.New("Client ID cannot be nil in SOLICIT when building ADVERTISE")
	}
	adv.AddOption(cid)

	// apply modifiers
	d := DHCPv6(&adv)
	for _, mod := range modifiers {
		d = mod(d)
	}
	return d, nil
}

// NewRequestFromAdvertise creates a new REQUEST packet based on an ADVERTISE
// packet options.
func NewRequestFromAdvertise(advertise DHCPv6, modifiers ...Modifier) (DHCPv6, error) {
	if advertise == nil {
		return nil, fmt.Errorf("ADVERTISE cannot be nil")
	}
	if advertise.Type() != MessageTypeAdvertise {
		return nil, fmt.Errorf("The passed ADVERTISE must have ADVERTISE type set")
	}
	adv, ok := advertise.(*Message)
	if !ok {
		return nil, fmt.Errorf("The passed ADVERTISE must be of Message type")
	}
	// build REQUEST from ADVERTISE
	req := Message{}
	req.SetMessage(MessageTypeRequest)
	req.SetTransactionID(adv.TransactionID())
	// add Client ID
	cid := adv.GetOneOption(OptionClientID)
	if cid == nil {
		return nil, fmt.Errorf("Client ID cannot be nil in ADVERTISE when building REQUEST")
	}
	req.AddOption(cid)
	// add Server ID
	sid := adv.GetOneOption(OptionServerID)
	if sid == nil {
		return nil, fmt.Errorf("Server ID cannot be nil in ADVERTISE when building REQUEST")
	}
	req.AddOption(sid)
	// add Elapsed Time
	req.AddOption(&OptElapsedTime{})
	// add IA_NA
	iaNa := adv.GetOneOption(OptionIANA)
	if iaNa == nil {
		return nil, fmt.Errorf("IA_NA cannot be nil in ADVERTISE when building REQUEST")
	}
	req.AddOption(iaNa)
	// add OptRequestedOption
	oro := OptRequestedOption{}
	oro.SetRequestedOptions([]OptionCode{
		OptionDNSRecursiveNameServer,
		OptionDomainSearchList,
	})
	req.AddOption(&oro)
	// add OPTION_VENDOR_CLASS, only if present in the original request
	// TODO implement OptionVendorClass
	vClass := adv.GetOneOption(OptionVendorClass)
	if vClass != nil {
		req.AddOption(vClass)
	}

	// apply modifiers
	d := DHCPv6(&req)
	for _, mod := range modifiers {
		d = mod(d)
	}
	return d, nil
}

// NewReplyFromMessage creates a new REPLY packet based on a
// Message. The function is to be used when generating a reply to
// REQUEST, CONFIRM, RENEW, REBIND, RELEASE and INFORMATION-REQUEST packets.
func NewReplyFromMessage(message DHCPv6, modifiers ...Modifier) (DHCPv6, error) {
	if message == nil {
		return nil, errors.New("Message cannot be nil")
	}
	switch message.Type() {
	case MessageTypeRequest, MessageTypeConfirm, MessageTypeRenew,
		MessageTypeRebind, MessageTypeRelease, MessageTypeInformationRequest:
	default:
		return nil, errors.New("Cannot create REPLY from the passed message type set")
	}
	msg, ok := message.(*Message)
	if !ok {
		return nil, errors.New("The passed MESSAGE must be of Message type")
	}
	// build REPLY from MESSAGE
	rep := Message{}
	rep.SetMessage(MessageTypeReply)
	rep.SetTransactionID(msg.TransactionID())
	// add Client ID
	cid := message.GetOneOption(OptionClientID)
	if cid == nil {
		return nil, errors.New("Client ID cannot be nil when building REPLY")
	}
	rep.AddOption(cid)

	// apply modifiers
	d := DHCPv6(&rep)
	for _, mod := range modifiers {
		d = mod(d)
	}
	return d, nil
}

// Type is the DHCPv6 message type.
func (d *Message) Type() MessageType {
	return d.messageType
}

// SetMessage sets the DHCP message type.
func (d *Message) SetMessage(messageType MessageType) {
	msgString := messageType.String()
	if msgString == "" {
		log.Printf("Warning: unknown DHCPv6 message type: %v", messageType)
	}
	if messageType == MessageTypeRelayForward || messageType == MessageTypeRelayReply {
		log.Printf("Warning: using a RELAY message type with a non-relay message: %v (%v)",
			msgString, messageType)
	}
	d.messageType = messageType
}

// TransactionID returns this message's transaction id.
func (d *Message) TransactionID() TransactionID {
	return d.transactionID
}

// SetTransactionID sets this message's transaction id.
func (d *Message) SetTransactionID(tid TransactionID) {
	d.transactionID = tid
}

// SetOptions replaces this message's options.
func (d *Message) SetOptions(options []Option) {
	d.options = options
}

// AddOption adds an option to this message.
func (d *Message) AddOption(option Option) {
	d.options.Add(option)
}

// UpdateOption updates the existing options with the passed option, adding it
// at the end if not present already
func (d *Message) UpdateOption(option Option) {
	d.options.Update(option)
}

// IsNetboot returns true if the machine is trying to netboot. It checks if
// "boot file" is one of the requested options, which is useful for
// SOLICIT/REQUEST packet types, it also checks if the "boot file" option is
// included in the packet, which is useful for ADVERTISE/REPLY packet.
func (d *Message) IsNetboot() bool {
	if d.IsOptionRequested(OptionBootfileURL) {
		return true
	}
	if optbf := d.GetOneOption(OptionBootfileURL); optbf != nil {
		return true
	}
	return false
}

// IsOptionRequested takes an OptionCode and returns true if that option is
// within the requested options of the DHCPv6 message.
func (d *Message) IsOptionRequested(requested OptionCode) bool {
	for _, optoro := range d.GetOption(OptionORO) {
		for _, o := range optoro.(*OptRequestedOption).RequestedOptions() {
			if o == requested {
				return true
			}
		}
	}
	return false
}

// String returns a short human-readable string for this message.
func (d *Message) String() string {
	return fmt.Sprintf("Message(messageType=%v transactionID=%s, %d options)",
		d.Type().String(), d.TransactionID(), len(d.options),
	)
}

// Summary prints all options associated with this message.
func (d *Message) Summary() string {
	ret := fmt.Sprintf(
		"Message\n"+
			"  messageType=%v\n"+
			"  transactionid=%s\n",
		d.Type().String(),
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

// ToBytes returns the serialized version of this message as defined by RFC
// 3315, Section 5.
func (d *Message) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(uint8(d.messageType))
	buf.WriteBytes(d.transactionID[:])
	buf.WriteBytes(d.options.ToBytes())
	return buf.Data()
}

// Options returns the current set of options associated with this message.
func (d *Message) Options() []Option {
	return d.options
}

// GetOption returns the options associated with the code.
func (d *Message) GetOption(code OptionCode) []Option {
	return d.options.Get(code)
}

// GetOneOption returns the first associated option with the code from this
// message.
func (d *Message) GetOneOption(code OptionCode) Option {
	return d.options.GetOne(code)
}

// IsRelay returns whether this is a relay message or not.
func (d *Message) IsRelay() bool {
	return false
}
