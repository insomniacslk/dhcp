package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptRootPathInterfaceMethods(t *testing.T) {
	o := OptRootPath{Path: "/foo/bar/baz"}
	require.Equal(t, OptionRootPath, o.Code(), "Code")
	require.Equal(t, 12, o.Length(), "Length")
	wantBytes := []byte{
		byte(OptionRootPath),
		12,
		'/', 'f', 'o', 'o', '/', 'b', 'a', 'r', '/', 'b', 'a', 'z',
	}
	require.Equal(t, wantBytes, o.ToBytes(), "ToBytes")
}

func TestParseOptRootPath(t *testing.T) {
	data := []byte{byte(OptionRootPath), 4, '/', 'f', 'o', 'o'}
	o, err := ParseOptRootPath(data[2:])
	require.NoError(t, err)
	require.Equal(t, &OptRootPath{Path: "/foo"}, o)
}

func TestOptRootPathString(t *testing.T) {
	o := OptRootPath{Path: "/foo/bar/baz"}
	require.Equal(t, "Root Path -> /foo/bar/baz", o.String())
}
