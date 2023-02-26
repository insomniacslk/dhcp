package dhcpv6

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func TestDHCP4oDHCP6ParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want *OptDHCP4oDHCP6Server
	}{
		{
			buf: []byte{
				0, 88, // DHCP4oDHCP6 option.
				0, 32, // length
				0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35,
				0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35,
			},
			want: &OptDHCP4oDHCP6Server{
				DHCP4oDHCP6Servers: []net.IP{
					net.IP{0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35},
					net.IP{0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35},
				},
			},
		},
		{
			buf: []byte{
				0, 88, // DHCP4oDHCP6 option.
				0, 6, // length
				0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe,
			},
			err: uio.ErrUnreadBytes,
		},
		{
			buf: []byte{
				0, 88, // DHCP4oDHCP6 option.
				0, 0, // length
			},
			want: &OptDHCP4oDHCP6Server{},
		},
		{
			buf:  []byte{0, 88, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.DHCP4oDHCP6Server(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DHCP4oDHCP6Server = %v, want %v", got, tt.want)
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

func TestParseOptDHCP4oDHCP6Server(t *testing.T) {
	opt := OptDHCP4oDHCP6Server{DHCP4oDHCP6Servers: []net.IP{net.ParseIP("2a03:2880:fffe:c:face:b00c:0:35")}}
	require.Contains(t, opt.String(), "[2a03:2880:fffe:c:face:b00c:0:35]", "String() should contain the correct DHCP4-over-DHCP6 server output")
}
