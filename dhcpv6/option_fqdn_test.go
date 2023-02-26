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

func TestFQDNParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want *OptFQDN
	}{
		{
			buf: []byte{
				0, 39, // FQDN option
				0, 34, // length
				0, // flags
				7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
				6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
			},
			want: &OptFQDN{
				Flags: 0,
				DomainName: &rfc1035label.Labels{
					Labels: []string{
						"example.com",
						"subnet.example.org",
					},
				},
			},
		},
		{
			buf: []byte{
				0, 39, // FQDN
				0, 23, // length
				0, // flags
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
			buf:  []byte{0, 39, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			got := mo.FQDN()
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(rfc1035label.Labels{})) {
				t.Errorf("FQDN = %v, want %v", got, tt.want)
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

func TestOptFQDNString(t *testing.T) {
	opt := &OptFQDN{
		DomainName: &rfc1035label.Labels{
			Labels: []string{"cnos.localhost"},
		},
	}
	require.Equal(t, "FQDN: {Flags=0 DomainName=[cnos.localhost]}", opt.String())
}
