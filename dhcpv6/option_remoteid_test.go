package dhcpv6

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptRemoteID(t *testing.T) {
	expected := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	remoteId := []byte("DSLAM01 eth2/1/01/21")
	expected = append(expected, remoteId...)
	var opt OptRemoteID
	if err := opt.FromBytes(expected); err != nil {
		t.Fatal(err)
	}
	if en := opt.EnterpriseNumber; en != 0xaabbccdd {
		t.Fatalf("Invalid Enterprise Number. Expected 0xaabbccdd, got %v", en)
	}
	if rid := opt.RemoteID; !bytes.Equal(rid, remoteId) {
		t.Fatalf("Invalid Remote ID. Expected %v, got %v", expected, rid)
	}
}

func TestOptRemoteIDToBytes(t *testing.T) {
	remoteId := []byte("DSLAM01 eth2/1/01/21")
	expected := append([]byte{0, 0, 0, 0}, remoteId...)
	opt := OptRemoteID{
		RemoteID: remoteId,
	}
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, expected) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, toBytes)
	}
}

func TestOptRemoteIDParseOptRemoteIDTooShort(t *testing.T) {
	buf := []byte{0xaa, 0xbb, 0xcc}
	var opt OptRemoteID
	err := opt.FromBytes(buf)
	require.Error(t, err, "A short option should return an error")
}

func TestOptRemoteIDString(t *testing.T) {
	buf := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	remoteId := []byte("Test1234")
	buf = append(buf, remoteId...)

	var opt OptRemoteID
	err := opt.FromBytes(buf)
	require.NoError(t, err)
	str := opt.String()
	require.Contains(
		t,
		str,
		"EnterpriseNumber=2864434397",
		"String() should contain the enterprisenum",
	)
	require.Contains(
		t,
		str,
		"RemoteID=0x5465737431323334",
		"String() should contain the remoteid bytes",
	)
}
