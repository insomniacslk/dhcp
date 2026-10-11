package dhcpv4

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptTimeOffset(t *testing.T) {
	o := OptTimeOffset(-5 * time.Hour)
	require.Equal(t, OptionTimeOffset, o.Code, "Code")
	require.Equal(t, []byte{0xff, 0xff, 0xb9, 0xb0}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Time Offset: -5h0m0s", o.String(), "String")

	o = OptTimeOffset(5 * time.Hour)
	require.Equal(t, []byte{0x00, 0x00, 0x46, 0x50}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Time Offset: 5h0m0s", o.String(), "String")

	o = OptTimeOffset(0)
	require.Equal(t, []byte{0, 0, 0, 0}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Time Offset: 0s", o.String(), "String")
}

func TestGetTimeOffset(t *testing.T) {
	// UTC-5 is the two's complement value 0xffffb9b0, not an unsigned duration.
	m, err := New(WithGeneric(OptionTimeOffset, []byte{0xff, 0xff, 0xb9, 0xb0}))
	require.NoError(t, err)
	require.Equal(t, -5*time.Hour, m.TimeOffset(0))
	require.Contains(t, m.Summary(), "Time Offset: -5h0m0s")

	m, err = New(WithGeneric(OptionTimeOffset, []byte{0x00, 0x00, 0x46, 0x50}))
	require.NoError(t, err)
	require.Equal(t, 5*time.Hour, m.TimeOffset(0))
	require.Contains(t, m.Summary(), "Time Offset: 5h0m0s")

	m, err = New(WithGeneric(OptionTimeOffset, []byte{0, 0, 0, 0}))
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), m.TimeOffset(time.Hour))

	// Too short.
	m, err = New(WithGeneric(OptionTimeOffset, []byte{0xff, 0xff}))
	require.NoError(t, err)
	require.Equal(t, 10*time.Second, m.TimeOffset(10*time.Second))

	// Too long.
	m, err = New(WithGeneric(OptionTimeOffset, []byte{1, 1, 1, 1, 1}))
	require.NoError(t, err)
	require.Equal(t, 10*time.Second, m.TimeOffset(10*time.Second))

	// Missing.
	m, err = New()
	require.NoError(t, err)
	require.Equal(t, 10*time.Second, m.TimeOffset(10*time.Second))
}

func TestTimeOffsetExtremes(t *testing.T) {
	for _, secs := range []int32{math.MinInt32, math.MaxInt32} {
		want := time.Duration(secs) * time.Second
		o := OptTimeOffset(want)
		var got TimeOffset
		require.NoError(t, got.FromBytes(o.Value.ToBytes()))
		require.Equal(t, want, time.Duration(got))
	}
}
