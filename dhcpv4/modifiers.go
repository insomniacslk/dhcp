package dhcpv4

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
