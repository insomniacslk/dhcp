package bsdp

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
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
	require.NoError(t, err)
	require.Empty(t, images, "no BootImages")
}

func TestNeedsReplyPort(t *testing.T) {
	require.True(t, needsReplyPort(123))
	require.False(t, needsReplyPort(0))
	require.False(t, needsReplyPort(dhcpv4.ClientPort))
}
