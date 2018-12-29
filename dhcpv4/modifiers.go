package dhcpv4

import (
	"net"

	"github.com/insomniacslk/dhcp/rfc1035label"
)

// WithTransactionID sets the Transaction ID for the DHCPv4 packet
func WithTransactionID(xid [4]byte) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		d.TransactionID = xid
		return d
	}
}

// WithBroadcast sets the packet to be broadcast or unicast
func WithBroadcast(broadcast bool) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		if broadcast {
			d.SetBroadcast()
		} else {
			d.SetUnicast()
		}
		return d
	}
}

// WithHwAddr sets the hardware address for a packet
func WithHwAddr(hwaddr net.HardwareAddr) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		d.ClientHWAddr = hwaddr
		return d
	}
}

// WithOption appends a DHCPv4 option provided by the user
func WithOption(opt Option) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		d.AddOption(opt)
		return d
	}
}

// WithUserClass adds a user class option to the packet.
// The rfc parameter allows you to specify if the userclass should be
// rfc compliant or not. More details in issue #113
func WithUserClass(uc []byte, rfc bool) Modifier {
	// TODO let the user specify multiple user classes
	return func(d *DHCPv4) *DHCPv4 {
		ouc := OptUserClass{
			UserClasses: [][]byte{uc},
			Rfc3004:     rfc,
		}
		d.AddOption(&ouc)
		return d
	}
}

// WithNetboot adds bootfile URL and bootfile param options to a DHCPv4 packet.
func WithNetboot(d *DHCPv4) *DHCPv4 {
	params := d.GetOneOption(OptionParameterRequestList)

	var (
		OptParams                 *OptParameterRequestList
		foundOptionTFTPServerName bool
		foundOptionBootfileName   bool
	)
	if params != nil {
		OptParams = params.(*OptParameterRequestList)
		for _, option := range OptParams.RequestedOpts {
			if option == OptionTFTPServerName {
				foundOptionTFTPServerName = true
			} else if option == OptionBootfileName {
				foundOptionBootfileName = true
			}
		}
		if !foundOptionTFTPServerName {
			OptParams.RequestedOpts = append(OptParams.RequestedOpts, OptionTFTPServerName)
		}
		if !foundOptionBootfileName {
			OptParams.RequestedOpts = append(OptParams.RequestedOpts, OptionBootfileName)
		}
	} else {
		OptParams = &OptParameterRequestList{
			RequestedOpts: []OptionCode{OptionTFTPServerName, OptionBootfileName},
		}
		d.AddOption(OptParams)
	}
	return d
}

// WithRequestedOptions adds requested options to the packet
func WithRequestedOptions(optionCodes ...OptionCode) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		params := d.GetOneOption(OptionParameterRequestList)
		if params == nil {
			params = &OptParameterRequestList{}
			d.AddOption(params)
		}
		opts := params.(*OptParameterRequestList)
		for _, optionCode := range optionCodes {
			opts.RequestedOpts = append(opts.RequestedOpts, optionCode)
		}
		return d
	}
}

// WithRelay adds parameters required for DHCPv4 to be relayed by the relay
// server with given ip
func WithRelay(ip net.IP) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		d.SetUnicast()
		d.GatewayIPAddr = ip
		d.HopCount = 1
		return d
	}
}

// WithNetmask adds or updates an OptSubnetMask
func WithNetmask(mask net.IPMask) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		osm := OptSubnetMask{
			SubnetMask: mask,
		}
		d.UpdateOption(&osm)
		return d
	}
}

// WithLeaseTime adds or updates an OptIPAddressLeaseTime
func WithLeaseTime(leaseTime uint32) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		olt := OptIPAddressLeaseTime{
			LeaseTime: leaseTime,
		}
		d.UpdateOption(&olt)
		return d
	}
}

// WithDNS adds or updates an OptionDomainNameServer
func WithDNS(dnses ...net.IP) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		odns := OptDomainNameServer{NameServers: dnses}
		d.UpdateOption(&odns)
		return d
	}
}

// WithDomainSearchList adds or updates an OptionDomainSearch
func WithDomainSearchList(searchList ...string) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		labels := rfc1035label.Labels{
			Labels: searchList,
		}
		odsl := OptDomainSearch{DomainSearch: &labels}
		d.UpdateOption(&odsl)
		return d
	}
}

// WithRouter adds or updates an OptionRouter
func WithRouter(routers ...net.IP) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		ortr := OptRouter{Routers: routers}
		d.UpdateOption(&ortr)
		return d
	}
}
