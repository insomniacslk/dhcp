package bsdp

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func RequireHasOption(t *testing.T, opts dhcpv4.Options, opt dhcpv4.Option) {
	require.NotNil(t, opts, "must pass list of options")
	require.NotNil(t, opt, "must pass option")
	require.True(t, opts.Has(opt.Code()))
	actual := opts.GetOne(opt.Code())
	require.Equal(t, opt, actual)
}

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
	ack.UpdateOption(&OptVendorSpecificInformation{
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
	require.True(t, m.Options.Has(dhcpv4.OptionVendorSpecificInformation))
	require.True(t, m.Options.Has(dhcpv4.OptionParameterRequestList))
	require.True(t, m.Options.Has(dhcpv4.OptionMaximumDHCPMessageSize))

	opt := m.GetOneOption(dhcpv4.OptionVendorSpecificInformation)
	require.NotNil(t, opt, "vendor opts not present")
	vendorInfo := opt.(*OptVendorSpecificInformation)
	require.True(t, vendorInfo.Options.Has(OptionMessageType))
	require.True(t, vendorInfo.Options.Has(OptionVersion))

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
	require.True(t, vendorInfo.Options.Has(OptionReplyPort))

	opt = vendorInfo.GetOneOption(OptionReplyPort)
	require.Equal(t, replyPort, opt.(*OptReplyPort).Port)
}

func newAck(hwAddr net.HardwareAddr, transactionID [4]byte) *dhcpv4.DHCPv4 {
	ack, _ := dhcpv4.New()
	ack.OpCode = dhcpv4.OpcodeBootReply
	ack.TransactionID = transactionID
	ack.HWType = iana.HWTypeEthernet
	ack.ClientHWAddr = hwAddr
	ack.UpdateOption(&dhcpv4.OptMessageType{MessageType: dhcpv4.MessageTypeAck})
	return ack
}

func TestInformSelectForAck_Broadcast(t *testing.T) {
	hwAddr := net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	tid := [4]byte{0x22, 0, 0, 0}
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
	ack.UpdateOption(&dhcpv4.OptServerIdentifier{ServerID: serverID})

	m, err := InformSelectForAck(*ack, 0, bootImage)
	require.NoError(t, err)
	require.Equal(t, dhcpv4.OpcodeBootRequest, m.OpCode)
	require.Equal(t, ack.HWType, m.HWType)
	require.Equal(t, ack.ClientHWAddr, m.ClientHWAddr)
	require.Equal(t, ack.TransactionID, m.TransactionID)
	require.True(t, m.IsBroadcast())

	// Validate options.
	require.True(t, m.Options.Has(dhcpv4.OptionClassIdentifier))
	require.True(t, m.Options.Has(dhcpv4.OptionParameterRequestList))
	require.True(t, m.Options.Has(dhcpv4.OptionDHCPMessageType))
	opt := m.GetOneOption(dhcpv4.OptionDHCPMessageType)
	require.Equal(t, dhcpv4.MessageTypeInform, opt.(*dhcpv4.OptMessageType).MessageType)

	// Validate vendor opts.
	require.True(t, m.Options.Has(dhcpv4.OptionVendorSpecificInformation))
	opt = m.GetOneOption(dhcpv4.OptionVendorSpecificInformation)
	vendorInfo := opt.(*OptVendorSpecificInformation)
	RequireHasOption(t, vendorInfo.Options, &OptMessageType{Type: MessageTypeSelect})
	require.True(t, vendorInfo.Options.Has(OptionVersion))
	RequireHasOption(t, vendorInfo.Options, &OptSelectedBootImageID{ID: bootImage.ID})
	RequireHasOption(t, vendorInfo.Options, &OptServerIdentifier{ServerID: serverID})
}

func TestInformSelectForAck_NoServerID(t *testing.T) {
	hwAddr := net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	tid := [4]byte{0x22, 0, 0, 0}
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
	hwAddr := net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	tid := [4]byte{0x22, 0, 0, 0}
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
	ack.UpdateOption(&dhcpv4.OptServerIdentifier{ServerID: serverID})

	_, err := InformSelectForAck(*ack, 11223, bootImage)
	require.Error(t, err, "expect error for > 1024 replyPort")
}

func TestInformSelectForAck_ReplyPort(t *testing.T) {
	hwAddr := net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	tid := [4]byte{0x22, 0, 0, 0}
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
	ack.UpdateOption(&dhcpv4.OptServerIdentifier{ServerID: serverID})

	replyPort := uint16(999)
	m, err := InformSelectForAck(*ack, replyPort, bootImage)
	require.NoError(t, err)

	require.True(t, m.Options.Has(dhcpv4.OptionVendorSpecificInformation))
	opt := m.GetOneOption(dhcpv4.OptionVendorSpecificInformation)
	vendorInfo := opt.(*OptVendorSpecificInformation)
	RequireHasOption(t, vendorInfo.Options, &OptReplyPort{Port: replyPort})
}

func TestNewReplyForInformList_NoDefaultImage(t *testing.T) {
	inform, _ := NewInformList(net.HardwareAddr{1, 2, 3, 4, 5, 6}, net.IP{1, 2, 3, 4}, dhcpv4.ClientPort)
	_, err := NewReplyForInformList(inform, ReplyConfig{})
	require.Error(t, err)
}

func TestNewReplyForInformList_NoImages(t *testing.T) {
	inform, _ := NewInformList(net.HardwareAddr{1, 2, 3, 4, 5, 6}, net.IP{1, 2, 3, 4}, dhcpv4.ClientPort)
	fakeImage := BootImage{
		ID: BootImageID{ImageType: BootImageTypeMacOSX},
	}
	_, err := NewReplyForInformList(inform, ReplyConfig{
		Images:       []BootImage{},
		DefaultImage: &fakeImage,
	})
	require.Error(t, err)

	_, err = NewReplyForInformList(inform, ReplyConfig{
		Images:        nil,
		SelectedImage: &fakeImage,
	})
	require.Error(t, err)
}

func TestNewReplyForInformList(t *testing.T) {
	inform, _ := NewInformList(net.HardwareAddr{1, 2, 3, 4, 5, 6}, net.IP{1, 2, 3, 4}, dhcpv4.ClientPort)
	images := []BootImage{
		BootImage{
			ID: BootImageID{
				IsInstall: true,
				ImageType: BootImageTypeMacOSX,
				Index:     0x7070,
			},
			Name: "image-1",
		},
		BootImage{
			ID: BootImageID{
				IsInstall: true,
				ImageType: BootImageTypeMacOSX,
				Index:     0x8080,
			},
			Name: "image-2",
		},
	}
	config := ReplyConfig{
		Images:         images,
		DefaultImage:   &images[0],
		ServerIP:       net.IP{9, 9, 9, 9},
		ServerHostname: "bsdp.foo.com",
		ServerPriority: 0x7070,
	}
	ack, err := NewReplyForInformList(inform, config)
	require.NoError(t, err)
	require.Equal(t, net.IP{1, 2, 3, 4}, ack.ClientIPAddr)
	require.Equal(t, net.IPv4zero, ack.YourIPAddr)
	require.Equal(t, "bsdp.foo.com", ack.ServerHostName)

	// Validate options.
	RequireHasOption(t, ack.Options, &dhcpv4.OptMessageType{MessageType: dhcpv4.MessageTypeAck})
	RequireHasOption(t, ack.Options, &dhcpv4.OptServerIdentifier{ServerID: net.IP{9, 9, 9, 9}})
	RequireHasOption(t, ack.Options, &dhcpv4.OptClassIdentifier{Identifier: AppleVendorID})
	require.NotNil(t, ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation))

	// Vendor-specific options.
	vendorOpts := ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation).(*OptVendorSpecificInformation)
	RequireHasOption(t, vendorOpts.Options, &OptMessageType{Type: MessageTypeList})
	RequireHasOption(t, vendorOpts.Options, &OptDefaultBootImageID{ID: images[0].ID})
	RequireHasOption(t, vendorOpts.Options, &OptServerPriority{Priority: 0x7070})
	RequireHasOption(t, vendorOpts.Options, &OptBootImageList{Images: images})

	// Add in selected boot image, ensure it's in the generated ACK.
	config.SelectedImage = &images[0]
	ack, err = NewReplyForInformList(inform, config)
	require.NoError(t, err)
	vendorOpts = ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation).(*OptVendorSpecificInformation)
	RequireHasOption(t, vendorOpts.Options, &OptSelectedBootImageID{ID: images[0].ID})
}

func TestNewReplyForInformSelect_NoSelectedImage(t *testing.T) {
	inform, _ := NewInformList(net.HardwareAddr{1, 2, 3, 4, 5, 6}, net.IP{1, 2, 3, 4}, dhcpv4.ClientPort)
	_, err := NewReplyForInformSelect(inform, ReplyConfig{})
	require.Error(t, err)
}

func TestNewReplyForInformSelect_NoImages(t *testing.T) {
	inform, _ := NewInformList(net.HardwareAddr{1, 2, 3, 4, 5, 6}, net.IP{1, 2, 3, 4}, dhcpv4.ClientPort)
	fakeImage := BootImage{
		ID: BootImageID{ImageType: BootImageTypeMacOSX},
	}
	_, err := NewReplyForInformSelect(inform, ReplyConfig{
		Images:        []BootImage{},
		SelectedImage: &fakeImage,
	})
	require.Error(t, err)

	_, err = NewReplyForInformSelect(inform, ReplyConfig{
		Images:        nil,
		SelectedImage: &fakeImage,
	})
	require.Error(t, err)
}

func TestNewReplyForInformSelect(t *testing.T) {
	inform, _ := NewInformList(net.HardwareAddr{1, 2, 3, 4, 5, 6}, net.IP{1, 2, 3, 4}, dhcpv4.ClientPort)
	images := []BootImage{
		BootImage{
			ID: BootImageID{
				IsInstall: true,
				ImageType: BootImageTypeMacOSX,
				Index:     0x7070,
			},
			Name: "image-1",
		},
		BootImage{
			ID: BootImageID{
				IsInstall: true,
				ImageType: BootImageTypeMacOSX,
				Index:     0x8080,
			},
			Name: "image-2",
		},
	}
	config := ReplyConfig{
		Images:         images,
		SelectedImage:  &images[0],
		ServerIP:       net.IP{9, 9, 9, 9},
		ServerHostname: "bsdp.foo.com",
		ServerPriority: 0x7070,
	}
	ack, err := NewReplyForInformSelect(inform, config)
	require.NoError(t, err)
	require.Equal(t, net.IP{1, 2, 3, 4}, ack.ClientIPAddr)
	require.Equal(t, net.IPv4zero, ack.YourIPAddr)
	require.Equal(t, "bsdp.foo.com", ack.ServerHostName)

	// Validate options.
	RequireHasOption(t, ack.Options, &dhcpv4.OptMessageType{MessageType: dhcpv4.MessageTypeAck})
	RequireHasOption(t, ack.Options, &dhcpv4.OptServerIdentifier{ServerID: net.IP{9, 9, 9, 9}})
	RequireHasOption(t, ack.Options, &dhcpv4.OptServerIdentifier{ServerID: net.IP{9, 9, 9, 9}})
	RequireHasOption(t, ack.Options, &dhcpv4.OptClassIdentifier{Identifier: AppleVendorID})
	require.NotNil(t, ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation))

	vendorOpts := ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation).(*OptVendorSpecificInformation)
	RequireHasOption(t, vendorOpts.Options, &OptMessageType{Type: MessageTypeSelect})
	RequireHasOption(t, vendorOpts.Options, &OptSelectedBootImageID{ID: images[0].ID})
}

func TestMessageTypeForPacket(t *testing.T) {
	var (
		pkt            *dhcpv4.DHCPv4
		gotMessageType *MessageType
	)

	list := new(MessageType)
	*list = MessageTypeList

	testcases := []struct {
		tcName          string
		opts            []dhcpv4.Option
		wantMessageType *MessageType
	}{
		{
			tcName: "No options",
			opts:   []dhcpv4.Option{},
		},
		{
			tcName: "Some options, no vendor opts",
			opts: []dhcpv4.Option{
				&dhcpv4.OptHostName{HostName: "foobar1234"},
			},
		},
		{
			tcName: "Vendor opts, no message type",
			opts: []dhcpv4.Option{
				&dhcpv4.OptHostName{HostName: "foobar1234"},
				&OptVendorSpecificInformation{
					Options: []dhcpv4.Option{
						Version1_1,
					},
				},
			},
		},
		{
			tcName: "Vendor opts, with message type",
			opts: []dhcpv4.Option{
				&dhcpv4.OptHostName{HostName: "foobar1234"},
				&OptVendorSpecificInformation{
					Options: []dhcpv4.Option{
						Version1_1,
						&OptMessageType{Type: MessageTypeList},
					},
				},
			},
			wantMessageType: list,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.tcName, func(t *testing.T) {
			pkt, _ = dhcpv4.New()
			for _, opt := range tt.opts {
				pkt.UpdateOption(opt)
			}
			gotMessageType = MessageTypeFromPacket(pkt)
			require.Equal(t, tt.wantMessageType, gotMessageType)
			if tt.wantMessageType != nil {
				require.Equal(t, *tt.wantMessageType, *gotMessageType)
			}
		})
	}
}
