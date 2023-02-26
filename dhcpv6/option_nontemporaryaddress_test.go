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

func TestIANAParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []*OptIANA
	}{
		{
			buf: []byte{
				0, 3, // IANA option code
				0, 40, // length
				1, 0, 0, 0, // IAID
				0, 0, 0, 1, // T1
				0, 0, 0, 2, // T2
				0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime
			},
			want: []*OptIANA{
				&OptIANA{
					IaId: [4]byte{1, 0, 0, 0},
					T1:   1 * time.Second,
					T2:   2 * time.Second,
					Options: IdentityOptions{Options: Options{&OptIAAddress{
						IPv6Addr:          net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
						PreferredLifetime: 2 * time.Second,
						ValidLifetime:     4 * time.Second,
						Options:           AddressOptions{Options: Options{}},
					}}},
				},
			},
		},
		{
			buf: []byte{
				0, 3, // IANA option code
				0, 40, // length
				1, 0, 0, 0, // IAID
				0, 0, 0, 1, // T1
				0, 0, 0, 2, // T2
				0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime

				0, 3, // IANA option code
				0, 40, // length
				1, 2, 3, 4, // IAID
				0, 0, 0, 9, // T1
				0, 0, 0, 8, // T2
				0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime
			},
			want: []*OptIANA{
				&OptIANA{
					IaId: [4]byte{1, 0, 0, 0},
					T1:   1 * time.Second,
					T2:   2 * time.Second,
					Options: IdentityOptions{Options: Options{&OptIAAddress{
						IPv6Addr:          net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
						PreferredLifetime: 2 * time.Second,
						ValidLifetime:     4 * time.Second,
						Options:           AddressOptions{Options: Options{}},
					}}},
				},
				&OptIANA{
					IaId: [4]byte{1, 2, 3, 4},
					T1:   9 * time.Second,
					T2:   8 * time.Second,
					Options: IdentityOptions{Options: Options{&OptIAAddress{
						IPv6Addr:          net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
						PreferredLifetime: 2 * time.Second,
						ValidLifetime:     4 * time.Second,
						Options:           AddressOptions{Options: Options{}},
					}}},
				},
			},
		},

		{
			buf:  nil,
			want: nil,
		},
		{
			buf:  []byte{0, 3, 0, 1, 0},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 3, // IANA option code
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
				0, 3, // IANA option code
				0, 36, // length
				1, 0, 0, 0, // IAID
				0, 0, 0, 1, // T1
				0, 0, 0, 2, // T2
				0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0xb2, 0x7a, // PreferredLifetime
				// Missing ValidLifetime
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
			if got := mo.IANA(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IANA = %v, want %v", got, tt.want)
			}
			var wantOne *OptIANA
			if len(tt.want) >= 1 {
				wantOne = tt.want[0]
			}
			if got := mo.OneIANA(); !reflect.DeepEqual(got, wantOne) {
				t.Errorf("OneIANA = %v, want %v", got, wantOne)
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

func TestOptIANAGetOneOption(t *testing.T) {
	oaddr := &OptIAAddress{
		IPv6Addr: net.ParseIP("::1"),
	}
	opt := OptIANA{
		Options: IdentityOptions{[]Option{&OptStatusCode{}, oaddr}},
	}
	require.Equal(t, oaddr, opt.Options.OneAddress())
}

func TestOptIANAAddOption(t *testing.T) {
	opt := OptIANA{}
	opt.Options.Add(OptElapsedTime(0))
	require.Equal(t, 1, len(opt.Options.Options))
	require.Equal(t, OptionElapsedTime, opt.Options.Options[0].Code())
}

func TestOptIANAGetOneOptionMissingOpt(t *testing.T) {
	oaddr := &OptIAAddress{
		IPv6Addr: net.ParseIP("::1"),
	}
	opt := OptIANA{
		Options: IdentityOptions{[]Option{&OptStatusCode{}, oaddr}},
	}
	require.Equal(t, nil, opt.Options.GetOne(OptionDNSRecursiveNameServer))
}

func TestOptIANADelOption(t *testing.T) {
	optiaaddr := OptIAAddress{}
	optsc := OptStatusCode{}

	iana1 := OptIANA{
		Options: IdentityOptions{[]Option{
			&optsc,
			&optiaaddr,
			&optiaaddr,
		}},
	}
	iana1.Options.Del(OptionIAAddr)
	require.Equal(t, iana1.Options.Options, Options{&optsc})

	iana2 := OptIANA{
		Options: IdentityOptions{[]Option{
			&optiaaddr,
			&optsc,
			&optiaaddr,
		}},
	}
	iana2.Options.Del(OptionIAAddr)
	require.Equal(t, iana2.Options.Options, Options{&optsc})
}

func TestOptIANAString(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
		0, 0, 0xb2, 0x7a, // PreferredLifetime
		0, 0, 0xc0, 0x8a, // ValidLifetime
	}
	var opt OptIANA
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
