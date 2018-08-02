package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserClassModifier(t *testing.T) {
	d, _ := New()
	userClass := WithUserClass([]byte("linuxboot"))
	d = userClass(d)
	require.Equal(t, "User Class Information -> linuxboot", d.options[0].String())
}
