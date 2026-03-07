package dhcpv6

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func TestParseRelayPort(t *testing.T) {
	var opt optRelayPort
	err := opt.FromBytes([]byte{0x12, 0x32})
	require.NoError(t, err)
	require.Equal(t, optRelayPort{DownstreamSourcePort: 0x1232}, opt)
}

func TestRelayPortParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want *optRelayPort
	}{
		{
			buf: []byte{
				0, 135, // Relay Port
				0, 2, // length
				0x12, 0x34,
			},
			want: &optRelayPort{ DownstreamSourcePort: 0x1234 },
		},
		{
			buf: []byte{
				0, 135,
				0, 2,
				0, 0,
			},
			want: &optRelayPort{ DownstreamSourcePort: 0 },
		},
		{
			buf: []byte{
				0, 135,
				0, 0,
			},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 135,
				0,
			},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var ro RelayOptions
			if err := ro.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := ro.RelayPort(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RelayPort = %v want %v", got, tt.want)
			}

			if tt.want != nil {
				var m RelayOptions
				m.Add(OptRelayPort(tt.want.DownstreamSourcePort))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestRelayPortToBytes(t *testing.T) {
	op := OptRelayPort(0x3845)
	require.Equal(t, []byte{0x38, 0x45}, op.ToBytes())
}
