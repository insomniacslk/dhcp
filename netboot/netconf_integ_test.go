// +build integration

package netboot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIfUp(t *testing.T) {
	// this assumes that eth0 exists and is configurable
	ifname := "eth0"
	iface, err := IfUp(ifname, 2*time.Second)
	require.NoError(t, err)
	assert.Equal(t, ifname, iface.Name)
}

func TestIfUpTimeout(t *testing.T) {
	// this assumes that eth0 exists and is configurable
	ifname := "eth0"
	_, err := IfUp(ifname, 0*time.Second)
	require.Error(t, err)
}
