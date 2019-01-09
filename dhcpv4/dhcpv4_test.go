package dhcpv4

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestGetExternalIPv4Addrs(t *testing.T) {
	addrs4and6 := []net.Addr{
		&net.IPAddr{IP: net.IP{1, 2, 3, 4}},
		&net.IPAddr{IP: net.IP{4, 3, 2, 1}},
		&net.IPNet{IP: net.IP{4, 3, 2, 0}},
		&net.IPAddr{IP: net.IP{1, 2, 3, 4, 1, 1, 1, 1}},
		&net.IPAddr{IP: net.IP{4, 3, 2, 1, 1, 1, 1, 1}},
		&net.IPAddr{},                         // nil IP
		&net.IPAddr{IP: net.IP{127, 0, 0, 1}}, // loopback IP
	}

	expected := []net.IP{
		net.IP{1, 2, 3, 4},
		net.IP{4, 3, 2, 1},
		net.IP{4, 3, 2, 0},
	}
	actual, err := GetExternalIPv4Addrs(addrs4and6)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestFromBytes(t *testing.T) {
	data := []byte{
		1,                      // dhcp request
		1,                      // ethernet hw type
		6,                      // hw addr length
		3,                      // hop count
		0xaa, 0xbb, 0xcc, 0xdd, // transaction ID, big endian (network)
		0, 3, // number of seconds
		0, 1, // broadcast
		0, 0, 0, 0, // client IP address
		0, 0, 0, 0, // your IP address
		0, 0, 0, 0, // server IP address
		0, 0, 0, 0, // gateway IP address
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // client MAC address + padding
	}
	// server host name
	expectedHostname := []byte{}
	for i := 0; i < 64; i++ {
		expectedHostname = append(expectedHostname, 0)
	}
	data = append(data, expectedHostname...)
	// boot file name
	expectedBootfilename := []byte{}
	for i := 0; i < 128; i++ {
		expectedBootfilename = append(expectedBootfilename, 0)
	}
	data = append(data, expectedBootfilename...)
	// magic cookie, then no options
	data = append(data, []byte{99, 130, 83, 99}...)

	d, err := FromBytes(data)
	require.NoError(t, err)
	require.Equal(t, d.Opcode(), OpcodeBootRequest)
	require.Equal(t, d.HwType(), iana.HwTypeEthernet)
	require.Equal(t, d.HopCount(), byte(3))
	require.Equal(t, d.TransactionID(), TransactionID{0xaa, 0xbb, 0xcc, 0xdd})
	require.Equal(t, d.NumSeconds(), uint16(3))
	require.Equal(t, d.Flags(), uint16(1))
	require.True(t, d.ClientIPAddr().Equal(net.IPv4zero))
	require.True(t, d.YourIPAddr().Equal(net.IPv4zero))
	require.True(t, d.GatewayIPAddr().Equal(net.IPv4zero))
	require.Equal(t, d.ClientHwAddr(), net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff})
	hostname := d.ServerHostName()
	require.Equal(t, hostname[:], expectedHostname)
	bootfileName := d.BootFileName()
	require.Equal(t, bootfileName[:], expectedBootfilename)
	// no need to check Magic Cookie as it is already validated in FromBytes
	// above
}

func TestFromBytesZeroLength(t *testing.T) {
	data := []byte{}
	_, err := FromBytes(data)
	require.Error(t, err)
}

func TestFromBytesShortLength(t *testing.T) {
	data := []byte{1, 1, 6, 0}
	_, err := FromBytes(data)
	require.Error(t, err)
}

func TestFromBytesInvalidOptions(t *testing.T) {
	data := []byte{
		1,                      // dhcp request
		1,                      // ethernet hw type
		6,                      // hw addr length
		0,                      // hop count
		0xaa, 0xbb, 0xcc, 0xdd, // transaction ID
		3, 0, // number of seconds
		1, 0, // broadcast
		0, 0, 0, 0, // client IP address
		0, 0, 0, 0, // your IP address
		0, 0, 0, 0, // server IP address
		0, 0, 0, 0, // gateway IP address
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // client MAC address + padding
	}
	// server host name
	for i := 0; i < 64; i++ {
		data = append(data, 0)
	}
	// boot file name
	for i := 0; i < 128; i++ {
		data = append(data, 0)
	}
	// invalid magic cookie, forcing option parsing to fail
	data = append(data, []byte{99, 130, 83, 98}...)
	_, err := FromBytes(data)
	require.Error(t, err)
}

func TestSettersAndGetters(t *testing.T) {
	data := []byte{
		1,                      // dhcp request
		1,                      // ethernet hw type
		6,                      // hw addr length
		3,                      // hop count
		0xaa, 0xbb, 0xcc, 0xdd, // transaction ID, big endian (network)
		0, 3, // number of seconds
		0, 1, // broadcast
		1, 2, 3, 4, // client IP address
		5, 6, 7, 8, // your IP address
		9, 10, 11, 12, // server IP address
		13, 14, 15, 16, // gateway IP address
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // client MAC address + padding
	}
	// server host name
	expectedHostname := []byte{}
	for i := 0; i < 64; i++ {
		expectedHostname = append(expectedHostname, 0)
	}
	data = append(data, expectedHostname...)
	// boot file name
	expectedBootfilename := []byte{}
	for i := 0; i < 128; i++ {
		expectedBootfilename = append(expectedBootfilename, 0)
	}
	data = append(data, expectedBootfilename...)
	// magic cookie, then no options
	data = append(data, []byte{99, 130, 83, 99}...)
	d, err := FromBytes(data)
	require.NoError(t, err)

	// getter/setter for Opcode
	require.Equal(t, OpcodeBootRequest, d.Opcode())
	d.SetOpcode(OpcodeBootReply)
	require.Equal(t, OpcodeBootReply, d.Opcode())

	// getter/setter for HwType
	require.Equal(t, iana.HwTypeEthernet, d.HwType())
	d.SetHwType(iana.HwTypeARCNET)
	require.Equal(t, iana.HwTypeARCNET, d.HwType())

	// getter/setter for HopCount
	require.Equal(t, uint8(3), d.HopCount())
	d.SetHopCount(1)
	require.Equal(t, uint8(1), d.HopCount())

	// getter/setter for TransactionID
	require.Equal(t, TransactionID{0xaa, 0xbb, 0xcc, 0xdd}, d.TransactionID())
	d.SetTransactionID(TransactionID{0xee, 0xff, 0x00, 0x11})
	require.Equal(t, TransactionID{0xee, 0xff, 0x00, 0x11}, d.TransactionID())

	// getter/setter for TransactionID
	require.Equal(t, uint16(3), d.NumSeconds())
	d.SetNumSeconds(15)
	require.Equal(t, uint16(15), d.NumSeconds())

	// getter/setter for Flags
	require.Equal(t, uint16(1), d.Flags())
	d.SetFlags(0)
	require.Equal(t, uint16(0), d.Flags())

	// getter/setter for ClientIPAddr
	require.True(t, d.ClientIPAddr().Equal(net.IPv4(1, 2, 3, 4)))
	d.SetClientIPAddr(net.IPv4(4, 3, 2, 1))
	require.True(t, d.ClientIPAddr().Equal(net.IPv4(4, 3, 2, 1)))

	// getter/setter for YourIPAddr
	require.True(t, d.YourIPAddr().Equal(net.IPv4(5, 6, 7, 8)))
	d.SetYourIPAddr(net.IPv4(8, 7, 6, 5))
	require.True(t, d.YourIPAddr().Equal(net.IPv4(8, 7, 6, 5)))

	// getter/setter for ServerIPAddr
	require.True(t, d.ServerIPAddr().Equal(net.IPv4(9, 10, 11, 12)))
	d.SetServerIPAddr(net.IPv4(12, 11, 10, 9))
	require.True(t, d.ServerIPAddr().Equal(net.IPv4(12, 11, 10, 9)))

	// getter/setter for GatewayIPAddr
	require.True(t, d.GatewayIPAddr().Equal(net.IPv4(13, 14, 15, 16)))
	d.SetGatewayIPAddr(net.IPv4(16, 15, 14, 13))
	require.True(t, d.GatewayIPAddr().Equal(net.IPv4(16, 15, 14, 13)))

	// getter/setter for ClientHwAddr
	require.Equal(t, net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, d.ClientHwAddr())
	d.SetFlags(0)

	// getter/setter for ServerHostName
	serverhostname := d.ServerHostName()
	require.Equal(t, expectedHostname, serverhostname[:])
	newHostname := []byte{'t', 'e', 's', 't'}
	for i := 0; i < 60; i++ {
		newHostname = append(newHostname, 0)
	}
	d.SetServerHostName(newHostname)
	serverhostname = d.ServerHostName()
	require.Equal(t, newHostname, serverhostname[:])

	// getter/setter for BootFileName
	bootfilename := d.BootFileName()
	require.Equal(t, expectedBootfilename, bootfilename[:])
	newBootfilename := []byte{'t', 'e', 's', 't'}
	for i := 0; i < 124; i++ {
		newBootfilename = append(newBootfilename, 0)
	}
	d.SetBootFileName(newBootfilename)
	bootfilename = d.BootFileName()
	require.Equal(t, newBootfilename, bootfilename[:])
}

func TestToStringMethods(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	// OpcodeToString
	d.SetOpcode(OpcodeBootRequest)
	require.Equal(t, "BootRequest", d.OpcodeToString())
	d.SetOpcode(OpcodeBootReply)
	require.Equal(t, "BootReply", d.OpcodeToString())
	d.SetOpcode(OpcodeType(0))
	require.Equal(t, "Unknown", d.OpcodeToString())

	// HwTypeToString
	d.SetHwType(iana.HwTypeEthernet)
	require.Equal(t, "Ethernet", d.HwTypeToString())
	d.SetHwType(iana.HwTypeARCNET)
	require.Equal(t, "ARCNET", d.HwTypeToString())
	d.SetHwType(iana.HwTypeType(0))
	require.Equal(t, "Invalid", d.HwTypeToString())

	// FlagsToString
	d.SetUnicast()
	require.Equal(t, "Unicast", d.FlagsToString())
	d.SetBroadcast()
	require.Equal(t, "Broadcast", d.FlagsToString())
	d.SetFlags(0xffff)
	require.Equal(t, "Broadcast (reserved bits not zeroed)", d.FlagsToString())

	// ClientHwAddrToString
	d.SetClientHwAddr(net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff})
	require.Equal(t, "aa:bb:cc:dd:ee:ff", d.ClientHwAddrToString())

	// ServerHostNameToString
	d.SetServerHostName([]byte("my.host.local"))
	require.Equal(t, "my.host.local", d.ServerHostNameToString())

	// BootFileNameToString
	d.SetBootFileName([]byte("/my/boot/file"))
	require.Equal(t, "/my/boot/file", d.BootFileNameToString())
}

func TestNewToBytes(t *testing.T) {
	// the following bytes match what dhcpv4.New would create. Keep them in
	// sync!
	expected := []byte{
		1,                      // Opcode BootRequest
		1,                      // HwType Ethernet
		6,                      // HwAddrLen
		0,                      // HopCount
		0x11, 0x22, 0x33, 0x44, // TransactionID
		0, 0, // NumSeconds
		0, 0, // Flags
		0, 0, 0, 0, // ClientIPAddr
		0, 0, 0, 0, // YourIPAddr
		0, 0, 0, 0, // ServerIPAddr
		0, 0, 0, 0, // GatewayIPAddr
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // ClientHwAddr
	}
	// ServerHostName
	for i := 0; i < 64; i++ {
		expected = append(expected, 0)
	}
	// BootFileName
	for i := 0; i < 128; i++ {
		expected = append(expected, 0)
	}
	// Magic Cookie
	expected = append(expected, magicCookie[:]...)
	// End
	expected = append(expected, 0xff)

	d, err := New()
	require.NoError(t, err)
	// fix TransactionID to match the expected one, since it's randomly
	// generated in New()
	d.SetTransactionID(TransactionID{0x11, 0x22, 0x33, 0x44})
	got := d.ToBytes()
	require.Equal(t, expected, got)
}

func TestGetOption(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}

	hostnameOpt := &OptionGeneric{OptionCode: OptionHostName, Data: []byte("darkstar")}
	bootFileOpt1 := &OptBootfileName{[]byte("boot.img")}
	bootFileOpt2 := &OptBootfileName{[]byte("boot2.img")}
	d.AddOption(hostnameOpt)
	d.AddOption(&OptBootfileName{[]byte("boot.img")})
	d.AddOption(&OptBootfileName{[]byte("boot2.img")})

	require.Equal(t, d.GetOption(OptionHostName), []Option{hostnameOpt})
	require.Equal(t, d.GetOption(OptionBootfileName), []Option{bootFileOpt1, bootFileOpt2})
	require.Equal(t, d.GetOption(OptionRouter), []Option{})

	require.Equal(t, d.GetOneOption(OptionHostName), hostnameOpt)
	require.Equal(t, d.GetOneOption(OptionBootfileName), bootFileOpt1)
	require.Equal(t, d.GetOneOption(OptionRouter), nil)
}

func TestAddOption(t *testing.T) {
	d, err := New()
	require.NoError(t, err)

	hostnameOpt := &OptionGeneric{OptionCode: OptionHostName, Data: []byte("darkstar")}
	bootFileOpt1 := &OptionGeneric{OptionCode: OptionBootfileName, Data: []byte("boot.img")}
	bootFileOpt2 := &OptionGeneric{OptionCode: OptionBootfileName, Data: []byte("boot2.img")}
	d.AddOption(hostnameOpt)
	d.AddOption(bootFileOpt1)
	d.AddOption(bootFileOpt2)

	options := d.Options()
	require.Equal(t, len(options), 4)
	require.Equal(t, options[3].Code(), OptionEnd)
}

func TestUpdateOption(t *testing.T) {
	d, err := New()
	require.NoError(t, err)
	require.Equal(t, 1, len(d.options))
	require.Equal(t, OptionEnd, d.options[0].Code())
	// test that it will add the option since it's missing
	d.UpdateOption(&OptDomainName{DomainName: "slackware.it"})
	require.Equal(t, 2, len(d.options))
	require.Equal(t, OptionDomainName, d.options[0].Code())
	require.Equal(t, OptionEnd, d.options[1].Code())
	// test that it won't add another option of the same type
	d.UpdateOption(&OptDomainName{DomainName: "slackware.it"})
	require.Equal(t, 2, len(d.options))
	require.Equal(t, OptionDomainName, d.options[0].Code())
	require.Equal(t, OptionEnd, d.options[1].Code())
}

func TestStrippedOptions(t *testing.T) {
	// Normal set of options that terminate with OptionEnd.
	d, err := New()
	require.NoError(t, err)
	opts := []Option{
		&OptBootfileName{[]byte("boot.img")},
		&OptClassIdentifier{"something"},
		&OptionGeneric{OptionCode: OptionEnd},
	}
	d.SetOptions(opts)
	stripped := d.StrippedOptions()
	require.Equal(t, len(opts), len(stripped))
	for i := range stripped {
		require.Equal(t, opts[i], stripped[i])
	}

	// Set of options with additional options after OptionEnd
	opts = append(opts, &OptMaximumDHCPMessageSize{uint16(1234)})
	d.SetOptions(opts)
	stripped = d.StrippedOptions()
	require.Equal(t, len(opts)-1, len(stripped))
	for i := range stripped {
		require.Equal(t, opts[i], stripped[i])
	}
}

func TestDHCPv4NewRequestFromOffer(t *testing.T) {
	offer, err := New()
	require.NoError(t, err)
	offer.SetBroadcast()
	offer.AddOption(&OptMessageType{MessageType: MessageTypeOffer})
	req, err := NewRequestFromOffer(offer)
	require.Error(t, err)

	// Now add the option so it doesn't error out.
	offer.AddOption(&OptServerIdentifier{ServerID: net.IPv4(192, 168, 0, 1)})

	// Broadcast request
	req, err = NewRequestFromOffer(offer)
	require.NoError(t, err)
	require.NotNil(t, req.MessageType())
	require.Equal(t, MessageTypeRequest, *req.MessageType())
	require.False(t, req.IsUnicast())
	require.True(t, req.IsBroadcast())

	// Unicast request
	offer.SetUnicast()
	req, err = NewRequestFromOffer(offer)
	require.NoError(t, err)
	require.True(t, req.IsUnicast())
	require.False(t, req.IsBroadcast())
}

func TestDHCPv4NewRequestFromOfferWithModifier(t *testing.T) {
	offer, err := New()
	require.NoError(t, err)
	offer.AddOption(&OptMessageType{MessageType: MessageTypeOffer})
	offer.AddOption(&OptServerIdentifier{ServerID: net.IPv4(192, 168, 0, 1)})
	userClass := WithUserClass([]byte("linuxboot"), false)
	req, err := NewRequestFromOffer(offer, userClass)
	require.NoError(t, err)
	require.NotEqual(t, (*MessageType)(nil), *req.MessageType())
	require.Equal(t, MessageTypeRequest, *req.MessageType())
	require.Equal(t, "User Class Information -> linuxboot", req.options[3].String())
}

func TestNewReplyFromRequest(t *testing.T) {
	discover, err := New()
	require.NoError(t, err)
	discover.SetGatewayIPAddr(net.IPv4(192, 168, 0, 1))
	reply, err := NewReplyFromRequest(discover)
	require.NoError(t, err)
	require.Equal(t, discover.TransactionID(), reply.TransactionID())
	require.Equal(t, discover.GatewayIPAddr(), reply.GatewayIPAddr())
}

func TestNewReplyFromRequestWithModifier(t *testing.T) {
	discover, err := New()
	require.NoError(t, err)
	discover.SetGatewayIPAddr(net.IPv4(192, 168, 0, 1))
	userClass := WithUserClass([]byte("linuxboot"), false)
	reply, err := NewReplyFromRequest(discover, userClass)
	require.NoError(t, err)
	require.Equal(t, discover.TransactionID(), reply.TransactionID())
	require.Equal(t, discover.GatewayIPAddr(), reply.GatewayIPAddr())
	require.Equal(t, "User Class Information -> linuxboot", reply.options[0].String())
}

func TestDHCPv4MessageTypeNil(t *testing.T) {
	m, err := New()
	require.NoError(t, err)
	require.Nil(t, m.MessageType())
}

func TestNewDiscovery(t *testing.T) {
	hwAddr := net.HardwareAddr{1, 2, 3, 4, 5, 6}
	m, err := NewDiscovery(hwAddr)
	require.NoError(t, err)
	require.NotNil(t, m.MessageType())
	require.Equal(t, MessageTypeDiscover, *m.MessageType())

	// Validate fields of DISCOVER packet.
	require.Equal(t, OpcodeBootRequest, m.Opcode())
	require.Equal(t, iana.HwTypeEthernet, m.HwType())
	require.Equal(t, hwAddr, m.ClientHwAddr())
	require.True(t, m.IsBroadcast())
	require.True(t, HasOption(m, OptionParameterRequestList))
	require.True(t, HasOption(m, OptionEnd))
}

func TestNewInform(t *testing.T) {
	hwAddr := net.HardwareAddr{1, 2, 3, 4, 5, 6}
	localIP := net.IPv4(10, 10, 11, 11)
	m, err := NewInform(hwAddr, localIP)

	require.NoError(t, err)
	require.Equal(t, OpcodeBootRequest, m.Opcode())
	require.Equal(t, iana.HwTypeEthernet, m.HwType())
	require.Equal(t, hwAddr, m.ClientHwAddr())
	require.NotNil(t, m.MessageType())
	require.Equal(t, MessageTypeInform, *m.MessageType())
	require.True(t, m.ClientIPAddr().Equal(localIP))
}

func TestIsOptionRequested(t *testing.T) {
	pkt, err := New()
	require.NoError(t, err)
	require.False(t, pkt.IsOptionRequested(OptionDomainNameServer))

	optprl := OptParameterRequestList{RequestedOpts: []OptionCode{OptionDomainNameServer}}
	pkt.AddOption(&optprl)
	require.True(t, pkt.IsOptionRequested(OptionDomainNameServer))
}

// TODO
//      test Summary() and String()
