package bsdp

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseBootImageListFromAck(t *testing.T) {
	expectedBootImages := []BootImage{
		BootImage{
			ID: BootImageID{
				IsInstall: true,
				ImageType: BootImageTypeMacOSX,
				Index:     0x1010,
			},
			Name: "bsdp-1",
		},
		BootImage{
			ID: BootImageID{
				IsInstall: false,
				ImageType: BootImageTypeMacOS9,
				Index:     0x1111,
			},
			Name: "bsdp-2",
		},
	}
	ack, _ := dhcpv4.New()
	ack.AddOption(&OptVendorSpecificInformation{
		[]dhcpv4.Option{&OptBootImageList{expectedBootImages}},
	})

	images, err := ParseBootImageListFromAck(*ack)
	require.NoError(t, err)
	require.NotEmpty(t, images, "should get BootImages")
	require.Equal(t, expectedBootImages, images, "should get same BootImages")
}

func TestParseBootImageListFromAckNoVendorOption(t *testing.T) {
	ack, _ := dhcpv4.New()
	images, err := ParseBootImageListFromAck(*ack)
	require.Error(t, err)
	require.Empty(t, images, "no BootImages")
}

func TestNeedsReplyPort(t *testing.T) {
	require.True(t, needsReplyPort(123))
	require.False(t, needsReplyPort(0))
	require.False(t, needsReplyPort(dhcpv4.ClientPort))
}

func TestNewInformList_NoReplyPort(t *testing.T) {
	hwAddr := net.HardwareAddr{1, 2, 3, 4, 5, 6}
	localIP := net.IPv4(10, 10, 11, 11)
	m, err := NewInformList(hwAddr, localIP, 0)

	require.NoError(t, err)
	require.True(t, dhcpv4.HasOption(m, dhcpv4.OptionVendorSpecificInformation))
	require.True(t, dhcpv4.HasOption(m, dhcpv4.OptionParameterRequestList))
	require.True(t, dhcpv4.HasOption(m, dhcpv4.OptionMaximumDHCPMessageSize))
	require.True(t, dhcpv4.HasOption(m, dhcpv4.OptionEnd))

	opt := m.GetOneOption(dhcpv4.OptionVendorSpecificInformation)
	require.NotNil(t, opt, "vendor opts not present")
	vendorInfo := opt.(*OptVendorSpecificInformation)
	require.True(t, dhcpv4.HasOption(vendorInfo, OptionMessageType))
	require.True(t, dhcpv4.HasOption(vendorInfo, OptionVersion))

	opt = vendorInfo.GetOneOption(OptionMessageType)
	require.Equal(t, MessageTypeList, opt.(*OptMessageType).Type)
}

func TestNewInformList_ReplyPort(t *testing.T) {
	hwAddr := net.HardwareAddr{1, 2, 3, 4, 5, 6}
	localIP := net.IPv4(10, 10, 11, 11)
	replyPort := uint16(11223)

	// Bad reply port
	_, err := NewInformList(hwAddr, localIP, replyPort)
	require.Error(t, err)

	// Good reply port
	replyPort = uint16(999)
	m, err := NewInformList(hwAddr, localIP, replyPort)
	require.NoError(t, err)

	opt := m.GetOneOption(dhcpv4.OptionVendorSpecificInformation)
	vendorInfo := opt.(*OptVendorSpecificInformation)
	require.True(t, dhcpv4.HasOption(vendorInfo, OptionReplyPort))

	opt = vendorInfo.GetOneOption(OptionReplyPort)
	require.Equal(t, replyPort, opt.(*OptReplyPort).Port)
}

func newAck(hwAddr []byte, transactionID uint32) *dhcpv4.DHCPv4 {
	ack, _ := dhcpv4.New()
	ack.SetTransactionID(transactionID)
	ack.SetHwType(iana.HwTypeEthernet)
	ack.SetClientHwAddr(hwAddr)
	ack.SetHwAddrLen(uint8(len(hwAddr)))
	ack.AddOption(&dhcpv4.OptMessageType{MessageType: dhcpv4.MessageTypeAck})
	ack.AddOption(&dhcpv4.OptionGeneric{OptionCode: dhcpv4.OptionEnd})
	return ack
}

func TestInformSelectForAck_Broadcast(t *testing.T) {
	hwAddr := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	tid := uint32(22)
	serverID := net.IPv4(1, 2, 3, 4)
	bootImage := BootImage{
		ID: BootImageID{
			IsInstall: true,
			ImageType: BootImageTypeMacOSX,
			Index:     0x1000,
		},
		Name: "bsdp-1",
	}
	ack := newAck(hwAddr, tid)
	ack.SetBroadcast()
	ack.AddOption(&dhcpv4.OptServerIdentifier{ServerID: serverID})

	m, err := InformSelectForAck(*ack, 0, bootImage)
	require.NoError(t, err)
	require.Equal(t, dhcpv4.OpcodeBootRequest, m.Opcode())
	require.Equal(t, ack.HwType(), m.HwType())
	require.Equal(t, ack.ClientHwAddr(), m.ClientHwAddr())
	require.Equal(t, ack.TransactionID(), m.TransactionID())
	require.True(t, m.IsBroadcast())

	// Validate options.
	require.True(t, dhcpv4.HasOption(m, dhcpv4.OptionClassIdentifier))
	require.True(t, dhcpv4.HasOption(m, dhcpv4.OptionParameterRequestList))
	require.True(t, dhcpv4.HasOption(m, dhcpv4.OptionDHCPMessageType))
	opt := m.GetOneOption(dhcpv4.OptionDHCPMessageType)
	require.Equal(t, dhcpv4.MessageTypeInform, opt.(*dhcpv4.OptMessageType).MessageType)
	require.True(t, dhcpv4.HasOption(m, dhcpv4.OptionEnd))

	// Validate vendor opts.
	require.True(t, dhcpv4.HasOption(m, dhcpv4.OptionVendorSpecificInformation))
	opt = m.GetOneOption(dhcpv4.OptionVendorSpecificInformation)
	vendorInfo := opt.(*OptVendorSpecificInformation)
	require.True(t, dhcpv4.HasOption(vendorInfo, OptionMessageType))
	opt = vendorInfo.GetOneOption(OptionMessageType)
	require.Equal(t, MessageTypeSelect, opt.(*OptMessageType).Type)
	require.True(t, dhcpv4.HasOption(vendorInfo, OptionVersion))
	require.True(t, dhcpv4.HasOption(vendorInfo, OptionSelectedBootImageID))
	opt = vendorInfo.GetOneOption(OptionSelectedBootImageID)
	require.Equal(t, bootImage.ID, opt.(*OptSelectedBootImageID).ID)
	require.True(t, dhcpv4.HasOption(vendorInfo, OptionServerIdentifier))
	opt = vendorInfo.GetOneOption(OptionServerIdentifier)
	require.True(t, serverID.Equal(opt.(*OptServerIdentifier).ServerID))
}

func TestInformSelectForAck_NoServerID(t *testing.T) {
	hwAddr := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	tid := uint32(22)
	bootImage := BootImage{
		ID: BootImageID{
			IsInstall: true,
			ImageType: BootImageTypeMacOSX,
			Index:     0x1000,
		},
		Name: "bsdp-1",
	}
	ack := newAck(hwAddr, tid)

	_, err := InformSelectForAck(*ack, 0, bootImage)
	require.Error(t, err, "expect error for no server identifier option")
}

func TestInformSelectForAck_BadReplyPort(t *testing.T) {
	hwAddr := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	tid := uint32(22)
	serverID := net.IPv4(1, 2, 3, 4)
	bootImage := BootImage{
		ID: BootImageID{
			IsInstall: true,
			ImageType: BootImageTypeMacOSX,
			Index:     0x1000,
		},
		Name: "bsdp-1",
	}
	ack := newAck(hwAddr, tid)
	ack.SetBroadcast()
	ack.AddOption(&dhcpv4.OptServerIdentifier{ServerID: serverID})

	_, err := InformSelectForAck(*ack, 11223, bootImage)
	require.Error(t, err, "expect error for > 1024 replyPort")
}

func TestInformSelectForAck_ReplyPort(t *testing.T) {
	hwAddr := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	tid := uint32(22)
	serverID := net.IPv4(1, 2, 3, 4)
	bootImage := BootImage{
		ID: BootImageID{
			IsInstall: true,
			ImageType: BootImageTypeMacOSX,
			Index:     0x1000,
		},
		Name: "bsdp-1",
	}
	ack := newAck(hwAddr, tid)
	ack.SetBroadcast()
	ack.AddOption(&dhcpv4.OptServerIdentifier{ServerID: serverID})

	replyPort := uint16(999)
	m, err := InformSelectForAck(*ack, replyPort, bootImage)
	require.NoError(t, err)

	require.True(t, dhcpv4.HasOption(m, dhcpv4.OptionVendorSpecificInformation))
	opt := m.GetOneOption(dhcpv4.OptionVendorSpecificInformation)
	vendorInfo := opt.(*OptVendorSpecificInformation)
	require.True(t, dhcpv4.HasOption(vendorInfo, OptionReplyPort))
	opt = vendorInfo.GetOneOption(OptionReplyPort)
	require.Equal(t, replyPort, opt.(*OptReplyPort).Port)
}
