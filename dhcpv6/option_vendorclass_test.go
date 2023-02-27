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

func TestVendorClassParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []*OptVendorClass
	}{
		{
			buf: []byte{
				0, 16, // Vendor Class
				0, 14, // length
				0, 0, 0, 16,
				0, 4,
				'S', 'L', 'A', 'M',
				0, 2,
				'h', 'h',
			},
			want: []*OptVendorClass{
				&OptVendorClass{
					EnterpriseNumber: 16,
					Data:             [][]byte{[]byte("SLAM"), []byte("hh")},
				},
			},
		},
		{
			buf: []byte{
				0, 16,
				0, 0,
			},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 16,
				0, 4,
				0, 0, 0, 6,
			},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: []byte{0, 16, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.VendorClasses(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VendorClass = %v, want %v", got, tt.want)
			}
			for _, v := range tt.want {
				if got := mo.VendorClass(v.EnterpriseNumber); !reflect.DeepEqual(got, v.Data) {
					t.Errorf("VendorClass(%d) = %v, want %v", v.EnterpriseNumber, got, v.Data)
				}
			}
			if got := mo.VendorClass(100); got != nil {
				t.Errorf("VendorClass(100) = %v, want nil", got)
			}

			if tt.want != nil {
				var m MessageOptions
				for _, o := range tt.want {
					m.Add(o)
				}
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptVendorClassString(t *testing.T) {
	data := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // EnterpriseNumber
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e', 's', 't',
	}
	var opt OptVendorClass
	err := opt.FromBytes(data)
	require.NoError(t, err)

	str := opt.String()
	require.Contains(
		t,
		str,
		"EnterpriseNumber=2864434397",
		"String() should contain the enterprisenum",
	)
	require.Contains(
		t,
		str,
		"Data=[linuxboot, test]",
		"String() should contain the list of vendor classes",
	)
}
