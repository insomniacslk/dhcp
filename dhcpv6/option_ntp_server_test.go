package dhcpv6

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/rfc1035label"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuboptionSrvAddr(t *testing.T) {
	ip := net.ParseIP("2a03:2880:fffe:c:face:b00c:0:35")
	so := NTPSuboptionSrvAddr(ip)
	assert.Equal(t, NTPSuboptionSrvAddrCode, so.Code())
	expected := append([]byte{0x00, 0x01, 0x00, 0x10}, ip...)
	assert.Equal(t, expected, so.ToBytes())
}

func TestSuboptionMCAddr(t *testing.T) {
	ip := net.ParseIP("2a03:2880:fffe:c:face:b00c:0:35")
	so := NTPSuboptionMCAddr(ip)
	assert.Equal(t, NTPSuboptionMCAddrCode, so.Code())
	expected := append([]byte{0x00, 0x02, 0x00, 0x10}, ip...)
	assert.Equal(t, expected, so.ToBytes())
}

func TestSuboptionSrvFQDN(t *testing.T) {
	fqdn, err := rfc1035label.FromBytes([]byte("\x03ntp\x07example\x03com"))
	require.NoError(t, err)
	so := NTPSuboptionSrvFQDN(*fqdn)
	assert.Equal(t, NTPSuboptionSrvFQDNCode, so.Code())
	expected := append([]byte{0x00, 0x03, 0x00, 0x10}, fqdn.ToBytes()...)
	assert.Equal(t, expected, so.ToBytes())
}

func TestSuboptionGeneric(t *testing.T) {
	data := []byte{
		0xff, 0xff, // unknown sub-option type
		0x00, 0x04, // length, 4 bytes
		0x74, 0x65, 0x73, 0x74, // the ASCII bytes for the string "test"
	}
	o, err := ParseOptNTPServer(data)
	require.NoError(t, err)
	require.Equal(t, 1, len(o.Suboptions))
	assert.IsType(t, &OptionGeneric{}, o.Suboptions[0])
	og := o.Suboptions[0].(*OptionGeneric)
	assert.Equal(t, []byte("test"), og.ToBytes())
}

func TestParseOptNTPServer(t *testing.T) {
	ip := net.ParseIP("2a03:2880:fffe:c:face:b00c:0:35")
	fqdn, err := rfc1035label.FromBytes([]byte("\x03ntp\x07example\x03com"))
	require.NoError(t, err)

	// add server address sub-option
	data := []byte{
		0x00, 0x01, // sub-option type
		0x00, 0x10, // length (16, IPv6 address)
	}
	data = append(data, []byte(ip)...)

	// add server FQDN sub-option
	data = append(data, []byte{
		0x00, 0x03, // sub-option type
		0x00, 0x10, // length (16, the FQDN "ntp.example.com." as rfc1035 label)
	}...)
	data = append(data, fqdn.ToBytes()...)

	o, err := ParseOptNTPServer(data)
	require.NoError(t, err)
	require.NotNil(t, o)
	assert.Equal(t, 2, len(o.Suboptions))

	optAddr, ok := o.Suboptions[0].(*NTPSuboptionSrvAddr)
	require.True(t, ok)
	assert.Equal(t, ip, net.IP(*optAddr))

	optFQDN, ok := o.Suboptions[1].(*NTPSuboptionSrvFQDN)
	require.True(t, ok)
	assert.Equal(t, *fqdn, rfc1035label.Labels(*optFQDN))
}
