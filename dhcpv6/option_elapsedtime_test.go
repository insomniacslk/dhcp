package dhcpv6

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/uio/uio"
)

func TestElapsedTimeParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want time.Duration
	}{
		{
			buf: []byte{
				0, 8, // Elapsed Time option
				0, 2, // length
				0, 2,
			},
			want: 20 * time.Millisecond,
		},
		{
			buf: []byte{
				0, 8, // Elapsed Time option
				0, 1, // length
				0,
			},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 8, // Elapsed Time option
				0, 3, // length
				0, 2, 2,
			},
			err: uio.ErrUnreadBytes,
		},
		{
			buf: []byte{0, 8, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.ElapsedTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ElapsedTime = %v, want %v", got, tt.want)
			}

			if tt.err == nil {
				var m MessageOptions
				m.Add(OptElapsedTime(tt.want))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptElapsedTimeString(t *testing.T) {
	opt := OptElapsedTime(100 * time.Millisecond)
	expected := "Elapsed Time: 100ms"
	if optString := opt.String(); optString != expected {
		t.Fatalf("Invalid elapsed time string. Expected %v, got %v", expected, optString)
	}
}
