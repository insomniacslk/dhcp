package dhcpv6

import (
	"bytes"
	"testing"
)

func TestDuidUuid(t *testing.T) {
	buf := []byte{
		0x00, 0x04,                                                                                     // type
		0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00, 0x08, 0x00, 0x09, // uuid
	}
	duid, err := DuidFromBytes(buf)
	if err != nil {
		t.Fatal(err)
	}
	if dt := duid.Type; dt != DUID_UUID {
		t.Fatalf("Invalid Preferred Lifetime. Expected 4, got %d", dt)
	}
	if uuid := duid.Uuid; !bytes.Equal(uuid[:], buf[2:]) {
		t.Fatalf("Invalid UUID. Expected %v, got %v", buf[2:], uuid)
	}
	if mac := duid.LinkLayerAddr; mac != nil {
		t.Fatalf("Invalid MAC. Expected nil, got %v", mac)
	}
}

func TestDuidUuidToBytes(t *testing.T) {
	uuid := [16]byte{0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00, 0x08, 0x00, 0x09}
	expected := []byte{00, 04}
	expected = append(expected, uuid[:]...)
	duid := Duid{
		Type: DUID_UUID,
		Uuid: uuid,
	}
	toBytes := duid.ToBytes()
	if !bytes.Equal(toBytes, expected) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, toBytes)
	}
}
