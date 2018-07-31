package dhcpv4

import (
	"bytes"
	"testing"
)

func TestLabelsFromBytes(t *testing.T) {
	labels, err := labelsFromBytes([]byte{
		0x9, 's', 'l', 'a', 'c', 'k', 'w', 'a', 'r', 'e',
		0x2, 'i', 't',
		0x0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 1 {
		t.Fatalf("Invalid labels length. Expected: 1, got: %v", len(labels))
	}
	if labels[0] != "slackware.it" {
		t.Fatalf("Invalid label. Expected: %v, got: %v'", "slackware.it", labels[0])
	}
}

func TestLabelsFromBytesZeroLength(t *testing.T) {
	labels, err := labelsFromBytes([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 0 {
		t.Fatalf("Invalid labels length. Expected: 0, got: %v", len(labels))
	}
}

func TestLabelsFromBytesInvalidLength(t *testing.T) {
	labels, err := labelsFromBytes([]byte{0x3, 0xaa, 0xbb}) // short length
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if len(labels) != 0 {
		t.Fatalf("Invalid labels length. Expected: 0, got: %v", len(labels))
	}
	if labels != nil {
		t.Fatalf("Invalid label. Expected nil, got %v", labels)
	}
}

func TestLabelToBytes(t *testing.T) {
	encodedLabel := labelToBytes("slackware.it")
	expected := []byte{
		0x9, 's', 'l', 'a', 'c', 'k', 'w', 'a', 'r', 'e',
		0x2, 'i', 't',
		0x0,
	}
	if !bytes.Equal(encodedLabel, expected) {
		t.Fatalf("Invalid label. Expected: %v, got: %v", expected, encodedLabel)
	}
}

func TestLabelToBytesZeroLength(t *testing.T) {
	encodedLabel := labelToBytes("")
	expected := []byte{0}
	if !bytes.Equal(encodedLabel, expected) {
		t.Fatalf("Invalid label. Expected: %v, got: %v", expected, encodedLabel)
	}
}
