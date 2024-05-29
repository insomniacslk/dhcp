package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptAutoConfigure(t *testing.T) {
	o := OptAutoConfigure(0)
	require.Equal(t, OptionAutoConfigure, o.Code, "Code")
	require.Equal(t, []byte{0}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Auto-Configure: DoNotAutoConfigure", o.String())

	o = OptAutoConfigure(1)
	require.Equal(t, OptionAutoConfigure, o.Code, "Code")
	require.Equal(t, []byte{1}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Auto-Configure: AutoConfigure", o.String())

	o = OptAutoConfigure(2)
	require.Equal(t, OptionAutoConfigure, o.Code, "Code")
	require.Equal(t, []byte{2}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Auto-Configure: UNKNOWN (2)", o.String())
}

func TestGetAutoConfigure(t *testing.T) {
	m, _ := New(WithGeneric(OptionAutoConfigure, []byte{1}))
	o, ok := m.AutoConfigure()
	require.True(t, ok)
	require.Equal(t, AutoConfigure, o)

	// Missing
	m, _ = New()
	_, ok = m.AutoConfigure()
	require.False(t, ok, "should get error if option missing")

	// Short byte stream
	m, _ = New(WithGeneric(OptionAutoConfigure, []byte{}))
	_, ok = m.AutoConfigure()
	require.False(t, ok, "should get error from short byte stream")

	// Bad length
	m, _ = New(WithGeneric(OptionAutoConfigure, []byte{2, 2}))
	_, ok = m.AutoConfigure()
	require.False(t, ok, "should get error from bad length")
}
