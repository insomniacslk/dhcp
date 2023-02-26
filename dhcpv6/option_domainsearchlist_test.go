package dhcpv6

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/insomniacslk/dhcp/rfc1035label"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func TestDomainSearchListParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want *rfc1035label.Labels
	}{
		{
			buf: []byte{
				0, 24, // Domain Search List option
				0, 33, // length
				7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
				6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
			},
			want: &rfc1035label.Labels{
				Labels: []string{
					"example.com",
					"subnet.example.org",
				},
			},
		},
		{
			buf: []byte{
				0, 24, // Domain Search List option
				0, 22, // length
				7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
				6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', // truncated
			},
			err: rfc1035label.ErrBufferTooShort,
		},
		{
			buf:  nil,
			want: nil,
		},
		{
			buf:  []byte{0, 24, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			got := mo.DomainSearchList()
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(rfc1035label.Labels{})) {
				t.Errorf("DomainSearchList = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				var m MessageOptions
				m.Add(OptDomainSearchList(tt.want))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptDomainSearchListString(t *testing.T) {
	opt := OptDomainSearchList(&rfc1035label.Labels{
		Labels: []string{
			"example.com",
			"subnet.example.org",
		},
	})
	require.Contains(t, opt.String(), "example.com subnet.example.org", "String() should contain the correct domain search output")
}
