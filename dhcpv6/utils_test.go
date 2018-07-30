package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNetboot(t *testing.T) {
	msg1 := DHCPv6Message{}
	require.False(t, IsNetboot(&msg1))

	msg2 := DHCPv6Message{}
	optro := OptRequestedOption{}
	optro.AddRequestedOption(OPT_BOOTFILE_URL)
	msg2.AddOption(&optro)
	require.True(t, IsNetboot(&msg2))

	msg3 := DHCPv6Message{}
	optbf := OptBootFileURL{}
	msg3.AddOption(&optbf)
	require.True(t, IsNetboot(&msg3))
}

func TestIsUsingUEFIArchTypeTrue(t *testing.T) {
	msg := DHCPv6Message{}
	opt := OptClientArchType{ArchType: EFIBC}
	msg.AddOption(&opt)
	require.True(t, IsUsingUEFI(&msg))
}

func TestIsUsingUEFIArchTypeFalse(t *testing.T) {
	msg := DHCPv6Message{}
	opt := OptClientArchType{ArchType: Intelx86PC}
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
	relay, err := EncapsulateRelay(message, RELAY_FORW, nil, nil)
	require.NoError(t, err)
	transactionID, err := GetTransactionID(relay)
	require.NoError(t, err)
	require.Equal(t, transactionID, message.(*DHCPv6Message).TransactionID())
}
