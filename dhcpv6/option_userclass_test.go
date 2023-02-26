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

func TestUserClassParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want [][]byte
	}{
		{
			buf: joinBytes([]byte{
				0, 15, // User Class
				0, 19, // length
				0, 8,
			}, []byte("bladibla"), []byte{0, 7}, []byte("foo=bar")),
			want: [][]byte{[]byte("bladibla"), []byte("foo=bar")},
		},
		{
			buf: nil,
		},
		{
			buf: []byte{
				0, 15,
				0, 0,
			},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: []byte{0, 15, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.UserClasses(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserClass = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				var m MessageOptions
				m.Add(&OptUserClass{UserClasses: tt.want})
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptUserClassString(t *testing.T) {
	data := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e', 's', 't',
	}
	var opt OptUserClass
	err := opt.FromBytes(data)
	require.NoError(t, err)

	require.Contains(
		t,
		opt.String(),
		"User Class: [linuxboot, test]",
		"String() should contain the list of user classes",
	)
}
