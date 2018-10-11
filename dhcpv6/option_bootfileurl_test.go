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
	if optLen := opt.Length(); optLen != len(expected) {
		t.Fatalf("Invalid length. Expected %v, got %v", len(expected), optLen)
	}
	if url := opt.BootFileURL; !bytes.Equal(url, expected) {
		t.Fatalf("Invalid boot file URL. Expected %v, got %v", expected, url)
	}
	require.Contains(t, opt.String(), "BootFileUrl=https://insomniac.slackware.it", "String() should contain the correct BootFileUrl output")
}

func TestOptBootFileURLToBytes(t *testing.T) {
	urlString := []byte("https://insomniac.slackware.it")
	expected := []byte{00, 59, 00, byte(len(urlString))}
	expected = append(expected, urlString...)
	opt := OptBootFileURL{
		BootFileURL: urlString,
	}
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, expected) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, toBytes)
	}
}
