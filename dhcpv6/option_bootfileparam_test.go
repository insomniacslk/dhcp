package dhcpv6

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/uio/uio"
)

func TestBootFileParamLargeParameter(t *testing.T) {
	param := []string{
		"foo=bar",
		strings.Repeat("a", 1<<16),
	}
	var m MessageOptions
	m.Add(OptBootFileParam(param...))
	want := append([]byte{
		0, 60, // Boot File Param
		0, 9, // length
		0, 7,
	}, []byte("foo=bar")...)

	got := m.ToBytes()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
	}
}

func joinBytes(p ...[]byte) []byte {
	return bytes.Join(p, nil)
}

func TestBootFileParamParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []string
	}{
		{
			buf: joinBytes([]byte{
				0, 60, // Boot File Param
				0, 25, // length
				0, 14,
			}, []byte("root=/dev/sda1"), []byte{0, 7}, []byte("foo=bar")),
			want: []string{"root=/dev/sda1", "foo=bar"},
		},
		{
			buf: nil,
		},
		{
			buf: []byte{0, 60, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.BootFileParam(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BootFileParam = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				var m MessageOptions
				m.Add(OptBootFileParam(tt.want...))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}
