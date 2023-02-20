package dhcpv6

import (
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	var m Message
	m.Options.Add(OptElapsedTime(2 * time.Second))
}
