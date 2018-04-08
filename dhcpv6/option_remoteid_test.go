package dhcpv6

import (
	"bytes"
	"testing"
)

func TestOptRemoteId(t *testing.T) {
	expected := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	remoteId := []byte("DSLAM01 eth2/1/01/21")
	expected = append(expected, remoteId...)
	opt, err := ParseOptRemoteId(expected)
	if err != nil {
		t.Fatal(err)
	}
	if optLen := opt.Length(); optLen != len(expected) {
		t.Fatalf("Invalid length. Expected %v, got %v", len(expected), optLen)
	}
	if en := opt.EnterpriseNumber(); en != 0xaabbccdd {
		t.Fatalf("Invalid Enterprise Number. Expected 0xaabbccdd, got %v", en)
	}
	if rid := opt.RemoteID(); !bytes.Equal(rid, remoteId) {
		t.Fatalf("Invalid Remote ID. Expected %v, got %v", expected, rid)
	}
}

func TestOptRemoteIdToBytes(t *testing.T) {
	remoteId := []byte("DSLAM01 eth2/1/01/21")
	expected := []byte{00, 37, 00, byte(len(remoteId) + 4), 00, 00, 00, 00}
	expected = append(expected, remoteId...)
	opt := OptRemoteId{
		remoteId: remoteId,
	}
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, expected) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, toBytes)
	}
}
