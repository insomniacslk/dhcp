package dhcpv6

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptFQDN(t *testing.T) {
	data := []byte{
		0, // Flags
		'c', 'n', 'o', 's', '.', 'l', 'o', 'c', 'a', 'l',
		'h', 'o', 's', 't',
	}
	opt, err := ParseOptFQDN(data)

	require.NoError(t, err)
	require.Equal(t, OptionFQDN, opt.Code())
	require.Equal(t, uint8(0), opt.Flags)
	require.Equal(t, "cnos.localhost", opt.DomainName)
	require.Equal(t, "OptFQDN{flags=0, domainname=cnos.localhost}", opt.String())
}

func TestOptFQDNToBytes(t *testing.T) {
	opt := OptFQDN{
		Flags:      0,
		DomainName: "cnos.localhost",
	}
	want := []byte{
		0, // Flags
		'c', 'n', 'o', 's', '.', 'l', 'o', 'c', 'a', 'l',
		'h', 'o', 's', 't',
	}
	b := opt.ToBytes()
	if !bytes.Equal(b, want) {
		t.Fatalf("opt.ToBytes()=%v, want %v", b, want)
	}
}
