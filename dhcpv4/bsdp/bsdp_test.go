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

// TODO(get9): Remove when #99 lands.
func newInform() *dhcpv4.DHCPv4 {
	p, _ := dhcpv4.New()
	p.SetClientIPAddr(net.IP{1, 2, 3, 4})
	p.SetGatewayIPAddr(net.IP{4, 3, 2, 1})
	p.SetHwType(iana.HwTypeEthernet)
	hwAddr := [16]byte{1, 2, 3, 4, 5, 6}
	p.SetClientHwAddr(hwAddr[:])
	p.SetHwAddrLen(6)
	return p
}

func TestNewReplyForInformList_NoDefaultImage(t *testing.T) {
	inform := newInform()
	_, err := NewReplyForInformList(inform, ReplyConfig{})
	require.Error(t, err)
}

func TestNewReplyForInformList_NoImages(t *testing.T) {
	inform := newInform()
	fakeImage := BootImage{
		ID: BootImageID{ImageType: BootImageTypeMacOSX},
	}
	_, err := NewReplyForInformList(inform, ReplyConfig{
		Images:       []BootImage{},
		DefaultImage: &fakeImage,
	})
	require.Error(t, err)
}

// TODO (get9): clean up when #99 lands.
func TestNewReplyForInformList(t *testing.T) {
	inform := newInform()
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
	require.Equal(t, net.IP{1, 2, 3, 4}, ack.ClientIPAddr())
	require.Equal(t, net.IPv4zero, ack.YourIPAddr())
	require.Equal(t, net.IP{4, 3, 2, 1}, ack.GatewayIPAddr())
	require.Equal(t, "bsdp.foo.com", ack.ServerHostNameToString())

	// Validate options.
	require.Equal(
		t,
		&dhcpv4.OptMessageType{MessageType: dhcpv4.MessageTypeAck},
		ack.GetOneOption(dhcpv4.OptionDHCPMessageType).(*dhcpv4.OptMessageType),
	)
	require.Equal(
		t,
		&dhcpv4.OptServerIdentifier{ServerID: net.IP{9, 9, 9, 9}},
		ack.GetOneOption(dhcpv4.OptionServerIdentifier).(*dhcpv4.OptServerIdentifier),
	)
	require.Equal(
		t,
		&dhcpv4.OptClassIdentifier{Identifier: AppleVendorID},
		ack.GetOneOption(dhcpv4.OptionClassIdentifier).(*dhcpv4.OptClassIdentifier),
	)
	require.NotNil(t, ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation))
	require.Equal(t, &dhcpv4.OptionGeneric{OptionCode: dhcpv4.OptionEnd}, ack.Options()[len(ack.Options())-1])

	vendorOpts := ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation).(*OptVendorSpecificInformation)
	require.Equal(
		t,
		&OptMessageType{Type: MessageTypeList},
		vendorOpts.GetOneOption(OptionMessageType).(*OptMessageType),
	)
	require.Equal(
		t,
		&OptServerPriority{Priority: 0x7070},
		vendorOpts.GetOneOption(OptionServerPriority).(*OptServerPriority),
	)
	require.Equal(
		t,
		&OptDefaultBootImageID{ID: images[0].ID},
		vendorOpts.GetOneOption(OptionDefaultBootImageID).(*OptDefaultBootImageID),
	)
	require.Equal(
		t,
		&OptBootImageList{Images: images},
		vendorOpts.GetOneOption(OptionBootImageList).(*OptBootImageList),
	)

	// Add in selected boot image, ensure it's in the generated ACK.
	config.SelectedImage = &images[0]
	ack, err = NewReplyForInformList(inform, config)
	require.NoError(t, err)
	vendorOpts = ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation).(*OptVendorSpecificInformation)
	require.Equal(
		t,
		&OptSelectedBootImageID{ID: images[0].ID},
		vendorOpts.GetOneOption(OptionSelectedBootImageID).(*OptSelectedBootImageID),
	)
}

func TestNewReplyForInformSelect_NoSelectedImage(t *testing.T) {
	inform := newInform()
	_, err := NewReplyForInformSelect(inform, ReplyConfig{})
	require.Error(t, err)
}

func TestNewReplyForInformSelect_NoImages(t *testing.T) {
	inform := newInform()
	fakeImage := BootImage{
		ID: BootImageID{ImageType: BootImageTypeMacOSX},
	}
	_, err := NewReplyForInformSelect(inform, ReplyConfig{
		Images:        []BootImage{},
		SelectedImage: &fakeImage,
	})
	require.Error(t, err)
}

func TestNewReplyForInformSelect(t *testing.T) {
	inform := newInform()
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
	require.Equal(t, net.IP{1, 2, 3, 4}, ack.ClientIPAddr())
	require.Equal(t, net.IPv4zero, ack.YourIPAddr())
	require.Equal(t, net.IP{4, 3, 2, 1}, ack.GatewayIPAddr())
	require.Equal(t, "bsdp.foo.com", ack.ServerHostNameToString())

	// Validate options.
	require.Equal(
		t,
		&dhcpv4.OptMessageType{MessageType: dhcpv4.MessageTypeAck},
		ack.GetOneOption(dhcpv4.OptionDHCPMessageType).(*dhcpv4.OptMessageType),
	)
	require.Equal(
		t,
		&dhcpv4.OptServerIdentifier{ServerID: net.IP{9, 9, 9, 9}},
		ack.GetOneOption(dhcpv4.OptionServerIdentifier).(*dhcpv4.OptServerIdentifier),
	)
	require.Equal(
		t,
		&dhcpv4.OptClassIdentifier{Identifier: AppleVendorID},
		ack.GetOneOption(dhcpv4.OptionClassIdentifier).(*dhcpv4.OptClassIdentifier),
	)
	require.NotNil(t, ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation))
	require.Equal(t, &dhcpv4.OptionGeneric{OptionCode: dhcpv4.OptionEnd}, ack.Options()[len(ack.Options())-1])

	vendorOpts := ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation).(*OptVendorSpecificInformation)
	require.Equal(
		t,
		&OptMessageType{Type: MessageTypeSelect},
		vendorOpts.GetOneOption(OptionMessageType).(*OptMessageType),
	)
	require.Equal(
		t,
		&OptSelectedBootImageID{ID: images[0].ID},
		vendorOpts.GetOneOption(OptionSelectedBootImageID).(*OptSelectedBootImageID),
	)
}
