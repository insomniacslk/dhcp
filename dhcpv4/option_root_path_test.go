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
	o, err := ParseOptRootPath(data)
	require.NoError(t, err)
	require.Equal(t, &OptRootPath{Path: "/foo"}, o)

	// Short byte stream
	data = []byte{byte(OptionRootPath)}
	_, err = ParseOptRootPath(data)
	require.Error(t, err, "should get error from short byte stream")

	// Wrong code
	data = []byte{43, 2, 1, 1}
	_, err = ParseOptRootPath(data)
	require.Error(t, err, "should get error from wrong code")

	// Bad length
	data = []byte{byte(OptionRootPath), 6, 1, 1, 1}
	_, err = ParseOptRootPath(data)
	require.Error(t, err, "should get error from bad length")
}

func TestOptRootPathString(t *testing.T) {
	o := OptRootPath{Path: "/foo/bar/baz"}
	require.Equal(t, "Root Path -> /foo/bar/baz", o.String())
}
