package dhcpv6

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func TestServerIDParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want DUID
	}{
		{
			buf: []byte{
				0, 2, // Server ID option
				0, 10, // length
				0, 3, // DUID_LL
				0, 1, // hwtype ethernet
				0, 1, 2, 3, 4, 5, // HW addr
			},
			want: &DUIDLL{HWType: iana.HWTypeEthernet, LinkLayerAddr: net.HardwareAddr{0, 1, 2, 3, 4, 5}},
		},
		{
			buf:  nil,
			want: nil,
		},
		{
			buf:  []byte{0, 1, 0, 1, 0},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.ServerID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerID = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerID(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		want optServerID
		err  error
	}{
		{
			buf: []byte{
				0, 3, // DUID_LL
				0, 1, // hwtype ethernet
				0, 1, 2, 3, 4, 5, // hw addr
			},
			want: optServerID{
				&DUIDLL{
					HWType:        iana.HWTypeEthernet,
					LinkLayerAddr: net.HardwareAddr([]byte{0, 1, 2, 3, 4, 5}),
				},
			},
		},
		{
			buf: []byte{0},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: []byte{0, 3, 0},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: nil,
			err: uio.ErrBufferTooShort,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var opt optServerID
			if err := opt.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if tt.err == nil {
				if !reflect.DeepEqual(tt.want, opt) {
					t.Errorf("FromBytes = %v, want %v", opt, tt.want)
				}

				out := tt.want.ToBytes()
				if diff := cmp.Diff(tt.buf, out); diff != "" {
					t.Errorf("ToBytes mismatch: (-want, +got):\n%s", diff)
				}
			}
		})
	}
}

func TestOptionServerIDString(t *testing.T) {
	opt := OptServerID(
		&DUIDLL{
			HWType:        iana.HWTypeEthernet,
			LinkLayerAddr: net.HardwareAddr([]byte{0xde, 0xad, 0, 0, 0xbe, 0xef}),
		},
	)
	require.Equal(t, OptionServerID, opt.Code())
	require.Contains(
		t,
		opt.String(),
		"Server ID: DUID-LL{HWType=Ethernet HWAddr=de:ad:00:00:be:ef}",
		"String() should contain the correct cid output",
	)
}
