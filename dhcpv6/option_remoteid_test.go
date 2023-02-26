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

func TestRemoteIDParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want *OptRemoteID
	}{
		{
			buf: []byte{
				0, 37, // Remote ID
				0, 8, // length
				0, 0, 0, 16,
				'S', 'L', 'A', 'M',
			},
			want: &OptRemoteID{
				EnterpriseNumber: 16,
				RemoteID:         []byte("SLAM"),
			},
		},
		{
			buf: []byte{
				0, 37,
				0, 0,
			},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 37,
				0, 4,
				0, 0, 0, 6,
			},
			want: &OptRemoteID{
				EnterpriseNumber: 6,
				RemoteID:         []byte{},
			},
		},
		{
			buf: []byte{0, 37, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var ro RelayOptions
			if err := ro.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := ro.RemoteID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteID = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				var m RelayOptions
				m.Add(tt.want)
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptRemoteIDString(t *testing.T) {
	opt := &OptRemoteID{
		EnterpriseNumber: 123,
		RemoteID:         []byte("Test1234"),
	}
	str := opt.String()
	require.Contains(
		t,
		str,
		"EnterpriseNumber=123",
		"String() should contain the enterprisenum",
	)
	require.Contains(
		t,
		str,
		"RemoteID=0x5465737431323334",
		"String() should contain the remoteid bytes",
	)
}
