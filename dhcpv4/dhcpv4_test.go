package dhcpv4

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func RequireEqualIPAddr(t *testing.T, a, b net.IP, msg ...interface{}) {
	if !net.IP.Equal(a, b) {
		t.Fatalf("Invalid %s. %v != %v", msg, a, b)
	}
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
	require.Equal(t, d.HwAddrLen(), byte(6))
	require.Equal(t, d.HopCount(), byte(3))
	require.Equal(t, d.TransactionID(), uint32(0xaabbccdd))
	require.Equal(t, d.NumSeconds(), uint16(3))
	require.Equal(t, d.Flags(), uint16(1))
	RequireEqualIPAddr(t, d.ClientIPAddr(), net.IPv4zero)
	RequireEqualIPAddr(t, d.YourIPAddr(), net.IPv4zero)
	RequireEqualIPAddr(t, d.GatewayIPAddr(), net.IPv4zero)
	clientHwAddr := d.ClientHwAddr()
	require.Equal(t, clientHwAddr[:], []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
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

	// getter/setter for HwAddrLen
	require.Equal(t, uint8(6), d.HwAddrLen())
	d.SetHwAddrLen(12)
	require.Equal(t, uint8(12), d.HwAddrLen())

	// getter/setter for HopCount
	require.Equal(t, uint8(3), d.HopCount())
	d.SetHopCount(1)
	require.Equal(t, uint8(1), d.HopCount())

	// getter/setter for TransactionID
	require.Equal(t, uint32(0xaabbccdd), d.TransactionID())
	d.SetTransactionID(0xeeff0011)
	require.Equal(t, uint32(0xeeff0011), d.TransactionID())

	// getter/setter for TransactionID
	require.Equal(t, uint16(3), d.NumSeconds())
	d.SetNumSeconds(15)
	require.Equal(t, uint16(15), d.NumSeconds())

	// getter/setter for Flags
	require.Equal(t, uint16(1), d.Flags())
	d.SetFlags(0)
	require.Equal(t, uint16(0), d.Flags())

	// getter/setter for ClientIPAddr
	RequireEqualIPAddr(t, net.IPv4(1, 2, 3, 4), d.ClientIPAddr())
	d.SetClientIPAddr(net.IPv4(4, 3, 2, 1))
	RequireEqualIPAddr(t, net.IPv4(4, 3, 2, 1), d.ClientIPAddr())

	// getter/setter for YourIPAddr
	RequireEqualIPAddr(t, net.IPv4(5, 6, 7, 8), d.YourIPAddr())
	d.SetYourIPAddr(net.IPv4(8, 7, 6, 5))
	RequireEqualIPAddr(t, net.IPv4(8, 7, 6, 5), d.YourIPAddr())

	// getter/setter for ServerIPAddr
	RequireEqualIPAddr(t, net.IPv4(9, 10, 11, 12), d.ServerIPAddr())
	d.SetServerIPAddr(net.IPv4(12, 11, 10, 9))
	RequireEqualIPAddr(t, net.IPv4(12, 11, 10, 9), d.ServerIPAddr())

	// getter/setter for GatewayIPAddr
	RequireEqualIPAddr(t, net.IPv4(13, 14, 15, 16), d.GatewayIPAddr())
	d.SetGatewayIPAddr(net.IPv4(16, 15, 14, 13))
	RequireEqualIPAddr(t, net.IPv4(16, 15, 14, 13), d.GatewayIPAddr())

	// getter/setter for ClientHwAddr
	hwaddr := d.ClientHwAddr()
	require.Equal(t, []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, hwaddr[:])
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
	require.Equal(t, "Invalid", d.OpcodeToString())

	// HwTypeToString
	d.SetHwType(iana.HwTypeEthernet)
	require.Equal(t, "Ethernet", d.HwTypeToString())
	d.SetHwType(iana.HwTypeARCNET)
	require.Equal(t, "ARCNET", d.HwTypeToString())

	// FlagsToString
	d.SetUnicast()
	require.Equal(t, "Unicast", d.FlagsToString())
	d.SetBroadcast()
	require.Equal(t, "Broadcast", d.FlagsToString())
	d.SetFlags(0xffff)
	require.Equal(t, "Broadcast (reserved bits not zeroed)", d.FlagsToString())

	// ClientHwAddrToString
	d.SetHwAddrLen(6)
	d.SetClientHwAddr([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
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
	expected = append(expected, MagicCookie...)

	d, err := New()
	require.NoError(t, err)
	// fix TransactionID to match the expected one, since it's randomly
	// generated in New()
	d.SetTransactionID(0x11223344)
	got := d.ToBytes()
	require.Equal(t, expected, got)
}

func TestGetOption(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}

	hostnameOpt := &OptionGeneric{OptionCode: OptionHostName, Data: []byte("darkstar")}
	bootFileOpt1 := &OptionGeneric{OptionCode: OptionBootfileName, Data: []byte("boot.img")}
	bootFileOpt2 := &OptionGeneric{OptionCode: OptionBootfileName, Data: []byte("boot2.img")}
	d.AddOption(hostnameOpt)
	d.AddOption(bootFileOpt1)
	d.AddOption(bootFileOpt2)

	require.Equal(t, d.GetOption(OptionHostName), []Option{hostnameOpt})
	require.Equal(t, d.GetOption(OptionBootfileName), []Option{bootFileOpt1, bootFileOpt2})
	require.Equal(t, d.GetOption(OptionRouter), []Option{})

	require.Equal(t, d.GetOneOption(OptionHostName), hostnameOpt)
	require.Equal(t, d.GetOneOption(OptionBootfileName), bootFileOpt1)
	require.Equal(t, d.GetOneOption(OptionRouter), nil)
}

func TestDHCPv4MessageTypeNil(t *testing.T) {
	m, err := New()
	require.NoError(t, err)
	require.Equal(t, (*MessageType)(nil), m.MessageType())
}

func TestDHCPv4MessageTypeDiscovery(t *testing.T) {
	m, err := NewDiscoveryForInterface("lo")
	require.NoError(t, err)
	require.NotEqual(t, (*MessageType)(nil), m.MessageType())
	require.Equal(t, MessageTypeDiscover, *m.MessageType())
}

// TODO
//      test broadcast/unicast flags
//      test Options setter/getter
//      test Summary() and String()
