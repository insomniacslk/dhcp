package dhcpv6

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestPreferenceParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want uint8
	}{
		{
			buf: []byte{
				0, 7, // Preference option
				0, 1, // length
				25,
			},
			want: 25,
		},
		{
			buf:  nil,
			want: 0,
		},
		{
			buf: []byte{
				0, 7, // Preference option
				0, 0, // length
			},
			want: 0,
			err:  errBufferNotLengthOne,
		},
		{
			buf: []byte{
				0, 7, // Preference option
				0, 2, // length
				0, 0,
			},
			want: 0,
			err:  errBufferNotLengthOne,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.Preference(); got != tt.want {
				t.Errorf("Preference = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreference(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		want optPreference
		err  error
	}{
		{
			buf:  []byte{0},
			want: optPreference{},
		},
		{
			buf: []byte{5},
			want: optPreference{
				prefValue: 5,
			},
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var opt optPreference
			if err := opt.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if tt.err == nil {
				if !reflect.DeepEqual(tt.want, opt) {
					t.Errorf("FromBytes = %v, want %v", opt, tt.want)
				}

				out := tt.want.ToBytes()
				if diff := cmp.Diff(tt.buf, out); diff != "" {
					t.Errorf("ToBytes mismatch: (-want, +got):\n%s", diff)
				}
			}
		})
	}
}

func TestOptionPreferenceString(t *testing.T) {
	opt := OptPreference(12)
	require.Equal(t, OptionPreference, opt.Code())
	require.Contains(
		t,
		opt.String(),
		"Preference: 12",
		"String() should contain the correct output",
	)
}
