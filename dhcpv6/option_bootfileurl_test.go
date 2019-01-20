package dhcpv6

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptBootFileURL(t *testing.T) {
	expected := []byte("https://insomniac.slackware.it")
	opt, err := ParseOptBootFileURL(expected)
	if err != nil {
		t.Fatal(err)
	}
	if url := opt.BootFileURL; !bytes.Equal(url, expected) {
		t.Fatalf("Invalid boot file URL. Expected %v, got %v", expected, url)
	}
	require.Contains(t, opt.String(), "BootFileUrl=https://insomniac.slackware.it", "String() should contain the correct BootFileUrl output")
}

func TestOptBootFileURLToBytes(t *testing.T) {
	urlString := []byte("https://insomniac.slackware.it")
	opt := OptBootFileURL{
		BootFileURL: urlString,
	}
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, urlString) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", urlString, toBytes)
	}
}
