package dhcpv6

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/uio/uio"
)

func TestVendorOptsParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []*OptVendorOpts
	}{
		{
			buf: []byte{
				0, 17, // VendorOpts option
				0, 10, // length
				0, 0, 0, 16,
				0, 5, // type
				0, 2, // length
				0xa, 0xb,

				0, 17, // VendorOpts option
				0, 9, // length
				0, 0, 0, 14,
				0, 9, // type
				0, 1, // length
				0xa,
			},
			want: []*OptVendorOpts{
				&OptVendorOpts{
					EnterpriseNumber: 16,
					VendorOpts: Options{
						&OptionGeneric{OptionCode: 5, OptionData: []byte{0xa, 0xb}},
					},
				},
				&OptVendorOpts{
					EnterpriseNumber: 14,
					VendorOpts: Options{
						&OptionGeneric{OptionCode: 9, OptionData: []byte{0xa}},
					},
				},
			},
		},
		{
			buf:  nil,
			want: nil,
		},
		{
			buf:  []byte{0, 17, 0, 1, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
		{
			buf:  []byte{0, 17, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.VendorOpts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VendorOpts = %v, want %v", got, tt.want)
			}
			for _, vo := range tt.want {
				if got := mo.VendorOpt(vo.EnterpriseNumber); !reflect.DeepEqual(got, vo.VendorOpts) {
					t.Errorf("VendorOpt(%d) = %v, want %v", vo.EnterpriseNumber, got, vo.VendorOpts)
				}
			}
			if got := mo.VendorOpt(100); got != nil {
				t.Errorf("VendorOpt(100) = %v, not nil", got)
			}

			if tt.want != nil {
				var m MessageOptions
				for _, opt := range tt.want {
					m.Add(opt)
				}
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}
