package bsdp

import (
	"errors"
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// MaxDHCPMessageSize is the size set in DHCP option 57 (DHCP Maximum Message Size).
// BSDP includes its own sub-option (12) to indicate to NetBoot servers that the
// client can support larger message sizes, and modern NetBoot servers will
// prefer this BSDP-specific option over the DHCP standard option.
const MaxDHCPMessageSize = 1500

// AppleVendorID is the string constant set in the vendor class identifier (DHCP
// option 60) that is sent by the server.
const AppleVendorID = "AAPLBSDPC"

// ReplyConfig is a struct containing some common configuration values for a
// BSDP reply (ACK).
type ReplyConfig struct {
	ServerIP                     net.IP
	ServerHostname, BootFileName string
	ServerPriority               int
	Images                       []BootImage
	DefaultImage, SelectedImage  *BootImage
}

// ParseBootImageListFromAck parses the list of boot images presented in the
// ACK[LIST] packet and returns them as a list of BootImages.
func ParseBootImageListFromAck(ack dhcpv4.DHCPv4) ([]BootImage, error) {
	var images []BootImage
	opt := ack.GetOneOption(dhcpv4.OptionVendorSpecificInformation)
	if opt == nil {
		return nil, errors.New("ParseBootImageListFromAck: could not find vendor-specific option")
	}
	vendorOpt, err := ParseOptVendorSpecificInformation(opt.ToBytes())
	if err != nil {
		return nil, err
	}
	bootImageOpts := vendorOpt.GetOption(OptionBootImageList)
	for _, opt := range bootImageOpts {
		images = append(images, opt.(*OptBootImageList).Images...)
	}
	return images, nil
}

func needsReplyPort(replyPort uint16) bool {
	return replyPort != 0 && replyPort != dhcpv4.ClientPort
}

func MessageTypeFromPacket(packet *dhcpv4.DHCPv4) (MessageType, bool) {
	var (
		messageType MessageType
		vendorOpts  *OptVendorSpecificInformation
		err         error
	)
	for _, opt := range packet.GetOption(dhcpv4.OptionVendorSpecificInformation) {
		if vendorOpts, err = ParseOptVendorSpecificInformation(opt.ToBytes()); err != nil {
			return messageType, false
		}
		if o := vendorOpts.GetOneOption(OptionMessageType); o != nil {
			if optMessageType, ok := o.(*OptMessageType); ok {
				return optMessageType.Type, true
			}
			return messageType, false
		}
	}
	return messageType, false
}

// NewInformListForInterface creates a new INFORM packet for interface ifname
// with configuration options specified by config.
func NewInformListForInterface(ifname string, replyPort uint16) (*dhcpv4.DHCPv4, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	// Get currently configured IP.
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	localIPs, err := dhcpv4.GetExternalIPv4Addrs(addrs)
	if err != nil {
		return nil, fmt.Errorf("could not get local IPv4 addr for %s: %v", iface.Name, err)
	}
	if localIPs == nil || len(localIPs) == 0 {
		return nil, fmt.Errorf("could not get local IPv4 addr for %s", iface.Name)
	}
	return NewInformList(iface.HardwareAddr, localIPs[0], replyPort)
}

// NewInformList creates a new INFORM packet for interface with hardware address
// `hwaddr` and IP `localIP`. Packet will be sent out on port `replyPort`.
func NewInformList(hwaddr net.HardwareAddr, localIP net.IP, replyPort uint16) (*dhcpv4.DHCPv4, error) {
	d, err := dhcpv4.NewInform(hwaddr, localIP)
	if err != nil {
		return nil, err
	}

	// Validate replyPort first
	if needsReplyPort(replyPort) && replyPort >= 1024 {
		return nil, errors.New("replyPort must be a privileged port")
	}

	// These are vendor-specific options used to pass along BSDP information.
	vendorOpts := []dhcpv4.Option{
		&OptMessageType{MessageTypeList},
		&OptVersion{Version1_1},
	}
	if needsReplyPort(replyPort) {
		vendorOpts = append(vendorOpts, &OptReplyPort{replyPort})
	}
	d.AddOption(&OptVendorSpecificInformation{vendorOpts})

	d.AddOption(&dhcpv4.OptParameterRequestList{
		RequestedOpts: []dhcpv4.OptionCode{
			dhcpv4.OptionVendorSpecificInformation,
			dhcpv4.OptionClassIdentifier,
		},
	})
	d.AddOption(&dhcpv4.OptMaximumDHCPMessageSize{Size: MaxDHCPMessageSize})

	vendorClassID, err := MakeVendorClassIdentifier()
	if err != nil {
		return nil, err
	}
	d.AddOption(&dhcpv4.OptClassIdentifier{Identifier: vendorClassID})
	return d, nil
}

// InformSelectForAck constructs an INFORM[SELECT] packet given an ACK to the
// previously-sent INFORM[LIST].
func InformSelectForAck(ack dhcpv4.DHCPv4, replyPort uint16, selectedImage BootImage) (*dhcpv4.DHCPv4, error) {
	d, err := dhcpv4.New()
	if err != nil {
		return nil, err
	}

	if needsReplyPort(replyPort) && replyPort >= 1024 {
		return nil, errors.New("replyPort must be a privileged port")
	}
	d.SetOpcode(dhcpv4.OpcodeBootRequest)
	d.SetHwType(ack.HwType())
	d.SetHwAddrLen(ack.HwAddrLen())
	clientHwAddr := ack.ClientHwAddr()
	d.SetClientHwAddr(clientHwAddr[:])
	d.SetTransactionID(ack.TransactionID())
	if ack.IsBroadcast() {
		d.SetBroadcast()
	} else {
		d.SetUnicast()
	}

	// Data for OptionSelectedBootImageID
	vendorOpts := []dhcpv4.Option{
		&OptMessageType{MessageTypeSelect},
		&OptVersion{Version1_1},
		&OptSelectedBootImageID{selectedImage.ID},
	}

	// Find server IP address
	var serverIP net.IP
	if opt := ack.GetOneOption(dhcpv4.OptionServerIdentifier); opt != nil {
		serverIP = opt.(*dhcpv4.OptServerIdentifier).ServerID
	}
	if serverIP.To4() == nil {
		return nil, fmt.Errorf("could not parse server identifier from ACK")
	}
	vendorOpts = append(vendorOpts, &OptServerIdentifier{serverIP})

	// Validate replyPort if requested.
	if needsReplyPort(replyPort) {
		vendorOpts = append(vendorOpts, &OptReplyPort{replyPort})
	}

	vendorClassID, err := MakeVendorClassIdentifier()
	if err != nil {
		return nil, err
	}
	d.AddOption(&dhcpv4.OptClassIdentifier{Identifier: vendorClassID})
	d.AddOption(&dhcpv4.OptParameterRequestList{
		RequestedOpts: []dhcpv4.OptionCode{
			dhcpv4.OptionSubnetMask,
			dhcpv4.OptionRouter,
			dhcpv4.OptionBootfileName,
			dhcpv4.OptionVendorSpecificInformation,
			dhcpv4.OptionClassIdentifier,
		},
	})
	d.AddOption(&dhcpv4.OptMessageType{MessageType: dhcpv4.MessageTypeInform})
	d.AddOption(&OptVendorSpecificInformation{vendorOpts})
	return d, nil
}

// NewReplyForInformList constructs an ACK for the INFORM[LIST] packet `inform`
// with additional options in `config`.
func NewReplyForInformList(inform *dhcpv4.DHCPv4, config ReplyConfig) (*dhcpv4.DHCPv4, error) {
	if config.DefaultImage == nil {
		return nil, errors.New("NewReplyForInformList: no default boot image ID set")
	}
	if config.Images == nil || len(config.Images) == 0 {
		return nil, errors.New("NewReplyForInformList: no boot images provided")
	}
	reply, err := dhcpv4.NewReplyFromRequest(inform)
	if err != nil {
		return nil, err
	}
	reply.SetClientIPAddr(inform.ClientIPAddr())
	reply.SetYourIPAddr(net.IPv4zero)
	reply.SetGatewayIPAddr(inform.GatewayIPAddr())
	reply.SetServerIPAddr(config.ServerIP)
	reply.SetServerHostName([]byte(config.ServerHostname))

	reply.AddOption(&dhcpv4.OptMessageType{MessageType: dhcpv4.MessageTypeAck})
	reply.AddOption(&dhcpv4.OptServerIdentifier{ServerID: config.ServerIP})
	reply.AddOption(&dhcpv4.OptClassIdentifier{Identifier: AppleVendorID})

	// BSDP opts.
	vendorOpts := []dhcpv4.Option{
		&OptMessageType{Type: MessageTypeList},
		&OptServerPriority{Priority: config.ServerPriority},
		&OptDefaultBootImageID{ID: config.DefaultImage.ID},
		&OptBootImageList{Images: config.Images},
	}
	if config.SelectedImage != nil {
		vendorOpts = append(vendorOpts, &OptSelectedBootImageID{ID: config.SelectedImage.ID})
	}
	reply.AddOption(&OptVendorSpecificInformation{Options: vendorOpts})

	reply.AddOption(&dhcpv4.OptionGeneric{OptionCode: dhcpv4.OptionEnd})
	return reply, nil
}

// NewReplyForInformSelect constructs an ACK for the INFORM[Select] packet
// `inform` with additional options in `config`.
func NewReplyForInformSelect(inform *dhcpv4.DHCPv4, config ReplyConfig) (*dhcpv4.DHCPv4, error) {
	if config.SelectedImage == nil {
		return nil, errors.New("NewReplyForInformSelect: no selected boot image ID set")
	}
	if config.Images == nil || len(config.Images) == 0 {
		return nil, errors.New("NewReplyForInformSelect: no boot images provided")
	}
	reply, err := dhcpv4.NewReplyFromRequest(inform)
	if err != nil {
		return nil, err
	}

	reply.SetClientIPAddr(inform.ClientIPAddr())
	reply.SetYourIPAddr(net.IPv4zero)
	reply.SetGatewayIPAddr(inform.GatewayIPAddr())
	reply.SetServerIPAddr(config.ServerIP)
	reply.SetServerHostName([]byte(config.ServerHostname))
	reply.SetBootFileName([]byte(config.BootFileName))

	reply.AddOption(&dhcpv4.OptMessageType{MessageType: dhcpv4.MessageTypeAck})
	reply.AddOption(&dhcpv4.OptServerIdentifier{ServerID: config.ServerIP})
	reply.AddOption(&dhcpv4.OptClassIdentifier{Identifier: AppleVendorID})

	// BSDP opts.
	reply.AddOption(&OptVendorSpecificInformation{
		Options: []dhcpv4.Option{
			&OptMessageType{Type: MessageTypeSelect},
			&OptSelectedBootImageID{ID: config.SelectedImage.ID},
		},
	})

	reply.AddOption(&dhcpv4.OptionGeneric{OptionCode: dhcpv4.OptionEnd})
	return reply, nil
}
