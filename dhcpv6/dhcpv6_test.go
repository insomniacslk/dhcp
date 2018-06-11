package dhcpv6

import (
	"testing"

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

func withServerID(d DHCPv6) DHCPv6 {
	sid := OptServerId{}
	d.AddOption(&sid)
	return d
}

func TestNewAdvertiseFromSolicit(t *testing.T) {
	s := DHCPv6Message{}
	s.SetMessage(SOLICIT)
	s.SetTransactionID(0xabcdef)
	cid := OptClientId{}
	s.AddOption(&cid)

	a, err := NewAdvertiseFromSolicit(&s, withServerID)
	require.NoError(t, err)
	require.Equal(t, a.(*DHCPv6Message).TransactionID(), s.TransactionID())
	require.Equal(t, a.Type(), ADVERTISE)
}

func TestNewReplyFromRequest(t *testing.T) {
	req := DHCPv6Message{}
	req.SetMessage(REQUEST)
	req.SetTransactionID(0xabcdef)
	cid := OptClientId{}
	req.AddOption(&cid)
	sid := OptServerId{}
	req.AddOption(&sid)

	rep, err := NewReplyFromRequest(&req)
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), req.TransactionID())
	require.Equal(t, rep.Type(), REPLY)
}

func TestNewReplyFromRenew(t *testing.T) {
	ren := DHCPv6Message{}
	ren.SetMessage(RENEW)
	ren.SetTransactionID(0xabcdef)
	cid := OptClientId{}
	ren.AddOption(&cid)

	rep, err := NewReplyFromRenew(&ren)
	require.Error(t, err)

	sid := OptServerId{}
	ren.AddOption(&sid)
	rep, err = NewReplyFromRenew(&ren)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), ren.TransactionID())
	require.Equal(t, rep.Type(), REPLY)
}

func TestNewReplyFromRebind(t *testing.T) {
	reb := DHCPv6Message{}
	reb.SetMessage(REPLY)
	rep, err := NewReplyFromRebind(&reb)
	require.Error(t, err)

	reb.SetMessage(REBIND)
	reb.SetTransactionID(0xabcdef)
	cid := OptClientId{}
	reb.AddOption(&cid)

	rep, err = NewReplyFromRebind(&reb)
	require.NoError(t, err)
	require.Equal(t, rep.(*DHCPv6Message).TransactionID(), reb.TransactionID())
	require.Equal(t, rep.Type(), REPLY)
}

// TODO test NewSolicit
//      test String and Summary
