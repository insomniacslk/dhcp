package dhcpv4

import (
	"bytes"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
)

// NOTE: if one of the following Assert* fails where expected and got values are
// the same, you probably have to cast one of them to match the other one's
// type, e.g. comparing int and byte, even the same value, will fail.
func AssertEqual(t *testing.T, a, b interface{}, what string) {
	if a != b {
		t.Fatalf("Invalid %s. %v != %v", what, a, b)
	}
}

func AssertEqualBytes(t *testing.T, a, b []byte, what string) {
	if !bytes.Equal(a, b) {
		t.Fatalf("Invalid %s. %v != %v", what, a, b)
	}
}

func AssertEqualIPAddr(t *testing.T, a, b net.IP, what string) {
	if !net.IP.Equal(a, b) {
		t.Fatalf("Invalid %s. %v != %v", what, a, b)
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
	if err != nil {
		t.Fatal(err)
	}
	AssertEqual(t, d.Opcode(), OpcodeBootRequest, "opcode")
	AssertEqual(t, d.HwType(), iana.HwTypeEthernet, "hardware type")
	AssertEqual(t, d.HwAddrLen(), byte(6), "hardware address length")
	AssertEqual(t, d.HopCount(), byte(3), "hop count")
	AssertEqual(t, d.TransactionID(), uint32(0xaabbccdd), "transaction ID")
	AssertEqual(t, d.NumSeconds(), uint16(3), "number of seconds")
	AssertEqual(t, d.Flags(), uint16(1), "flags")
	AssertEqualIPAddr(t, d.ClientIPAddr(), net.IPv4zero, "client IP address")
	AssertEqualIPAddr(t, d.YourIPAddr(), net.IPv4zero, "your IP address")
	AssertEqualIPAddr(t, d.ServerIPAddr(), net.IPv4zero, "server IP address")
	AssertEqualIPAddr(t, d.GatewayIPAddr(), net.IPv4zero, "gateway IP address")
	hwaddr := d.ClientHwAddr()
	AssertEqualBytes(t, hwaddr[:], []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "flags")
	hostname := d.ServerHostName()
	AssertEqualBytes(t, hostname[:], expectedHostname, "server host name")
	bootfilename := d.BootFileName()
	AssertEqualBytes(t, bootfilename[:], expectedBootfilename, "boot file name")
	// no need to check Magic Cookie as it is already validated in FromBytes
	// above
}

func TestFromBytesZeroLength(t *testing.T) {
	data := []byte{}
	_, err := FromBytes(data)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestFromBytesShortLength(t *testing.T) {
	data := []byte{1, 1, 6, 0}
	_, err := FromBytes(data)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
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
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
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
	if err != nil {
		t.Fatal(err)
	}

	// getter/setter for Opcode
	AssertEqual(t, d.Opcode(), OpcodeBootRequest, "opcode")
	d.SetOpcode(OpcodeBootReply)
	AssertEqual(t, d.Opcode(), OpcodeBootReply, "opcode")

	// getter/setter for HwType
	AssertEqual(t, d.HwType(), iana.HwTypeEthernet, "hardware type")
	d.SetHwType(iana.HwTypeARCNET)
	AssertEqual(t, d.HwType(), iana.HwTypeARCNET, "hardware type")

	// getter/setter for HwAddrLen
	AssertEqual(t, d.HwAddrLen(), uint8(6), "hardware address length")
	d.SetHwAddrLen(12)
	AssertEqual(t, d.HwAddrLen(), uint8(12), "hardware address length")

	// getter/setter for HopCount
	AssertEqual(t, d.HopCount(), uint8(3), "hop count")
	d.SetHopCount(1)
	AssertEqual(t, d.HopCount(), uint8(1), "hop count")

	// getter/setter for TransactionID
	AssertEqual(t, d.TransactionID(), uint32(0xaabbccdd), "transaction ID")
	d.SetTransactionID(0xeeff0011)
	AssertEqual(t, d.TransactionID(), uint32(0xeeff0011), "transaction ID")

	// getter/setter for TransactionID
	AssertEqual(t, d.NumSeconds(), uint16(3), "number of seconds")
	d.SetNumSeconds(15)
	AssertEqual(t, d.NumSeconds(), uint16(15), "number of seconds")

	// getter/setter for Flags
	AssertEqual(t, d.Flags(), uint16(1), "flags")
	d.SetFlags(0)
	AssertEqual(t, d.Flags(), uint16(0), "flags")

	// getter/setter for ClientIPAddr
	AssertEqualIPAddr(t, d.ClientIPAddr(), net.IPv4(1, 2, 3, 4), "client IP address")
	d.SetClientIPAddr(net.IPv4(4, 3, 2, 1))
	AssertEqualIPAddr(t, d.ClientIPAddr(), net.IPv4(4, 3, 2, 1), "client IP address")

	// getter/setter for YourIPAddr
	AssertEqualIPAddr(t, d.YourIPAddr(), net.IPv4(5, 6, 7, 8), "your IP address")
	d.SetYourIPAddr(net.IPv4(8, 7, 6, 5))
	AssertEqualIPAddr(t, d.YourIPAddr(), net.IPv4(8, 7, 6, 5), "your IP address")

	// getter/setter for ServerIPAddr
	AssertEqualIPAddr(t, d.ServerIPAddr(), net.IPv4(9, 10, 11, 12), "server IP address")
	d.SetServerIPAddr(net.IPv4(12, 11, 10, 9))
	AssertEqualIPAddr(t, d.ServerIPAddr(), net.IPv4(12, 11, 10, 9), "server IP address")

	// getter/setter for GatewayIPAddr
	AssertEqualIPAddr(t, d.GatewayIPAddr(), net.IPv4(13, 14, 15, 16), "gateway IP address")
	d.SetGatewayIPAddr(net.IPv4(16, 15, 14, 13))
	AssertEqualIPAddr(t, d.GatewayIPAddr(), net.IPv4(16, 15, 14, 13), "gateway IP address")

	// getter/setter for ClientHwAddr
	hwaddr := d.ClientHwAddr()
	AssertEqualBytes(t, hwaddr[:], []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "client hardware address")
	d.SetFlags(0)

	// getter/setter for ServerHostName
	serverhostname := d.ServerHostName()
	AssertEqualBytes(t, serverhostname[:], expectedHostname, "server host name")
	newHostname := []byte{'t', 'e', 's', 't'}
	for i := 0; i < 60; i++ {
		newHostname = append(newHostname, 0)
	}
	d.SetServerHostName(newHostname)
	serverhostname = d.ServerHostName()
	AssertEqualBytes(t, serverhostname[:], newHostname, "server host name")

	// getter/setter for BootFileName
	bootfilename := d.BootFileName()
	AssertEqualBytes(t, bootfilename[:], expectedBootfilename, "boot file name")
	newBootfilename := []byte{'t', 'e', 's', 't'}
	for i := 0; i < 124; i++ {
		newBootfilename = append(newBootfilename, 0)
	}
	d.SetBootFileName(newBootfilename)
	bootfilename = d.BootFileName()
	AssertEqualBytes(t, bootfilename[:], newBootfilename, "boot file name")
}

func TestToStringMethods(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	// OpcodeToString
	d.SetOpcode(OpcodeBootRequest)
	AssertEqual(t, d.OpcodeToString(), "BootRequest", "OpcodeToString")
	d.SetOpcode(OpcodeBootReply)
	AssertEqual(t, d.OpcodeToString(), "BootReply", "OpcodeToString")
	d.SetOpcode(OpcodeType(0))
	AssertEqual(t, d.OpcodeToString(), "Invalid", "OpcodeToString")

	// HwTypeToString
	d.SetHwType(iana.HwTypeEthernet)
	AssertEqual(t, d.HwTypeToString(), "Ethernet", "HwTypeToString")
	d.SetHwType(iana.HwTypeARCNET)
	AssertEqual(t, d.HwTypeToString(), "ARCNET", "HwTypeToString")

	// FlagsToString
	d.SetUnicast()
	AssertEqual(t, d.FlagsToString(), "Unicast", "FlagsToString")
	d.SetBroadcast()
	AssertEqual(t, d.FlagsToString(), "Broadcast", "FlagsToString")
	d.SetFlags(0xffff)
	AssertEqual(t, d.FlagsToString(), "Broadcast (reserved bits not zeroed)", "FlagsToString")

	// ClientHwAddrToString
	d.SetHwAddrLen(6)
	d.SetClientHwAddr([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	AssertEqual(t, d.ClientHwAddrToString(), "aa:bb:cc:dd:ee:ff", "ClientHwAddrToString")

	// ServerHostNameToString
	d.SetServerHostName([]byte("my.host.local"))
	AssertEqual(t, d.ServerHostNameToString(), "my.host.local", "ServerHostNameToString")

	// BootFileNameToString
	d.SetBootFileName([]byte("/my/boot/file"))
	AssertEqual(t, d.BootFileNameToString(), "/my/boot/file", "BootFileNameToString")
}

func TestToBytes(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	// fix TransactionID to match the expected one, since it's randomly
	// generated in New()
	d.SetTransactionID(0x11223344)
	got := d.ToBytes()
	AssertEqualBytes(t, expected, got, "ToBytes")
}

// TODO
//      test broadcast/unicast flags
//      test Options setter/getter
//      test Summary() and String()
