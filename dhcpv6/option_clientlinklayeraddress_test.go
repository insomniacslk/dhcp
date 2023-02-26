package dhcpv6

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func TestClientLinkLayerAddressParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf        []byte
		err        error
		wantHWType iana.HWType
		wantHWAddr net.HardwareAddr
	}{
		{
			buf: []byte{
				0, 79, // Client Link Layer Address option
				0, 8, // length
				0, 1, // Ethernet
				1, 2, 3, 4, 5, 6,
			},
			wantHWType: iana.HWTypeEthernet,
			wantHWAddr: net.HardwareAddr{1, 2, 3, 4, 5, 6},
		},
		{
			buf: []byte{0, 79, 0, 1, 0},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: []byte{0, 79, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var ro RelayOptions
			if err := ro.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if gotHWType, gotHWAddr := ro.ClientLinkLayerAddress(); gotHWType != tt.wantHWType || !bytes.Equal(gotHWAddr, tt.wantHWAddr) {
				t.Errorf("ClientLinkLayerAddress = (%s, %v), want (%s, %v)", gotHWType, tt.wantHWType, gotHWAddr, tt.wantHWAddr)
			}

			if tt.err == nil {
				var m MessageOptions
				m.Add(OptClientLinkLayerAddress(tt.wantHWType, tt.wantHWAddr))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptClientLinkLayerAddressString(t *testing.T) {
	opt := OptClientLinkLayerAddress(iana.HWTypeEthernet, net.HardwareAddr{0xa4, 0x83, 0xe7, 0xe3, 0xdf, 0x88})
	require.Equal(t, "Client Link-Layer Address: Type=Ethernet LinkLayerAddress=a4:83:e7:e3:df:88", opt.String())
}
