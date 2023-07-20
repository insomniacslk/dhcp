package dhcpv6

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/u-root/uio/rand"
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
	require.Equal(t, Options{&opt}, d.Options.Options)
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
	expected := [][]byte{
		{01, 0xab, 0xcd, 0xef, 0x00, 0x00, 0x00, 0x00},
		[]byte("0000\x00\x01\x00\x0e\x00\x01000000000000"),
	}
	t.Parallel()
	for i, packet := range expected {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			d, err := FromBytes(packet)
			require.NoError(t, err)
			toBytes := d.ToBytes()
			require.Equal(t, packet, toBytes)
		})
	}
}

func TestFromBytesInvalid(t *testing.T) {
	expected := [][]byte{
		{},
		{30},
		{12},
	}
	t.Parallel()
	for i, packet := range expected {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := FromBytes(packet)
			require.Error(t, err)
		})
	}
}

func TestNewAdvertiseFromSolicit(t *testing.T) {
	s := Message{
		MessageType:   MessageTypeSolicit,
		TransactionID: TransactionID{0xa, 0xb, 0xc},
	}
	s.AddOption(OptClientID(&DUIDLLT{}))

	a, err := NewAdvertiseFromSolicit(&s, WithServerID(&DUIDLLT{}))
	require.NoError(t, err)
	require.Equal(t, a.TransactionID, s.TransactionID)
	require.Equal(t, a.Type(), MessageTypeAdvertise)
}

func TestNewReplyFromMessage(t *testing.T) {
	msg := Message{
		TransactionID: TransactionID{0xa, 0xb, 0xc},
		MessageType:   MessageTypeConfirm,
	}
	var duid DUIDLLT
	msg.AddOption(OptClientID(&duid))
	msg.AddOption(OptServerID(&duid))

	rep, err := NewReplyFromMessage(&msg, WithServerID(&duid))
	require.NoError(t, err)
	require.Equal(t, rep.TransactionID, msg.TransactionID)
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.MessageType = MessageTypeRenew
	rep, err = NewReplyFromMessage(&msg, WithServerID(&duid))
	require.NoError(t, err)
	require.Equal(t, rep.TransactionID, msg.TransactionID)
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.MessageType = MessageTypeRebind
	rep, err = NewReplyFromMessage(&msg, WithServerID(&duid))
	require.NoError(t, err)
	require.Equal(t, rep.TransactionID, msg.TransactionID)
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.MessageType = MessageTypeRelease
	rep, err = NewReplyFromMessage(&msg, WithServerID(&duid))
	require.NoError(t, err)
	require.Equal(t, rep.TransactionID, msg.TransactionID)
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.MessageType = MessageTypeInformationRequest
	rep, err = NewReplyFromMessage(&msg, WithServerID(&duid))
	require.NoError(t, err)
	require.Equal(t, rep.TransactionID, msg.TransactionID)
	require.Equal(t, rep.Type(), MessageTypeReply)

	msg.MessageType = MessageTypeSolicit
	_, err = NewReplyFromMessage(&msg)
	require.Error(t, err)

	msg.MessageType = MessageTypeSolicit
	msg.AddOption(&OptionGeneric{OptionCode: OptionRapidCommit})
	_, err = NewReplyFromMessage(&msg)
	require.NoError(t, err)
	msg.Options.Del(OptionRapidCommit)
}

func TestNewMessageTypeSolicit(t *testing.T) {
	hwAddr, err := net.ParseMAC("24:0A:9E:9F:EB:2B")
	require.NoError(t, err)

	duid := &DUIDLL{
		HWType:        iana.HWTypeEthernet,
		LinkLayerAddr: hwAddr,
	}

	s, err := NewSolicit(hwAddr, WithClientID(duid))
	require.NoError(t, err)

	require.Equal(t, s.Type(), MessageTypeSolicit)
	// Check CID
	cduid := s.Options.ClientID()
	require.NotNil(t, cduid)
	require.Equal(t, cduid, duid)

	// Check ORO
	oro := s.Options.RequestedOptions()
	require.Contains(t, oro, OptionDNSRecursiveNameServer)
	require.Contains(t, oro, OptionDomainSearchList)
	require.Equal(t, len(oro), 2)

	// Check IA_NA
	iaid := [4]byte{hwAddr[2], hwAddr[3], hwAddr[4], hwAddr[5]}
	iana := s.Options.OneIANA()
	require.NotNil(t, iana)
	require.Equal(t, iaid, iana.IaId)
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

func FuzzDHCPv6(f *testing.F) {

	var relayForwBytesDuidUUID_data = []byte{
		0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x26, 0x8a, 0x07, 0xff, 0xfe, 0x56,
		0xdc, 0xa4, 0x00, 0x12, 0x00, 0x06, 0x24, 0x8a,
		0x07, 0x56, 0xdc, 0xa4, 0x00, 0x09, 0x00, 0x5a,
		0x06, 0x7d, 0x9b, 0xca, 0x00, 0x01, 0x00, 0x12,
		0x00, 0x04, 0xb7, 0xfd, 0x0a, 0x8c, 0x1b, 0x14,
		0x10, 0xaa, 0xeb, 0x0a, 0x5b, 0x3f, 0xe8, 0x9d,
		0x0f, 0x56, 0x00, 0x06, 0x00, 0x0a, 0x00, 0x17,
		0x00, 0x18, 0x00, 0x17, 0x00, 0x18, 0x00, 0x01,
		0x00, 0x08, 0x00, 0x02, 0xff, 0xff, 0x00, 0x03,
		0x00, 0x28, 0x07, 0x56, 0xdc, 0xa4, 0x00, 0x00,
		0x0e, 0x10, 0x00, 0x00, 0x15, 0x18, 0x00, 0x05,
		0x00, 0x18, 0x26, 0x20, 0x01, 0x0d, 0xc0, 0x82,
		0x90, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xaf, 0xa0, 0x00, 0x00, 0x1c, 0x20, 0x00, 0x00,
		0x1d, 0x4c}

	f.Add(relayForwBytesDuidUUID_data)
	f.Add([]byte{01, 0xab, 0xcd, 0xef, 0x00, 0x00, 0x00, 0x00})
	f.Add([]byte{01, 0xa, 0xb, 0xc, 0x00, 0x00, 0x00, 0x00})
	f.Add([]byte("0000\x00\x01\x00\x0e\x00\x01000000000000"))

	f.Fuzz(func(t *testing.T, data []byte) {
		msg, err := FromBytes(data)
		if err != nil {
			return
		}
		msg.ToBytes()
	})
}
