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

func TestRelayMsgOptionParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want DHCPv6
	}{
		{
			buf: []byte{
				0, 9, // Relay Msg option
				0, 10, // length
				1,                // MessageTypeSolicit
				0xaa, 0xbb, 0xcc, // transaction ID
				0, 8, // option: elapsed time
				0, 2, // option length
				0, 0, // option value
			},
			want: &Message{
				MessageType:   MessageTypeSolicit,
				TransactionID: TransactionID{0xaa, 0xbb, 0xcc},
				Options:       MessageOptions{Options{OptElapsedTime(0)}},
			},
		},
		{
			buf: []byte{
				0, 9, // Relay Msg option
				0, 6, // length
				1,                // MessageTypeSolicit
				0xaa, 0xbb, 0xcc, // transaction ID
				0, 8, // option: elapsed time
				// no length/value for elapsed time option
			},
			err: uio.ErrUnreadBytes,
		},
		{
			buf:  []byte{0, 9, 0, 1, 0},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
		{
			buf:  []byte{0, 9, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var ro RelayOptions
			if err := ro.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := ro.RelayMessage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RelayMessage = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				var m RelayOptions
				m.Add(OptRelayMessage(tt.want))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestSample(t *testing.T) {
	// Nested relay message. This test only checks if it parses correctly, but
	// could/should be extended to check all the fields like done in other tests
	buf := []byte{
		12,                                             // relay
		1,                                              // hop count
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // linkAddr
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // peerAddr
		// relay msg
		0, 9, // opt relay msg
		0, 66, // opt len
		// relay fwd
		12,                                             // relay
		0,                                              // hop count
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // linkAddr
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // peerAddr
		// opt interface ID
		0, 18, // interface id
		0, 6, // opt len
		0xba, 0xbe, 0xb1, 0xb0, 0xbe, 0xbe, // opt value
		// relay msg
		0, 9, // relay msg
		0, 18, // msg len
		// dhcpv6 msg
		1,                // solicit
		0xaa, 0xbb, 0xcc, // transaction ID
		// client ID
		0, 1, // opt client id
		0, 10, // opt len
		0, 3, // duid type
		0, 1, // hw type
		5, 6, 7, 8, 9, 10, // duid value
	}
	_, err := FromBytes(buf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRelayMsgParseOptRelayMsgTooShort(t *testing.T) {
	var opt optRelayMsg
	err := opt.FromBytes([]byte{})
	require.Error(t, err, "ParseOptRelayMsg() should return an error if the encapsulated message is invalid")
}

func TestRelayMsgString(t *testing.T) {
	var opt optRelayMsg
	err := opt.FromBytes([]byte{
		1,                // MessageTypeSolicit
		0xaa, 0xbb, 0xcc, // transaction ID
		0, 8, // option: elapsed time
		0, 2, // option length
		0, 0, // option value
	})
	require.NoError(t, err)
	require.Contains(
		t,
		opt.String(),
		"Relay Message: Message",
		"String() should contain the relaymsg contents",
	)
}
