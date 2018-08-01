package dhcpv4

// WithUserClass adds a user class option to the packet
func WithUserClass(uc []byte) Modifier {
	// TODO let the user specify multiple user classes
	return func(d *DHCPv4) *DHCPv4 {
		ouc := OptUserClass{UserClasses: [][]byte{uc}}
		d.AddOption(&ouc)
		return d
	}
}