package dhcpv6

// Selector defines the signature for functions that can select DHCPv6
// structures. This is used to drop illegal packets.
type Selector func(d *Message) bool

// WithDefault returns true by default
func WithDefault() Selector {
	return func(d *Message) bool {
		return true
	}
}
