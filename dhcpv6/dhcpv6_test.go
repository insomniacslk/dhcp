package dhcpv6

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insomniacslk/dhcp/iana"
)

func randomReadMock(value []byte, n int, err error) func([]byte) (int, error) {
	return func(b []byte) (int, error) {
		copy(b, value)
		return n, err
	}
}

type GenerateTransactionIDTestSuite struct {
	suite.Suite
	random []byte
}

func (s *GenerateTransactionIDTestSuite) SetupTest() {
	s.random = make([]byte, 16)
}

func (s *GenerateTransactionIDTestSuite) TearDown() {
	randomRead = rand.Read
}

func (s *GenerateTransactionIDTestSuite) TestErrors() {
	// Error is returned from random number generator
	e := errors.New("mocked error")
	randomRead = randomReadMock(s.random, 0, e)
	_, err := GenerateTransactionID()
	s.Assert().Equal(e, err)

	// Less than 4 bytes are generated
	randomRead = randomReadMock(s.random, 2, nil)
	_, err = GenerateTransactionID()
	s.Assert().EqualError(err, "invalid random sequence: shorter than 3 bytes")
}

func (s *GenerateTransactionIDTestSuite) TestSuccess() {
	binary.BigEndian.PutUint32(s.random, 0x01020300)
	randomRead = randomReadMock(s.random, 3, nil)
	tid, err := GenerateTransactionID()
	s.Require().NoError(err)
	s.Assert().Equal(TransactionID{0x1, 0x2, 0x3}, tid)
}

func TestGenerateTransactionIDTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateTransactionIDTestSuite))
}

func TestNewMessage(t *testing.T) {
	d, err := NewMessage()
	require.NoError(t, err)
	require.NotNil(t, d)
	require.Equal(t, MessageTypeSolicit, d.Type())
	require.NotEqual(t, 0, d.TransactionID)
	require.Empty(t, d.Options)
}

func TestDecapsulateRelayIndex(t *testing.T) {
	m := Message{}
	r1, err := EncapsulateRelay(&m, MessageTypeRelayForward, net.IPv6linklocalallnodes, net.IPv6interfacelocalallnodes)
	require.NoError(t, err)
	r2, err := EncapsulateRelay(r1, MessageTypeRelayForward, net.IPv6loopback, net.IPv6linklocalallnodes)
	require.NoError(t, err)
	r3, err := EncapsulateRelay(r2, MessageTypeRelayForward, net.IPv6unspecified, net.IPv6linklocalallrouters)
	require.NoError(t, err)

	first, err := DecapsulateRelayIndex(r3, 0)
	require.NoError(t, err)
	relay, ok := first.(*RelayMessage)
	require.True(t, ok)
	require.Equal(t, relay.HopCount, uint8(1))
	require.Equal(t, relay.LinkAddr, net.IPv6loopback)
	require.Equal(t, relay.PeerAddr, net.IPv6linklocalallnodes)

	second, err := DecapsulateRelayIndex(r3, 1)
	require.NoError(t, err)
	relay, ok = second.(*RelayMessage)
	require.True(t, ok)
	require.Equal(t, relay.HopCount, uint8(0))
	require.Equal(t, relay.LinkAddr, net.IPv6linklocalallnodes)
	require.Equal(t, relay.PeerAddr, net.IPv6interfacelocalallnodes)

	third, err := DecapsulateRelayIndex(r3, 2)
	require.NoError(t, err)
	_, ok = third.(*Message)
	require.True(t, ok)

	rfirst, err := DecapsulateRelayIndex(r3, -1)
	require.NoError(t, err)
	relay, ok = rfirst.(*RelayMessage)
	require.True(t, ok)
	require.Equal(t, relay.HopCount, uint8(0))
	require.Equal(t, relay.LinkAddr, net.IPv6linklocalallnodes)
	require.Equal(t, relay.PeerAddr, net.IPv6interfacelocalallnodes)

	_, err = DecapsulateRelayIndex(r3, -2)
	require.Error(t, err)
}

func TestAddOption(t *testing.T) {
	d := Message{}
	require.Empty(t, d.Options)
	opt := OptionGeneric{OptionCode: 0, OptionData: []byte{}}
	d.AddOption(&opt)
	require.Equal(t, Options{&opt}, d.Options)
}

func TestToBytes(t *testing.T) {
	d := Message{
		MessageType:   MessageTypeSolicit,
		TransactionID: TransactionID{0xa, 0xb, 0xc},
	}
	d.AddOption(&OptionGeneric{OptionCode: 0, OptionData: []byte{}})

	bytes := d.ToBytes()
	expected := []byte{01, 0xa, 0xb, 0xc, 0x00, 0x00, 0x00, 0x00}
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
	s := Message{
		MessageType:   MessageTypeSolicit,
		TransactionID: TransactionID{0xa, 0xb, 0xc},
	}
	cid := OptClientId{}
	s.AddOption(&cid)
	duid := Duid{}

	a, err := NewAdvertiseFromSolicit(&s, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, a.TransactionID, s.TransactionID)
	require.Equal(t, a.Type(), MessageTypeAdvertise)
}

func TestNewReplyFromMessage(t *testing.T) {
	msg := Message{
		TransactionID: TransactionID{0xa, 0xb, 0xc},
		MessageType:   MessageTypeConfirm,
	}
	cid := OptClientId{}
	msg.AddOption(&cid)
	sid := OptServerId{}
	duid := Duid{}
	sid.Sid = duid
	msg.AddOption(&sid)

	rep, err := NewReplyFromMessage(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.TransactionID, msg.TransactionID)
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.MessageType = MessageTypeRenew
	rep, err = NewReplyFromMessage(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.TransactionID, msg.TransactionID)
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.MessageType = MessageTypeRebind
	rep, err = NewReplyFromMessage(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.TransactionID, msg.TransactionID)
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.MessageType = MessageTypeRelease
	rep, err = NewReplyFromMessage(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.TransactionID, msg.TransactionID)
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.MessageType = MessageTypeInformationRequest
	rep, err = NewReplyFromMessage(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.TransactionID, msg.TransactionID)
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.MessageType = MessageTypeSolicit
	rep, err = NewReplyFromMessage(&msg)
	require.Error(t, err)
}

func TestNewMessageTypeSolicitWithCID(t *testing.T) {
	hwAddr, err := net.ParseMAC("24:0A:9E:9F:EB:2B")
	require.NoError(t, err)

	duid := Duid{
		Type:          DUID_LL,
		HwType:        iana.HWTypeEthernet,
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
	msg := Message{}
	opt := OptClientArchType{ArchTypes: []iana.Arch{iana.EFI_BC}}
	msg.AddOption(&opt)
	require.True(t, IsUsingUEFI(&msg))
}

func TestIsUsingUEFIArchTypeFalse(t *testing.T) {
	msg := Message{}
	opt := OptClientArchType{ArchTypes: []iana.Arch{iana.INTEL_X86PC}}
	msg.AddOption(&opt)
	require.False(t, IsUsingUEFI(&msg))
}

func TestIsUsingUEFIUserClassTrue(t *testing.T) {
	msg := Message{}
	opt := OptUserClass{UserClasses: [][]byte{[]byte("ipxeUEFI")}}
	msg.AddOption(&opt)
	require.True(t, IsUsingUEFI(&msg))
}

func TestIsUsingUEFIUserClassFalse(t *testing.T) {
	msg := Message{}
	opt := OptUserClass{UserClasses: [][]byte{[]byte("ipxeLegacy")}}
	msg.AddOption(&opt)
	require.False(t, IsUsingUEFI(&msg))
}

func TestGetTransactionIDMessage(t *testing.T) {
	message, err := NewMessage()
	require.NoError(t, err)
	transactionID, err := GetTransactionID(message)
	require.NoError(t, err)
	require.Equal(t, transactionID, message.TransactionID)
}

func TestGetTransactionIDRelay(t *testing.T) {
	message, err := NewMessage()
	require.NoError(t, err)
	relay, err := EncapsulateRelay(message, MessageTypeRelayForward, nil, nil)
	require.NoError(t, err)
	transactionID, err := GetTransactionID(relay)
	require.NoError(t, err)
	require.Equal(t, transactionID, message.TransactionID)
}

// TODO test NewMessageTypeSolicit
//      test String and Summary
