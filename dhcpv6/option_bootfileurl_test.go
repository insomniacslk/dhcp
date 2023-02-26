package dhcpv6

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func TestBootFileURLParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want string
	}{
		{
			buf: []byte{
				0, 59, // Boot File URL
				0, 17, // length
				'h', 't', 't', 'p', ':', '/', '/', 'u', '-', 'r', 'o', 'o', 't', '.', 'o', 'r', 'g',
			},
			want: "http://u-root.org",
		},
		{
			buf: nil,
		},
		{
			buf: []byte{0, 59, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.BootFileURL(); got != tt.want {
				t.Errorf("BootFileURL = %v, want %v", got, tt.want)
			}

			if tt.want != "" {
				var m MessageOptions
				m.Add(OptBootFileURL(tt.want))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptBootFileURL(t *testing.T) {
	opt := OptBootFileURL("https://insomniac.slackware.it")
	require.Contains(t, opt.String(), "https://insomniac.slackware.it", "String() should contain the correct BootFileUrl output")
}
