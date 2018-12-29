package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptDomainName(t *testing.T) {
	o := OptDomainName("foo")
	require.Equal(t, OptionDomainName, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Domain Name: foo", o.String())
}

func TestParseOptDomainName(t *testing.T) {
	o := Options{
		OptionDomainName.Code(): []byte{'t', 'e', 's', 't'},
	}
	require.Equal(t, "test", GetDomainName(o))
	require.Equal(t, "", GetDomainName(Options{}))
}

func TestOptHostName(t *testing.T) {
	o := OptHostName("foo")
	require.Equal(t, OptionHostName, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Host Name: foo", o.String())
}

func TestParseOptHostName(t *testing.T) {
	o := Options{
		OptionHostName.Code(): []byte{'t', 'e', 's', 't'},
	}
	require.Equal(t, "test", GetHostName(o))
	require.Equal(t, "", GetHostName(Options{}))
}

func TestOptRootPath(t *testing.T) {
	o := OptRootPath("foo")
	require.Equal(t, OptionRootPath, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Root Path: foo", o.String())
}

func TestParseOptRootPath(t *testing.T) {
	o := OptionsFromList(OptRootPath("test"))
	require.Equal(t, "test", GetRootPath(o))
	require.Equal(t, "", GetRootPath(Options{}))
}

func TestOptBootFileName(t *testing.T) {
	o := OptBootFileName("foo")
	require.Equal(t, OptionBootfileName, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Bootfile Name: foo", o.String())
}

func TestParseOptBootFileName(t *testing.T) {
	o := OptionsFromList(OptBootFileName("test"))
	require.Equal(t, "test", GetBootFileName(o))
	require.Equal(t, "", GetBootFileName(Options{}))
}

func TestOptTFTPServerName(t *testing.T) {
	o := OptTFTPServerName("foo")
	require.Equal(t, OptionTFTPServerName, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "TFTP Server Name: foo", o.String())
}

func TestParseOptTFTPServerName(t *testing.T) {
	o := OptionsFromList(OptTFTPServerName("test"))
	require.Equal(t, "test", GetTFTPServerName(o))
	require.Equal(t, "", GetTFTPServerName(Options{}))
}

func TestOptClassIdentifier(t *testing.T) {
	o := OptClassIdentifier("foo")
	require.Equal(t, OptionClassIdentifier, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Class Identifier: foo", o.String())
}

func TestParseOptClassIdentifier(t *testing.T) {
	o := OptionsFromList(OptClassIdentifier("test"))
	require.Equal(t, "test", GetClassIdentifier(o))
	require.Equal(t, "", GetClassIdentifier(Options{}))
}
