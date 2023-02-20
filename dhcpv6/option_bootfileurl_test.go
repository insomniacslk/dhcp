package dhcpv6

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMessageWithBootFileURL(t *testing.T) {
	buf := []byte{
		0, 59, // boot file option
		0, 3, // length
		0x66, 0x6f, 0x6f, //
	}

	want := "foo"
	var mo MessageOptions
	if err := mo.FromBytes(buf); err != nil {
		t.Errorf("FromBytes = %v", err)
	} else if got := mo.BootFileURL(); !reflect.DeepEqual(got, want) {
		t.Errorf("BootFileURL = %v, want %v", got, want)
	}
}

func TestOptBootFileURL(t *testing.T) {
	expected := "https://insomniac.slackware.it"
	var opt optBootFileURL
	if err := opt.FromBytes([]byte(expected)); err != nil {
		t.Fatal(err)
	}
	if opt.url != expected {
		t.Fatalf("Invalid boot file URL. Expected %v, got %v", expected, opt)
	}
	require.Contains(t, opt.String(), "https://insomniac.slackware.it", "String() should contain the correct BootFileUrl output")
}

func TestOptBootFileURLToBytes(t *testing.T) {
	urlString := "https://insomniac.slackware.it"
	opt := OptBootFileURL(urlString)
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, []byte(urlString)) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", urlString, toBytes)
	}
}
