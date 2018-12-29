package dhcpv4

// OptMessageType returns a new DHCPv4 Message Type option.
func OptMessageType(m MessageType) Option {
	return Option{Code: OptionDHCPMessageType, Value: m}
}

// GetMessageType returns the DHCPv4 Message Type option in o.
func GetMessageType(o Options) MessageType {
	v := o.Get(OptionDHCPMessageType)
	if v == nil {
		return MessageTypeNone
	}
	var m MessageType
	if err := m.FromBytes(v); err != nil {
		return MessageTypeNone
	}
	return m
}
