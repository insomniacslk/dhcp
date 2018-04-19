package dhcpv6

import (
	"log"
)

// WithNetboot adds bootfile URL and bootfile param options to a DHCPv6 packet.
func WithNetboot(d DHCPv6) DHCPv6 {
	msg, ok := d.(*DHCPv6Message)
	if !ok {
		log.Printf("WithNetboot: not a DHCPv6Message")
		return d
	}
	// add OPT_BOOTFILE_URL and OPT_BOOTFILE_PARAM
	opt := msg.GetOneOption(OPTION_ORO)
	if opt == nil {
		opt = &OptRequestedOption{}
	}
	// TODO only add options if they are not there already
	oro := opt.(*OptRequestedOption)
	oro.AddRequestedOption(OPT_BOOTFILE_URL)
	oro.AddRequestedOption(OPT_BOOTFILE_PARAM)
	msg.UpdateOption(oro)
	return d
}

// WithUserClass adds a user class option to the packet
func WithUserClass(uc string) Modifier {
	return func(d DHCPv6) DHCPv6 {
		ouc := OptUserClass{UserClass: []byte("FbLoL")}
		d.AddOption(&ouc)
		return d
	}
}

// WithArchType adds an arch type option to the packet
func WithArchType(at ArchType) Modifier {
	return func(d DHCPv6) DHCPv6 {
		ao := OptClientArchType{ArchType: at}
		d.AddOption(&ao)
		return d
	}
}
