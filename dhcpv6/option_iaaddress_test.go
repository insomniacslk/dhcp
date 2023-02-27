package dhcpv6

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func TestIAAddressParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []*OptIAAddress
	}{
		{
			buf: []byte{
				0, 5, // IAAddr option
				0, 0x18, // length
				0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime
			},
			want: []*OptIAAddress{
				&OptIAAddress{
					IPv6Addr:          net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
					PreferredLifetime: 2 * time.Second,
					ValidLifetime:     4 * time.Second,
					Options:           AddressOptions{Options: Options{}},
				},
			},
		},
		{
			buf: []byte{
				0, 5, // IAAddr option
				0, 32, // length
				0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime
				0, 13, // option status code
				0, 4, // length
				0, 0, // StatusSuccess,
				'O', 'K', // OK
			},
			want: []*OptIAAddress{
				&OptIAAddress{
					IPv6Addr:          net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
					PreferredLifetime: 2 * time.Second,
					ValidLifetime:     4 * time.Second,
					Options: AddressOptions{Options: Options{
						&OptStatusCode{StatusCode: iana.StatusSuccess, StatusMessage: "OK"},
					}},
				},
			},
		},
		{
			buf: []byte{
				0, 5, // IAAddr option
				0, 0x18, // length
				0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime

				0, 5, // IAAddr option
				0, 0x18, // length
				0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime
			},
			want: []*OptIAAddress{
				&OptIAAddress{
					IPv6Addr:          net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
					PreferredLifetime: 2 * time.Second,
					ValidLifetime:     4 * time.Second,
					Options:           AddressOptions{Options: Options{}},
				},
				&OptIAAddress{
					IPv6Addr:          net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
					PreferredLifetime: 2 * time.Second,
					ValidLifetime:     4 * time.Second,
					Options:           AddressOptions{Options: Options{}},
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
				0, 5, // IAAddr option code
				0, 4, // length
				0, 0, 0, 1, // truncated IP
			},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 5, // IAAddr option
				0, 28, // length
				0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime
				0, 13, // option status code
				0, 1, // length
				// option too short
			},
			err: uio.ErrBufferTooShort,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var io IdentityOptions
			if err := io.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := io.Addresses(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Addresses = %v, want %v", got, tt.want)
			}

			var wantOneAddr *OptIAAddress
			if len(tt.want) >= 1 {
				wantOneAddr = tt.want[0]
			}
			if got := io.OneAddress(); !reflect.DeepEqual(got, wantOneAddr) {
				t.Errorf("OneAddress = %v, want %v", got, wantOneAddr)
			}

			if len(tt.want) >= 1 {
				var b IdentityOptions
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

func TestOptIAAddressString(t *testing.T) {
	ipaddr := []byte{0x24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	data := append(ipaddr, []byte{
		0x00, 0x00, 0x00, 70, // preferred lifetime
		0x00, 0x00, 0x00, 50, // valid lifetime
		0, 8, 0, 2, 0xaa, 0xbb, // options
	}...)
	var opt OptIAAddress
	err := opt.FromBytes(data)
	require.NoError(t, err)

	str := opt.String()
	require.Contains(
		t, str,
		"IP=2401:203:405:607:809:a0b:c0d:e0f",
		"String() should return the ipv6addr",
	)
	require.Contains(
		t, str,
		"PreferredLifetime=1m10s",
		"String() should return the preferredlifetime",
	)
	require.Contains(
		t, str,
		"ValidLifetime=50s",
		"String() should return the validlifetime",
	)
}
