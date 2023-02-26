package dhcpv6

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func TestStatusCodeParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want *OptStatusCode
	}{
		{
			buf: []byte{
				0, 13, // StatusCode option
				0, 15, // length
				0, 5, // StatusUseMulticast
				'u', 's', 'e', ' ', 'm', 'u', 'l', 't', 'i', 'c', 'a', 's', 't',
			},
			want: &OptStatusCode{
				StatusCode:    iana.StatusUseMulticast,
				StatusMessage: "use multicast",
			},
		},
		{
			buf:  nil,
			want: nil,
		},
		{
			buf:  []byte{0, 13, 0, 1, 0},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
		{
			buf:  []byte{0, 13, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.Status(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Status = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				var m MessageOptions
				m.Add(tt.want)
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptStatusCodeString(t *testing.T) {
	data := []byte{
		0, 5, // StatusUseMulticast
		'u', 's', 'e', ' ', 'm', 'u', 'l', 't', 'i', 'c', 'a', 's', 't',
	}
	var opt OptStatusCode
	err := opt.FromBytes(data)
	require.NoError(t, err)

	require.Contains(
		t,
		opt.String(),
		"Code=UseMulticast (5); Message=use multicast",
		"String() should contain the code and message",
	)
}
