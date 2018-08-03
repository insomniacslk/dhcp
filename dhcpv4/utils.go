package dhcpv4

// IsRequested function takes a DHCPv4 message and an OptionCode, and returns
// true if that option is within the requested options of the DHCPv6 message.
func IsRequested(pkt *DHCPv4, requested OptionCode) bool {
	for _, optprl := range pkt.GetOption(OptionParameterRequestList) {
		for _, o := range optprl.(*OptParameterRequestList).RequestedOpts {
			if o == requested {
				return true
			}
		}
	}
	return false
}
