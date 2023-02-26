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

func TestArchTypeParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want iana.Archs
	}{
		{
			buf: []byte{
				0, 61, // Client Arch Types option
				0, 2, // length
				0, 7, // EFI_X86_64
			},
			want: iana.Archs{iana.EFI_X86_64},
		},
		{
			buf: []byte{
				0, 61, // Client Arch Types option
				0, 4, // length
				0, 7, // EFI_X86_64
				0, 8, // EFI_XSCALE
			},
			want: iana.Archs{iana.EFI_X86_64, iana.EFI_XSCALE},
		},
		{
			buf:  nil,
			want: nil,
		},
		{
			buf:  []byte{0, 61, 0, 1, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
		{
			buf:  []byte{0, 61, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.ArchTypes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchTypes = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				var m MessageOptions
				m.Add(OptClientArchType(tt.want...))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptClientArchType(t *testing.T) {
	opt := OptClientArchType(iana.EFI_ITANIUM)
	require.Equal(t, OptionClientArchType, opt.Code())
	require.Contains(t, opt.String(), "EFI Itanium", "String() should contain the correct ArchType output")
}
