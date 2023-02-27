package dhcpv6

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func TestIAPrefixParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []*OptIAPrefix
	}{
		{
			buf: []byte{
				0, 26, // IAPrefix option code
				0, 25, // length
				0, 0, 0, 1, // PreferredLifetime
				0, 0, 0, 2, // ValidLifetime
				16,                                                                        // prefix length
				0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // v6-prefix
			},
			want: []*OptIAPrefix{
				&OptIAPrefix{
					PreferredLifetime: 1 * time.Second,
					ValidLifetime:     2 * time.Second,
					Prefix: &net.IPNet{
						IP:   net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
						Mask: net.CIDRMask(16, 128),
					},
					Options: PrefixOptions{Options: Options{}},
				},
			},
		},
		{
			buf: []byte{
				0, 26, // IAPrefix option code
				0, 25, // length
				0, 0, 0, 1, // PreferredLifetime
				0, 0, 0, 2, // ValidLifetime
				0,                                              // prefix length
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // v6-prefix
			},
			want: []*OptIAPrefix{
				&OptIAPrefix{
					PreferredLifetime: 1 * time.Second,
					ValidLifetime:     2 * time.Second,
					Prefix:            nil,
					Options:           PrefixOptions{Options: Options{}},
				},
			},
		},
		{
			buf: []byte{
				0, 26, // IAPrefix option code
				0, 25, // length
				0, 0, 0, 1, // PreferredLifetime
				0, 0, 0, 2, // ValidLifetime
				16,                                                                        // prefix length
				0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // v6-prefix

				0, 26, // IAPrefix option code
				0, 25, // length
				0, 0, 0, 15, // PreferredLifetime
				0, 0, 0, 14, // ValidLifetime
				32,                                                                        // prefix length
				0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // v6-prefix
			},
			want: []*OptIAPrefix{
				&OptIAPrefix{
					PreferredLifetime: 1 * time.Second,
					ValidLifetime:     2 * time.Second,
					Prefix: &net.IPNet{
						IP:   net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
						Mask: net.CIDRMask(16, 128),
					},
					Options: PrefixOptions{Options: Options{}},
				},
				&OptIAPrefix{
					PreferredLifetime: 15 * time.Second,
					ValidLifetime:     14 * time.Second,
					Prefix: &net.IPNet{
						IP:   net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
						Mask: net.CIDRMask(32, 128),
					},
					Options: PrefixOptions{Options: Options{}},
				},
			},
		},
		{
			buf:  []byte{0, 3, 0, 1, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
		{
			buf: []byte{
				0, 26, // IAPrefix option code
				0, 8, // length
				1, 0, 0, 0, // T1
				0, 0, 0, 1, // T2
				// truncated from here
			},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 26, // IANA option code
				0, 26, // length
				0, 0, 0, 1, // T1
				0, 0, 0, 2, // T2
				8,
				0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, // malformed options
			},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo PDOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.Prefixes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prefixes = %#v, want %#v", got, tt.want)
			}

			if len(tt.want) >= 1 {
				var b PDOptions
				for _, iana := range tt.want {
					b.Add(iana)
				}
				got := b.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptIAPrefixString(t *testing.T) {
	buf := []byte{
		0x00, 0x00, 0x00, 60, // preferredLifetime
		0x00, 0x00, 0x00, 50, // validLifetime
		36,                                                         // prefixLength
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // ipv6Prefix
	}
	var opt OptIAPrefix
	err := opt.FromBytes(buf)
	require.NoError(t, err)

	str := opt.String()
	require.Contains(
		t, str,
		"Prefix=2001:db8::/36",
		"String() should return the ipv6addr",
	)
	require.Contains(
		t, str,
		"PreferredLifetime=1m",
		"String() should return the preferredlifetime",
	)
	require.Contains(
		t, str,
		"ValidLifetime=50s",
		"String() should return the validlifetime",
	)
}
