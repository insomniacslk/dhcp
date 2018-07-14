package dhcpv6

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestBytesToTransactionID(t *testing.T) {
	// only the first three bytes should be used
	tid, err := BytesToTransactionID([]byte{0x11, 0x22, 0x33, 0xaa})
	require.NoError(t, err)
	require.Equal(t, uint32(0x112233), *tid)
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
	require.NotNil(t, *tid)
	require.True(t, *tid <= 0xffffff, "transaction ID should be smaller than 0xffffff")
}

func TestNewMessage(t *testing.T) {
	d, err := NewMessage()
	require.NoError(t, err)
	require.NotNil(t, d)
	require.Equal(t, SOLICIT, d.Type())
	require.NotEqual(t, 0, d.(*DHCPv6Message).transactionID)
	require.Empty(t, d.(*DHCPv6Message).options)
}

func TestDecapsulateRelayIndex(t *testing.T) {
	m := DHCPv6Message{}
	r1, err := EncapsulateRelay(&m, RELAY_FORW, net.IPv6linklocalallnodes, net.IPv6interfacelocalallnodes)
	require.NoError(t, err)
	r2, err := EncapsulateRelay(r1, RELAY_FORW, net.IPv6loopback, net.IPv6linklocalallnodes)
	require.NoError(t, err)
	r3, err := EncapsulateRelay(r2, RELAY_FORW, net.IPv6unspecified, net.IPv6linklocalallrouters)
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
	d.SetMessage(SOLICIT)
	require.Equal(t, SOLICIT, d.Type())
	d.SetMessage(ADVERTISE)
	require.Equal(t, ADVERTISE, d.Type())

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
	d.SetMessage(SOLICIT)
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
	s.SetMessage(SOLICIT)
	s.SetTransactionID(0xabcdef)
	cid := OptClientId{}
	s.AddOption(&cid)
	duid := Duid{}

	a, err := NewAdvertiseFromSolicit(&s, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, a.(*DHCPv6Message).TransactionID(), s.TransactionID())
	require.Equal(t, a.Type(), ADVERTISE)
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

	msg.SetMessage(CONFIRM)
	rep, err := NewReplyFromDHCPv6Message(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), msg.TransactionID())
	require.Equal(t, rep.Type(), REPLY)

	msg.SetMessage(RENEW)
	rep, err = NewReplyFromDHCPv6Message(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), msg.TransactionID())
	require.Equal(t, rep.Type(), REPLY)

	msg.SetMessage(REBIND)
	rep, err = NewReplyFromDHCPv6Message(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), msg.TransactionID())
	require.Equal(t, rep.Type(), REPLY)

	msg.SetMessage(RELEASE)
	rep, err = NewReplyFromDHCPv6Message(&msg, WithServerID(duid))
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), msg.TransactionID())
	require.Equal(t, rep.Type(), REPLY)

	msg.SetMessage(SOLICIT)
	rep, err = NewReplyFromDHCPv6Message(&msg)
	require.Error(t, err)

	relay := DHCPv6Relay{}
	rep, err = NewReplyFromDHCPv6Message(&relay)
	require.Error(t, err)
}

func TestNewSolicitWithCID(t *testing.T) {
	hwAddr, err := net.ParseMAC("24:0A:9E:9F:EB:2B")
	require.NoError(t, err)

	duid := Duid{
		Type:          DUID_LL,
		HwType:        iana.HwTypeEthernet,
		LinkLayerAddr: hwAddr,
	}

	s, err := NewSolicitWithCID(duid)
	require.NoError(t, err)

	require.Equal(t, s.Type(), SOLICIT)
	// Check CID
	cidOption := s.GetOneOption(OPTION_CLIENTID)
	require.NotNil(t, cidOption)
	cid, ok := cidOption.(*OptClientId)
	require.True(t, ok)
	require.Equal(t, cid.Cid, duid)

	// Check ORO
	oroOption := s.GetOneOption(OPTION_ORO)
	require.NotNil(t, oroOption)
	oro, ok := oroOption.(*OptRequestedOption)
	require.True(t, ok)
	opts := oro.RequestedOptions()
	require.Contains(t, opts, DNS_RECURSIVE_NAME_SERVER)
	require.Contains(t, opts, DOMAIN_SEARCH_LIST)
	require.Equal(t, len(opts), 2)
}

// TODO test NewSolicit
//      test String and Summary
