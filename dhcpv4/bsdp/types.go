package bsdp

import "github.com/insomniacslk/dhcp/dhcpv4"

// DefaultMacOSVendorClassIdentifier is a default vendor class identifier used
// on non-darwin hosts where the vendor class identifier cannot be determined.
// It should mostly be used for debugging if testing BSDP on a non-darwin
// system.
const DefaultMacOSVendorClassIdentifier = "AAPLBSDP/i386/MacMini6,1"

// Options (occur as sub-options of DHCP option 43).
const (
	OptionMessageType                   dhcpv4.OptionCode = 1
	OptionVersion                       dhcpv4.OptionCode = 2
	OptionServerIdentifier              dhcpv4.OptionCode = 3
	OptionServerPriority                dhcpv4.OptionCode = 4
	OptionReplyPort                     dhcpv4.OptionCode = 5
	OptionBootImageListPath             dhcpv4.OptionCode = 6 // Not used
	OptionDefaultBootImageID            dhcpv4.OptionCode = 7
	OptionSelectedBootImageID           dhcpv4.OptionCode = 8
	OptionBootImageList                 dhcpv4.OptionCode = 9
	OptionNetboot1_0Firmware            dhcpv4.OptionCode = 10
	OptionBootImageAttributesFilterList dhcpv4.OptionCode = 11
	OptionShadowMountPath               dhcpv4.OptionCode = 128
	OptionShadowFilePath                dhcpv4.OptionCode = 129
	OptionMachineName                   dhcpv4.OptionCode = 130
)

// OptionCodeToString maps BSDP OptionCodes to human-readable strings
// describing what they are.
var OptionCodeToString = map[dhcpv4.OptionCode]string{
	OptionMessageType:                   "BSDP Message Type",
	OptionVersion:                       "BSDP Version",
	OptionServerIdentifier:              "BSDP Server Identifier",
	OptionServerPriority:                "BSDP Server Priority",
	OptionReplyPort:                     "BSDP Reply Port",
	OptionBootImageListPath:             "", // Not used
	OptionDefaultBootImageID:            "BSDP Default Boot Image ID",
	OptionSelectedBootImageID:           "BSDP Selected Boot Image ID",
	OptionBootImageList:                 "BSDP Boot Image List",
	OptionNetboot1_0Firmware:            "BSDP Netboot 1.0 Firmware",
	OptionBootImageAttributesFilterList: "BSDP Boot Image Attributes Filter List",
	OptionShadowMountPath:               "BSDP Shadow Mount Path",
	OptionShadowFilePath:                "BSDP Shadow File Path",
	OptionMachineName:                   "BSDP Machine Name",
}
