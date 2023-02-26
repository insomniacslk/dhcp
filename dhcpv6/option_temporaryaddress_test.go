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

func TestIATAParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []*OptIATA
	}{
		{
			buf: []byte{
				0, 4, // IATA option code
				0, 32, // length
				1, 0, 0, 0, // IAID
				0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime
			},
			want: []*OptIATA{
				&OptIATA{
					IaId: [4]byte{1, 0, 0, 0},
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
				0, 4, // IATA option code
				0, 32, // length
				1, 0, 0, 0, // IAID
				0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime

				0, 4, // IATA option code
				0, 32, // length
				1, 2, 3, 4, // IAID
				0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, // IPv6
				0, 0, 0, 2, // PreferredLifetime
				0, 0, 0, 4, // ValidLifetime
			},
			want: []*OptIATA{
				&OptIATA{
					IaId: [4]byte{1, 0, 0, 0},
					Options: IdentityOptions{Options: Options{&OptIAAddress{
						IPv6Addr:          net.IP{0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0},
						PreferredLifetime: 2 * time.Second,
						ValidLifetime:     4 * time.Second,
						Options:           AddressOptions{Options: Options{}},
					}}},
				},
				&OptIATA{
					IaId: [4]byte{1, 2, 3, 4},
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
			buf:  []byte{0, 4, 0, 1, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
		{
			buf: []byte{
				0, 4, // IATA option code
				0, 3, // length
				1, 0, 0, // IAID too short
			},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
		{
			buf: []byte{
				0, 4, // IATA option code
				0, 28, // length
				1, 0, 0, 0, // IAID
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
			if got := mo.IATA(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IATA = %v, want %v", got, tt.want)
			}
			var wantOne *OptIATA
			if len(tt.want) >= 1 {
				wantOne = tt.want[0]
			}
			if got := mo.OneIATA(); !reflect.DeepEqual(got, wantOne) {
				t.Errorf("OneIATA = %v, want %v", got, wantOne)
			}

			if len(tt.want) >= 1 {
				var b MessageOptions
				for _, iata := range tt.want {
					b.Add(iata)
				}
				got := b.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptIATAGetOneOption(t *testing.T) {
	oaddr := &OptIAAddress{
		IPv6Addr: net.ParseIP("::1"),
	}
	opt := OptIATA{
		Options: IdentityOptions{[]Option{&OptStatusCode{}, oaddr}},
	}
	require.Equal(t, oaddr, opt.Options.OneAddress())
}

func TestOptIATAAddOption(t *testing.T) {
	opt := OptIATA{}
	opt.Options.Add(OptElapsedTime(0))
	require.Equal(t, 1, len(opt.Options.Options))
	require.Equal(t, OptionElapsedTime, opt.Options.Options[0].Code())
}

func TestOptIATAGetOneOptionMissingOpt(t *testing.T) {
	oaddr := &OptIAAddress{
		IPv6Addr: net.ParseIP("::1"),
	}
	opt := OptIATA{
		Options: IdentityOptions{[]Option{&OptStatusCode{}, oaddr}},
	}
	require.Equal(t, nil, opt.Options.GetOne(OptionDNSRecursiveNameServer))
}

func TestOptIATADelOption(t *testing.T) {
	optiaaddr := OptIAAddress{}
	optsc := OptStatusCode{}

	iana1 := OptIATA{
		Options: IdentityOptions{[]Option{
			&optsc,
			&optiaaddr,
			&optiaaddr,
		}},
	}
	iana1.Options.Del(OptionIAAddr)
	require.Equal(t, iana1.Options.Options, Options{&optsc})

	iana2 := OptIATA{
		Options: IdentityOptions{[]Option{
			&optiaaddr,
			&optsc,
			&optiaaddr,
		}},
	}
	iana2.Options.Del(OptionIAAddr)
	require.Equal(t, iana2.Options.Options, Options{&optsc})
}

func TestOptIATAString(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, 0, 0, 0xb2, 0x7a, 0, 0, 0xc0, 0x8a, // options
	}
	var opt OptIATA
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
		"Options={",
		"String() should return a list of options",
	)
}
