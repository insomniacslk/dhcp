package dhcpv6

import (
	"bytes"
	"net"
	"reflect"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestDuidInvalidTooShort(t *testing.T) {
	// too short DUID at all (must be at least 2 bytes)
	_, err := DUIDFromBytes([]byte{0})
	require.Error(t, err)

	// too short DUID_LL (must be at least 4 bytes)
	_, err = DUIDFromBytes([]byte{0, 3, 0xa})
	require.Error(t, err)

	// too short DUID_EN (must be at least 6 bytes)
	_, err = DUIDFromBytes([]byte{0, 2, 0xa, 0xb, 0xc})
	require.Error(t, err)

	// too short DUID_LLT (must be at least 8 bytes)
	_, err = DUIDFromBytes([]byte{0, 1, 0xa, 0xb, 0xc, 0xd, 0xe})
	require.Error(t, err)

	// too short DUID_UUID (must be at least 18 bytes)
	_, err = DUIDFromBytes([]byte{0, 4, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf})
	require.Error(t, err)
}

func TestFromBytes(t *testing.T) {
	for _, tt := range []struct {
		name     string
		buf      []byte
		want     DUID
		stringer string
	}{
		{
			name: "DUID-LLT",
			buf: []byte{
				0, 1, // DUID_LLT
				0, 1, // HwTypeEthernet
				0x01, 0x02, 0x03, 0x04, // time
				0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, // link-layer addr
			},
			want: &DUIDLLT{
				Time:          0x01020304,
				HWType:        iana.HWTypeEthernet,
				LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			},
			stringer: "DUID-LLT{HWType=Ethernet HWAddr=aa:bb:cc:dd:ee:ff Time=16909060}",
		},
		{
			name: "DUID-LL",
			buf: []byte{
				0, 3, // DUID_LL
				0, 1, // HwTypeEthernet
				0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, // link-layer addr
			},
			want: &DUIDLL{
				HWType:        iana.HWTypeEthernet,
				LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			},
			stringer: "DUID-LL{HWType=Ethernet HWAddr=aa:bb:cc:dd:ee:ff}",
		},
		{
			name: "DUID-EN",
			buf: []byte{
				0, 2, // DUID_EN
				0, 0, 0, 1, // EnterpriseNumber
				0x66, 0x6f, 0x6f, // "foo"
			},
			want: &DUIDEN{
				EnterpriseNumber:     0x1,
				EnterpriseIdentifier: []byte("foo"),
			},
			stringer: "DUID-EN{EnterpriseNumber=1 EnterpriseIdentifier=foo}",
		},
		{
			name: "DUID-UUID",
			buf: []byte{
				0x00, 0x04, // DUID_UUID
				0x01, 0x02, 0x03, 0x04, // UUID
				0x01, 0x02, 0x03, 0x04, // UUID
				0x01, 0x02, 0x03, 0x04, // UUID
				0x01, 0x02, 0x03, 0x04, // UUID
			},
			want: &DUIDUUID{
				UUID: [16]byte{
					0x01, 0x02, 0x03, 0x04,
					0x01, 0x02, 0x03, 0x04,
					0x01, 0x02, 0x03, 0x04,
					0x01, 0x02, 0x03, 0x04,
				},
			},
			stringer: "DUID-UUID{0x01020304010203040102030401020304}",
		},
		{
			name: "DUIDOpaque",
			buf: []byte{
				0x00, 0x05, // unknown DUID
				0x01, 0x02, 0x03, // Opaque
			},
			want: &DUIDOpaque{
				Type: 0x5,
				Data: []byte{0x01, 0x02, 0x03},
			},
			stringer: "DUID-Opaque{Type=5 Data=0x010203}",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// FromBytes
			got, err := DUIDFromBytes(tt.buf)
			if err != nil {
				t.Errorf("DUIDFromBytes = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DUIDFromBytes = %v, want %v", got, tt.want)
			}

			// ToBytes
			buf := tt.want.ToBytes()
			if !bytes.Equal(buf, tt.buf) {
				t.Errorf("ToBytes() = %#x, want %#x", buf, tt.buf)
			}

			// Stringer
			s := tt.want.String()
			if s != tt.stringer {
				t.Errorf("String() = %s, want %s", s, tt.stringer)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	for _, tt := range []struct {
		name string
		a    DUID
		b    DUID
		want bool
	}{
		{
			name: "DUID-LL-equal",
			a:    &DUIDLL{HWType: iana.HWTypeEthernet, LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
			b:    &DUIDLL{HWType: iana.HWTypeEthernet, LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
			want: true,
		},
		{
			name: "DUID-LL-not-equal",
			a:    &DUIDLL{HWType: iana.HWTypeEthernet, LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
			b:    &DUIDLL{HWType: iana.HWTypeEthernet, LinkLayerAddr: net.HardwareAddr{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}},
			want: false,
		},
		{
			name: "DUID-LL-and-DUID-EN",
			a:    &DUIDLL{HWType: iana.HWTypeEthernet, LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
			b:    &DUIDEN{EnterpriseNumber: 5, EnterpriseIdentifier: []byte("foo")},
			want: false,
		},
		{
			name: "DUID-EN-equal",
			a:    &DUIDEN{EnterpriseNumber: 5, EnterpriseIdentifier: []byte("foo")},
			b:    &DUIDEN{EnterpriseNumber: 5, EnterpriseIdentifier: []byte("foo")},
			want: true,
		},
		{
			name: "DUID-EN-not-equal",
			a:    &DUIDEN{EnterpriseNumber: 5, EnterpriseIdentifier: []byte("foo")},
			b:    &DUIDEN{EnterpriseNumber: 5, EnterpriseIdentifier: []byte("bar")},
			want: false,
		},
		{
			name: "DUID-LLT-equal",
			a:    &DUIDLLT{HWType: iana.HWTypeEthernet, Time: 10, LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
			b:    &DUIDLLT{HWType: iana.HWTypeEthernet, Time: 10, LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
			want: true,
		},
		{
			name: "DUID-LLT-not-equal",
			a:    &DUIDLLT{HWType: iana.HWTypeEthernet, Time: 10, LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
			b:    &DUIDLLT{HWType: iana.HWTypeEthernet, Time: 10, LinkLayerAddr: net.HardwareAddr{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}},
			want: false,
		},
		{
			name: "DUID-LLT-and-DUID-UUID",
			a:    &DUIDLLT{HWType: iana.HWTypeEthernet, Time: 10, LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
			b:    &DUIDUUID{UUID: [16]byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}},
			want: false,
		},
		{
			name: "DUID-UUID-equal",
			a:    &DUIDUUID{UUID: [16]byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}},
			b:    &DUIDUUID{UUID: [16]byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}},
			want: true,
		},
		{
			name: "DUID-UUID-not-equal",
			a:    &DUIDUUID{UUID: [16]byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}},
			b:    &DUIDUUID{UUID: [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
			want: false,
		},
		{
			name: "DUID-UUID-and-DUID-Opaque",
			a:    &DUIDOpaque{Type: 5, Data: []byte{0x1}},
			b:    &DUIDUUID{UUID: [16]byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}},
			want: false,
		},
		{
			name: "DUID-Opaque-equal",
			a:    &DUIDOpaque{Type: 5, Data: []byte{0x1}},
			b:    &DUIDOpaque{Type: 5, Data: []byte{0x1}},
			want: true,
		},
		{
			name: "DUID-Opaque-not-equal",
			a:    &DUIDOpaque{Type: 5, Data: []byte{0x1}},
			b:    &DUIDOpaque{Type: 5, Data: []byte{0x2}},
			want: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(tt.b); got != tt.want {
				t.Errorf("%s.Equal(%s) = %v, want %v", tt.a, tt.b, got, tt.want)
			}

			if got := tt.b.Equal(tt.a); got != tt.want {
				t.Errorf("%s.Equal(%s) = %v, want %v", tt.b, tt.a, got, tt.want)
			}
		})
	}
}
