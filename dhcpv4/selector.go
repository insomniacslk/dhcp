package dhcpv4

// Selector defines the signature for functions that can select DHCPv4
// structures. This is used to drop illegal packets.
type Selector func(d *DHCPv4) bool

// WithDefault returns true by default
func WithDefault() Selector {
	return func(d *DHCPv4) bool {
		return true
	}
}
