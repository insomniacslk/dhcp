package dhcpv4

import (
	"net"

	"github.com/insomniacslk/dhcp/rfc1035label"
)

// WithTransactionID sets the Transaction ID for the DHCPv4 packet
func WithTransactionID(xid TransactionID) Modifier {
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
		d.UpdateOption(opt)
		return d
	}
}

// WithUserClass adds a user class option to the packet.
// The rfc parameter allows you to specify if the userclass should be
// rfc compliant or not. More details in issue #113
func WithUserClass(uc []byte, rfc bool) Modifier {
	// TODO let the user specify multiple user classes
	return WithOption(&OptUserClass{
		UserClasses: [][]byte{uc},
		Rfc3004:     rfc,
	})
}

// WithNetboot adds bootfile URL and bootfile param options to a DHCPv4 packet.
func WithNetboot(d *DHCPv4) *DHCPv4 {
	return WithRequestedOptions(OptionTFTPServerName, OptionBootfileName)(d)
}

// WithRequestedOptions adds requested options to the packet.
func WithRequestedOptions(optionCodes ...OptionCode) Modifier {
	return func(d *DHCPv4) *DHCPv4 {
		params := d.GetOneOption(OptionParameterRequestList)
		if params == nil {
			d.UpdateOption(&OptParameterRequestList{OptionCodeList(optionCodes)})
		} else {
			opts := params.(*OptParameterRequestList)
			opts.RequestedOpts.Add(optionCodes...)
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
		d.HopCount += 1
		return d
	}
}

// WithNetmask adds or updates an OptSubnetMask
func WithNetmask(mask net.IPMask) Modifier {
	return WithOption(&OptSubnetMask{SubnetMask: mask})
}

// WithLeaseTime adds or updates an OptIPAddressLeaseTime
func WithLeaseTime(leaseTime uint32) Modifier {
	return WithOption(&OptIPAddressLeaseTime{LeaseTime: leaseTime})
}

// WithDNS adds or updates an OptionDomainNameServer
func WithDNS(dnses ...net.IP) Modifier {
	return WithOption(&OptDomainNameServer{NameServers: dnses})
}

// WithDomainSearchList adds or updates an OptionDomainSearch
func WithDomainSearchList(searchList ...string) Modifier {
	return WithOption(&OptDomainSearch{DomainSearch: &rfc1035label.Labels{
		Labels: searchList,
	}})
}

// WithRouter adds or updates an OptionRouter
func WithRouter(routers ...net.IP) Modifier {
	return WithOption(&OptRouter{Routers: routers})
}
