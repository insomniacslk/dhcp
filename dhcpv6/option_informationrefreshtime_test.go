package dhcpv6

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/uio/uio"
)

func TestInformationRefreshTimeParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want time.Duration
	}{
		{
			buf: []byte{
				0, 32, // IRT option
				0, 4, // length
				0, 0, 0, 3,
			},
			want: 3 * time.Second,
		},
		{
			buf: []byte{
				0, 32, // IRT option
				0, 6, // length
				0, 0, 0, 3, 0, 0,
			},
			err: uio.ErrUnreadBytes,
		},
		{
			buf: []byte{0, 32, 0, 1, 0},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: []byte{0, 32, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.InformationRefreshTime(0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InformationRefreshTime = %v, want %v", got, tt.want)
			}

			if tt.err == nil {
				var m MessageOptions
				m.Add(OptInformationRefreshTime(tt.want))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptInformationRefreshTime(t *testing.T) {
	var opt optInformationRefreshTime
	err := opt.FromBytes([]byte{0xaa, 0xbb, 0xcc, 0xdd})
	if err != nil {
		t.Fatal(err)
	}
	if informationRefreshTime := opt.InformationRefreshtime; informationRefreshTime != time.Duration(0xaabbccdd)*time.Second {
		t.Fatalf("Invalid information refresh time. Expected 0xaabb, got %v", informationRefreshTime)
	}
}

func TestOptInformationRefreshTimeToBytes(t *testing.T) {
	opt := OptInformationRefreshTime(0)
	expected := []byte{0, 0, 0, 0}
	if toBytes := opt.ToBytes(); !bytes.Equal(expected, toBytes) {
		t.Fatalf("Invalid ToBytes output. Expected %v, got %v", expected, toBytes)
	}
}

func TestOptInformationRefreshTimeString(t *testing.T) {
	opt := OptInformationRefreshTime(3600 * time.Second)
	expected := "Information Refresh Time: 1h0m0s"
	if optString := opt.String(); optString != expected {
		t.Fatalf("Invalid elapsed time string. Expected %v, got %v", expected, optString)
	}
}
