package dhcpv6

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

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
