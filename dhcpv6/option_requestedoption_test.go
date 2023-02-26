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

func TestOROParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want OptionCodes
	}{
		{
			buf: []byte{
				0, 6, // ORO option
				0, 2, // length
				0, 3, // IANA option
			},
			want: OptionCodes{OptionIANA},
		},
		{
			buf: []byte{
				0, 6, // ORO option
				0, 4, // length
				0, 3, // IANA
				0, 4, // IATA
			},
			want: OptionCodes{OptionIANA, OptionIATA},
		},
		{
			buf:  nil,
			want: nil,
		},
		{
			buf:  []byte{0, 6, 0, 1, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
		{
			buf:  []byte{0, 6, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.RequestedOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequestedOptions = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				var m MessageOptions
				m.Add(OptRequestedOption(tt.want...))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestParseMessageOptionsWithORO(t *testing.T) {
	buf := []byte{
		0, 6, // ORO option
		0, 2, // length
		0, 3, // IANA Option
		0, 6, // ORO
		0, 2, // length
		0, 4, // IATA
	}

	want := OptionCodes{OptionIANA, OptionIATA}
	var mo MessageOptions
	if err := mo.FromBytes(buf); err != nil {
		t.Errorf("FromBytes = %v", err)
	} else if got := mo.RequestedOptions(); !reflect.DeepEqual(got, want) {
		t.Errorf("ORO = %v, want %v", got, want)
	}
}

func TestOptRequestedOptionString(t *testing.T) {
	buf := []byte{0, 1, 0, 2}
	var o optRequestedOption
	err := o.FromBytes(buf)
	require.NoError(t, err)
	require.Contains(
		t,
		o.String(),
		"Client ID, Server ID",
		"String() should contain the options specified",
	)
	o.OptionCodes = append(o.OptionCodes, 12345)
	require.Contains(
		t,
		o.String(),
		"unknown",
		"String() should contain 'Unknown' for an illegal option",
	)
}
