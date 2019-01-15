package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptDomainName(t *testing.T) {
	o := OptDomainName{DomainName: "foo"}
	require.Equal(t, OptionDomainName, o.Code(), "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.ToBytes(), "ToBytes")
	require.Equal(t, "Domain Name -> foo", o.String())
}

func TestParseOptDomainName(t *testing.T) {
	data := []byte{'t', 'e', 's', 't'}
	o, err := ParseOptDomainName(data)
	require.NoError(t, err)
	require.Equal(t, &OptDomainName{DomainName: "test"}, o)
}

func TestOptHostName(t *testing.T) {
	o := OptHostName{HostName: "foo"}
	require.Equal(t, OptionHostName, o.Code(), "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.ToBytes(), "ToBytes")
	require.Equal(t, "Host Name -> foo", o.String())
}

func TestParseOptHostName(t *testing.T) {
	data := []byte{'t', 'e', 's', 't'}
	o, err := ParseOptHostName(data)
	require.NoError(t, err)
	require.Equal(t, &OptHostName{HostName: "test"}, o)
}

func TestOptRootPath(t *testing.T) {
	o := OptRootPath{Path: "/foo/bar/baz"}
	require.Equal(t, OptionRootPath, o.Code(), "Code")
	wantBytes := []byte{
		'/', 'f', 'o', 'o', '/', 'b', 'a', 'r', '/', 'b', 'a', 'z',
	}
	require.Equal(t, wantBytes, o.ToBytes(), "ToBytes")
	require.Equal(t, "Root Path -> /foo/bar/baz", o.String())
}

func TestParseOptRootPath(t *testing.T) {
	data := []byte{byte(OptionRootPath), 4, '/', 'f', 'o', 'o'}
	o, err := ParseOptRootPath(data[2:])
	require.NoError(t, err)
	require.Equal(t, &OptRootPath{Path: "/foo"}, o)
}

func TestOptBootfileName(t *testing.T) {
	opt := OptBootfileName{
		BootfileName: "linuxboot",
	}
	require.Equal(t, OptionBootfileName, opt.Code())
	require.Equal(t, []byte{'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't'}, opt.ToBytes())
	require.Equal(t, "Bootfile Name -> linuxboot", opt.String())
}

func TestParseOptBootfileName(t *testing.T) {
	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptBootfileName(expected)
	require.NoError(t, err)
	require.Equal(t, "linuxboot", opt.BootfileName)
}

func TestOptTFTPServer(t *testing.T) {
	opt := OptTFTPServerName{
		TFTPServerName: "linuxboot",
	}
	require.Equal(t, OptionTFTPServerName, opt.Code())
	require.Equal(t, []byte("linuxboot"), opt.ToBytes())
	require.Equal(t, "TFTP Server Name -> linuxboot", opt.String())
}

func TestParseOptTFTPServerName(t *testing.T) {
	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptTFTPServerName(expected)
	require.NoError(t, err)
	require.Equal(t, "linuxboot", string(opt.TFTPServerName))
}

func TestOptClassIdentifier(t *testing.T) {
	o := OptClassIdentifier{Identifier: "foo"}
	require.Equal(t, OptionClassIdentifier, o.Code(), "Code")
	require.Equal(t, []byte("foo"), o.ToBytes(), "ToBytes")
	require.Equal(t, "Class Identifier -> foo", o.String())
}

func TestParseOptClassIdentifier(t *testing.T) {
	data := []byte("test")
	o, err := ParseOptClassIdentifier(data)
	require.NoError(t, err)
	require.Equal(t, &OptClassIdentifier{Identifier: "test"}, o)
}
