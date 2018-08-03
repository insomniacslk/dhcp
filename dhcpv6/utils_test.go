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
	optro.AddRequestedOption(OptionBootfileURL)
	msg2.AddOption(&optro)
	require.True(t, IsNetboot(&msg2))

	msg3 := DHCPv6Message{}
	optbf := OptBootFileURL{}
	msg3.AddOption(&optbf)
	require.True(t, IsNetboot(&msg3))
}

func TestIsRequetsed(t *testing.T) {
	msg1 := DHCPv6Message{}
	require.False(t, IsRequested(&msg1, OptionDNSRecursiveNameServer))

	msg2 := DHCPv6Message{}
	optro := OptRequestedOption{}
	optro.AddRequestedOption(OptionDNSRecursiveNameServer)
	msg2.AddOption(&optro)
	require.True(t, IsRequested(&msg2, OptionDNSRecursiveNameServer))
}

func TestIsUsingUEFIArchTypeTrue(t *testing.T) {
	msg := DHCPv6Message{}
	opt := OptClientArchType{ArchType: EFI_BC}
	msg.AddOption(&opt)
	require.True(t, IsUsingUEFI(&msg))
}

func TestIsUsingUEFIArchTypeFalse(t *testing.T) {
	msg := DHCPv6Message{}
	opt := OptClientArchType{ArchType: INTEL_X86PC}
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
