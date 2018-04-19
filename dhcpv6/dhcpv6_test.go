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

// TODO test NewSolicit
//      test String and Summary
