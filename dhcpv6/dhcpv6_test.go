package dhcpv6

import (
	"encoding/binary"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insomniacslk/dhcp/iana"
)

func TestBytesToTransactionID(t *testing.T) {
	// Check if the function transforms the bytes for the exact length input
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, 0x01020304)
	tid, err := BytesToTransactionID(b)
	require.NoError(t, err)
	require.NotNil(t, tid)
	assert.Equal(t, *tid, uint32(0x000040302))

	binary.BigEndian.PutUint32(b, 0x01020304)
	tid, err = BytesToTransactionID(b)
	require.NoError(t, err)
	require.NotNil(t, tid)
	assert.Equal(t, *tid, uint32(0x00010203))

	// Check if the function transforms only the first bytes for a longer input
	b = make([]byte, 8)
	binary.LittleEndian.PutUint32(b, 0x01020304)
	binary.LittleEndian.PutUint32(b[4:], 0x11121314)
	tid, err = BytesToTransactionID(b)
	require.NoError(t, err)
	require.NotNil(t, tid)
	assert.Equal(t, *tid, uint32(0x000040302))

	binary.BigEndian.PutUint32(b, 0x01020304)
	binary.BigEndian.PutUint32(b[4:], 0x11121314)
	tid, err = BytesToTransactionID(b)
	require.NoError(t, err)
	require.NotNil(t, tid)
	assert.Equal(t, *tid, uint32(0x00010203))
}

func TestBytesToTransactionIDShortData(t *testing.T) {
	// short sequence, less than three bytes
	tid, err := BytesToTransactionID([]byte{0x11, 0x22})
	require.Error(t, err)
	require.Nil(t, tid)
}

func TestGenerateTransactionID(t *testing.T) {
	tid, err := GenerateTransactionID()
	require.NoError(t, err)
	require.NotNil(t, tid)
	require.True(t, *tid <= 0xffffff, "transaction ID should be smaller than 0xffffff")
}

func TestNewMessage(t *testing.T) {
	d, err := NewMessage()
	require.NoError(t, err)
	require.NotNil(t, d)
	require.Equal(t, MessageTypeSolicit, d.Type())
	require.NotEqual(t, 0, d.(*DHCPv6Message).transactionID)
	require.Empty(t, d.(*DHCPv6Message).options)
}

func TestDecapsulateRelayIndex(t *testing.T) {
	m := DHCPv6Message{}
	r1, err := EncapsulateRelay(&m, MessageTypeRelayForward, net.IPv6linklocalallnodes, net.IPv6interfacelocalallnodes)
	require.NoError(t, err)
	r2, err := EncapsulateRelay(r1, MessageTypeRelayForward, net.IPv6loopback, net.IPv6linklocalallnodes)
	require.NoError(t, err)
	r3, err := EncapsulateRelay(r2, MessageTypeRelayForward, net.IPv6unspecified, net.IPv6linklocalallrouters)
	require.NoError(t, err)

	first, err := DecapsulateRelayIndex(r3, 0)
	require.NoError(t, err)
	relay, ok := first.(*DHCPv6Relay)
	require.True(t, ok)
	require.Equal(t, relay.HopCount(), uint8(1))
	require.Equal(t, relay.LinkAddr(), net.IPv6loopback)
	require.Equal(t, relay.PeerAddr(), net.IPv6linklocalallnodes)

	second, err := DecapsulateRelayIndex(r3, 1)
	require.NoError(t, err)
	relay, ok = second.(*DHCPv6Relay)
	require.True(t, ok)
	require.Equal(t, relay.HopCount(), uint8(0))
	require.Equal(t, relay.LinkAddr(), net.IPv6linklocalallnodes)
	require.Equal(t, relay.PeerAddr(), net.IPv6interfacelocalallnodes)

	third, err := DecapsulateRelayIndex(r3, 2)
	require.NoError(t, err)
	_, ok = third.(*DHCPv6Message)
	require.True(t, ok)

	rfirst, err := DecapsulateRelayIndex(r3, -1)
	require.NoError(t, err)
	relay, ok = rfirst.(*DHCPv6Relay)
	require.True(t, ok)
	require.Equal(t, relay.HopCount(), uint8(0))
	require.Equal(t, relay.LinkAddr(), net.IPv6linklocalallnodes)
	require.Equal(t, relay.PeerAddr(), net.IPv6interfacelocalallnodes)

	_, err = DecapsulateRelayIndex(r3, -2)
	require.Error(t, err)
}

func TestSettersAndGetters(t *testing.T) {
	d := DHCPv6Message{}
	// Message
	d.SetMessage(MessageTypeSolicit)
	require.Equal(t, MessageTypeSolicit, d.Type())
	d.SetMessage(MessageTypeAdvertise)
	require.Equal(t, MessageTypeAdvertise, d.Type())

	// TransactionID
	d.SetTransactionID(12345)
	require.Equal(t, uint32(12345), d.TransactionID())

	// Options
	require.Empty(t, d.Options())
	expectedOptions := []Option{&OptionGeneric{OptionCode: 0, OptionData: []byte{}}}
	d.SetOptions(expectedOptions)
	require.Equal(t, expectedOptions, d.Options())
}

func TestAddOption(t *testing.T) {
	d := DHCPv6Message{}
	require.Empty(t, d.Options())
	opt := OptionGeneric{OptionCode: 0, OptionData: []byte{}}
	d.AddOption(&opt)
	require.Equal(t, []Option{&opt}, d.Options())
}

func TestToBytes(t *testing.T) {
	d := DHCPv6Message{}
	d.SetMessage(MessageTypeSolicit)
	d.SetTransactionID(0xabcdef)
	opt := OptionGeneric{OptionCode: 0, OptionData: []byte{}}
	d.AddOption(&opt)
	bytes := d.ToBytes()
	expected := []byte{01, 0xab, 0xcd, 0xef, 0x00, 0x00, 0x00, 0x00}
	require.Equal(t, expected, bytes)
}

func TestFromAndToBytes(t *testing.T) {
	expected := []byte{01, 0xab, 0xcd, 0xef, 0x00, 0x00, 0x00, 0x00}
	d, err := FromBytes(expected)
	require.NoError(t, err)
	toBytes := d.ToBytes()
	require.Equal(t, expected, toBytes)
}

func TestNewAdvertiseFromSolicit(t *testing.T) {
	s := DHCPv6Message{}
	s.SetMessage(MessageTypeSolicit)
	s.SetTransactionID(0xabcdef)
	cid := OptClientId{}
	s.AddOption(&cid)
	duid := Duid{}

	a, err := NewAdvertiseFromSolicit(&s, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, a.(*DHCPv6Message).TransactionID(), s.TransactionID())
	require.Equal(t, a.Type(), MessageTypeAdvertise)
}

func TestNewReplyFromDHCPv6Message(t *testing.T) {
	msg := DHCPv6Message{}
	msg.SetTransactionID(0xabcdef)
	cid := OptClientId{}
	msg.AddOption(&cid)
	sid := OptServerId{}
	duid := Duid{}
	sid.Sid = duid
	msg.AddOption(&sid)

	msg.SetMessage(MessageTypeConfirm)
	rep, err := NewReplyFromDHCPv6Message(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), msg.TransactionID())
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.SetMessage(MessageTypeRenew)
	rep, err = NewReplyFromDHCPv6Message(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), msg.TransactionID())
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.SetMessage(MessageTypeRebind)
	rep, err = NewReplyFromDHCPv6Message(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), msg.TransactionID())
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.SetMessage(MessageTypeRelease)
	rep, err = NewReplyFromDHCPv6Message(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), msg.TransactionID())
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.SetMessage(MessageTypeInformationRequest)
	rep, err = NewReplyFromDHCPv6Message(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), msg.TransactionID())
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.SetMessage(MessageTypeSolicit)
	rep, err = NewReplyFromDHCPv6Message(&msg)
	require.Error(t, err)

	relay := DHCPv6Relay{}
	rep, err = NewReplyFromDHCPv6Message(&relay)
	require.Error(t, err)
}

func TestNewMessageTypeSolicitWithCID(t *testing.T) {
	hwAddr, err := net.ParseMAC("24:0A:9E:9F:EB:2B")
	require.NoError(t, err)

	duid := Duid{
		Type:          DUID_LL,
		HwType:        iana.HwTypeEthernet,
		LinkLayerAddr: hwAddr,
	}

	s, err := NewSolicitWithCID(duid)
	require.NoError(t, err)

	require.Equal(t, s.Type(), MessageTypeSolicit)
	// Check CID
	cidOption := s.GetOneOption(OptionClientID)
	require.NotNil(t, cidOption)
	cid, ok := cidOption.(*OptClientId)
	require.True(t, ok)
	require.Equal(t, cid.Cid, duid)

	// Check ORO
	oroOption := s.GetOneOption(OptionORO)
	require.NotNil(t, oroOption)
	oro, ok := oroOption.(*OptRequestedOption)
	require.True(t, ok)
	opts := oro.RequestedOptions()
	require.Contains(t, opts, OptionDNSRecursiveNameServer)
	require.Contains(t, opts, OptionDomainSearchList)
	require.Equal(t, len(opts), 2)
}


func TestIsUsingUEFIArchTypeTrue(t *testing.T) {
	msg := DHCPv6Message{}
	opt := OptClientArchType{ArchTypes: []iana.ArchType{iana.EFI_BC}}
	msg.AddOption(&opt)
	require.True(t, IsUsingUEFI(&msg))
}

func TestIsUsingUEFIArchTypeFalse(t *testing.T) {
	msg := DHCPv6Message{}
	opt := OptClientArchType{ArchTypes: []iana.ArchType{iana.INTEL_X86PC}}
	msg.AddOption(&opt)
	require.False(t, IsUsingUEFI(&msg))
}

func TestIsUsingUEFIUserClassTrue(t *testing.T) {
	msg := DHCPv6Message{}
	opt := OptUserClass{UserClasses: [][]byte{[]byte("ipxeUEFI")}}
	msg.AddOption(&opt)
	require.True(t, IsUsingUEFI(&msg))
}

func TestIsUsingUEFIUserClassFalse(t *testing.T) {
	msg := DHCPv6Message{}
	opt := OptUserClass{UserClasses: [][]byte{[]byte("ipxeLegacy")}}
	msg.AddOption(&opt)
	require.False(t, IsUsingUEFI(&msg))
}

func TestGetTransactionIDMessage(t *testing.T) {
	message, err := NewMessage()
	require.NoError(t, err)
	transactionID, err := GetTransactionID(message)
	require.NoError(t, err)
	require.Equal(t, transactionID, message.(*DHCPv6Message).TransactionID())
}

func TestGetTransactionIDRelay(t *testing.T) {
	message, err := NewMessage()
	require.NoError(t, err)
	relay, err := EncapsulateRelay(message, MessageTypeRelayForward, nil, nil)
	require.NoError(t, err)
	transactionID, err := GetTransactionID(relay)
	require.NoError(t, err)
	require.Equal(t, transactionID, message.(*DHCPv6Message).TransactionID())
}

// TODO test NewMessageTypeSolicit
//      test String and Summary
