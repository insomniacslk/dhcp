package dhcpv6

import (
	"bytes"
	"testing"
	"time"
)

func TestOptInformationRefreshTime(t *testing.T) {
	var opt optInformationRefreshTime
	err := opt.FromBytes([]byte{0xaa, 0xbb, 0xcc, 0xdd})
	if err != nil {
		t.Fatal(err)
	}
	if informationRefreshTime := opt.InformationRefreshtime; informationRefreshTime != time.Duration(0xaabbccdd)*time.Second {
		t.Fatalf("Invalid information refresh time. Expected 0xaabb, got %v", informationRefreshTime)
	}
}

func TestOptInformationRefreshTimeToBytes(t *testing.T) {
	opt := OptInformationRefreshTime(0)
	expected := []byte{0, 0, 0, 0}
	if toBytes := opt.ToBytes(); !bytes.Equal(expected, toBytes) {
		t.Fatalf("Invalid ToBytes output. Expected %v, got %v", expected, toBytes)
	}
}

func TestOptInformationRefreshTimeString(t *testing.T) {
	opt := OptInformationRefreshTime(3600 * time.Second)
	expected := "Information Refresh Time: 1h0m0s"
	if optString := opt.String(); optString != expected {
		t.Fatalf("Invalid elapsed time string. Expected %v, got %v", expected, optString)
	}
}
