package dhcpv6

import (
	"bytes"
	"testing"
)

func TestOptElapsedTime(t *testing.T) {
	opt, err := ParseOptElapsedTime([]byte{0xaa, 0xbb})
	if err != nil {
		t.Fatal(err)
	}
	if optLen := opt.Length(); optLen != 2 {
		t.Fatalf("Invalid length. Expected 2, got %v", optLen)
	}
	if elapsedTime := opt.ElapsedTime; elapsedTime != 0xaabb {
		t.Fatalf("Invalid elapsed time. Expected 0xaabb, got %v", elapsedTime)
	}
}

func TestOptElapsedTimeToBytes(t *testing.T) {
	opt := OptElapsedTime{}
	expected := []byte{0, 8, 0, 2, 0, 0}
	if toBytes := opt.ToBytes(); !bytes.Equal(expected, toBytes) {
		t.Fatalf("Invalid ToBytes output. Expected %v, got %v", expected, toBytes)
	}
}

func TestOptElapsedTimeSetGetElapsedTime(t *testing.T) {
	opt := OptElapsedTime{}
	opt.ElapsedTime = 10
	if elapsedTime := opt.ElapsedTime; elapsedTime != 10 {
		t.Fatalf("Invalid elapsed time. Expected 10, got %v", elapsedTime)
	}
}

func TestOptElapsedTimeString(t *testing.T) {
	opt := OptElapsedTime{}
	opt.ElapsedTime = 10
	expected := "OptElapsedTime{elapsedtime=10}"
	if optString := opt.String(); optString != expected {
		t.Fatalf("Invalid elapsed time string. Expected %v, got %v", expected, optString)
	}
}
