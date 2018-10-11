package dhcpv6

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func captureStdout(t *testing.T, f func()) string {
	stdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	outC := make(chan string)
	go func() {
		var buf strings.Builder
		_, err := io.Copy(&buf, r)
		r.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "captureStdout(): copying pipe: %v\n", err)
			os.Exit(1)
		}
		outC <- buf.String()
	}()

	defer func() {
		w.Close()
		os.Stdout = stdout
	}()

	f()
	w.Close()
	out := <-outC

	return out
}

func TestOptRequestedOption(t *testing.T) {
	expected := []byte{0, 1, 0, 2}
	_, err := ParseOptRequestedOption(expected)
	require.NoError(t, err, "ParseOptRequestedOption() correct options should not error")
}

func TestOptRequestedOptionAddRequestedOptionDuplicate(t *testing.T) {
	opt := OptRequestedOption{}

	opt.AddRequestedOption(OptionDNSRecursiveNameServer)

	output := captureStdout(t, func() {
		opt.AddRequestedOption(OptionDNSRecursiveNameServer)
	})

	require.Contains(
		t,
		output,
		"appending duplicate",
		"AddRequestedOption() should complain to stdout when a duplicate entry is added",
	)
}

func TestOptRequestedOptionParseOptRequestedOptionTooShort(t *testing.T) {
	buf := []byte{0, 1, 0}
	_, err := ParseOptRequestedOption(buf)
	require.Error(t, err, "A short option should return an error (must be divisible by 2)")
}

func TestOptRequestedOptionString(t *testing.T) {
	buf := []byte{0, 1, 0, 2}
	opt, err := ParseOptRequestedOption(buf)
	require.NoError(t, err)
	require.Contains(
		t,
		opt.String(),
		"OPTION_CLIENTID, OPTION_SERVERID",
		"String() should contain the options specified",
	)
	opt.AddRequestedOption(12345)
	require.Contains(
		t,
		opt.String(),
		"Unknown",
		"String() should contain 'Unknown' for an illegal option",
	)
}
