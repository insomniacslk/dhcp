package dhcpv6

import (
	"log"
)

// WithClientID adds a client ID option to a DHCPv6 packet
func WithClientID(duid Duid) Modifier {
	return func(d DHCPv6) DHCPv6 {
		cid := OptClientId{Cid: duid}
		d.UpdateOption(&cid)
		return d
	}
}

// WithServerID adds a client ID option to a DHCPv6 packet
func WithServerID(duid Duid) Modifier {
	return func(d DHCPv6) DHCPv6 {
		sid := OptServerId{Sid: duid}
		d.UpdateOption(&sid)
		return d
	}
}

// WithNetboot adds bootfile URL and bootfile param options to a DHCPv6 packet.
func WithNetboot(d DHCPv6) DHCPv6 {
	msg, ok := d.(*DHCPv6Message)
	if !ok {
		log.Printf("WithNetboot: not a DHCPv6Message")
		return d
	}
	// add OptionBootfileURL and OptionBootfileParam
	opt := msg.GetOneOption(OptionORO)
	if opt == nil {
		opt = &OptRequestedOption{}
	}
	// TODO only add options if they are not there already
	oro := opt.(*OptRequestedOption)
	oro.AddRequestedOption(OptionBootfileURL)
	oro.AddRequestedOption(OptionBootfileParam)
	msg.UpdateOption(oro)
	return d
}

// WithUserClass adds a user class option to the packet
func WithUserClass(uc []byte) Modifier {
	// TODO let the user specify multiple user classes
	return func(d DHCPv6) DHCPv6 {
		ouc := OptUserClass{UserClasses: [][]byte{uc}}
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
