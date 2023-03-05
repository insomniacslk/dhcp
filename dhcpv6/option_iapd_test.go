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

func TestIAPDParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []*OptIAPD
	}{
		{
			buf: []byte{
				0, 25, // IAPD option code
				0, 41, // length
				1, 0, 0, 0, // IAID
				0, 0, 0, 1, // T1
				0, 0, 0, 2, // T2
				0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
				0, 0, 0, 2, // IAPrefix preferredLifetime
				0, 0, 0, 4, // IAPrefix validLifetime
				36,                                             // IAPrefix prefixLength
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix
			},
			want: []*OptIAPD{
				&OptIAPD{
					IaId: [4]byte{1, 0, 0, 0},
					T1:   1 * time.Second,
					T2:   2 * time.Second,
					Options: PDOptions{Options: Options{&OptIAPrefix{
						PreferredLifetime: 2 * time.Second,
						ValidLifetime:     4 * time.Second,
						Prefix: &net.IPNet{
							Mask: net.CIDRMask(36, 128),
							IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
						},
						Options: PrefixOptions{Options: Options{}},
					}}},
				},
			},
		},
		{
			buf: []byte{
				0, 25, // IAPD option code
				0, 41, // length
				1, 0, 0, 0, // IAID
				0, 0, 0, 1, // T1
				0, 0, 0, 2, // T2
				0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
				0, 0, 0, 2, // IAPrefix preferredLifetime
				0, 0, 0, 4, // IAPrefix validLifetime
				36,                                             // IAPrefix prefixLength
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix

				0, 25, // IAPD option code
				0, 41, // length
				1, 2, 3, 4, // IAID
				0, 0, 0, 5, // T1
				0, 0, 0, 6, // T2
				0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
				0, 0, 0, 2, // IAPrefix preferredLifetime
				0, 0, 0, 4, // IAPrefix validLifetime
				36,                                             // IAPrefix prefixLength
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix
			},
			want: []*OptIAPD{
				&OptIAPD{
					IaId: [4]byte{1, 0, 0, 0},
					T1:   1 * time.Second,
					T2:   2 * time.Second,
					Options: PDOptions{Options: Options{&OptIAPrefix{
						PreferredLifetime: 2 * time.Second,
						ValidLifetime:     4 * time.Second,
						Prefix: &net.IPNet{
							Mask: net.CIDRMask(36, 128),
							IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
						},
						Options: PrefixOptions{Options: Options{}},
					}}},
				},
				&OptIAPD{
					IaId: [4]byte{1, 2, 3, 4},
					T1:   5 * time.Second,
					T2:   6 * time.Second,
					Options: PDOptions{Options: Options{&OptIAPrefix{
						PreferredLifetime: 2 * time.Second,
						ValidLifetime:     4 * time.Second,
						Prefix: &net.IPNet{
							Mask: net.CIDRMask(36, 128),
							IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
						},
						Options: PrefixOptions{Options: Options{}},
					}}},
				},
			},
		},
		{
			buf:  nil,
			want: nil,
		},
		{
			buf:  []byte{0, 25, 0, 1, 0},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 25, // IAPD option code
				0, 8, // length
				1, 0, 0, 0, // IAID
				0, 0, 0, 1, // T1
				// truncated from here
			},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 25, // IAPD option code
				0, 36, // length
				1, 0, 0, 0, // IAID
				0, 0, 0, 1, // T1
				0, 0, 0, 2, // T2
				0, 26, 0, 4, // 26 = IAPrefix Option, 4 = length
				0, 0, 0, 2, // IAPrefix preferredLifetime
				// Missing stuff
			},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.IAPD(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAPD = %v, want %v", got, tt.want)
			}
			var wantOne *OptIAPD
			if len(tt.want) >= 1 {
				wantOne = tt.want[0]
			}
			if got := mo.OneIAPD(); !reflect.DeepEqual(got, wantOne) {
				t.Errorf("OneIAPD = %v, want %v", got, wantOne)
			}

			if len(tt.want) >= 1 {
				var b MessageOptions
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

func TestOptIAPDString(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
		0xaa, 0xbb, 0xcc, 0xdd, // IAPrefix preferredLifetime
		0xee, 0xff, 0x00, 0x11, // IAPrefix validLifetime
		36,                                             // IAPrefix prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix
	}
	var opt OptIAPD
	err := opt.FromBytes(data)
	require.NoError(t, err)

	str := opt.String()
	require.Contains(
		t, str,
		"IAID=0x01000000",
		"String() should return the IAID",
	)
	require.Contains(
		t, str,
		"T1=1s T2=2s",
		"String() should return the T1/T2 options",
	)
	require.Contains(
		t, str,
		"Options={",
		"String() should return a list of options",
	)
}
