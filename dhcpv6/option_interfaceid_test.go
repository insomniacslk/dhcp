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

func TestInterfaceIDParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []byte
	}{
		{
			buf: []byte{
				0, 18, // Interface ID
				0, 4, // length
				'S', 'L', 'A', 'M',
			},
			want: []byte("SLAM"),
		},
		{
			buf: []byte{
				0, 18,
				0, 0,
			},
		},
		{
			buf: []byte{0, 18, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var ro RelayOptions
			if err := ro.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := ro.InterfaceID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterfaceID = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				var m RelayOptions
				m.Add(OptInterfaceID(tt.want))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptInterfaceID(t *testing.T) {
	opt := OptInterfaceID([]byte("DSLAM01 eth2/1/01/21"))
	require.Contains(
		t,
		opt.String(),
		"68 83 76 65 77 48 49 32 101 116 104 50 47 49 47 48 49 47 50 49",
		"String() should return the interfaceId as bytes",
	)
}
