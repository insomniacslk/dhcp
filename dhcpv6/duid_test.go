package dhcpv6

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestDuidInvalidTooShort(t *testing.T) {
	// too short DUID at all (must be at least 2 bytes)
	_, err := DuidFromBytes([]byte{0})
	require.Error(t, err)

	// too short DUID_LL (must be at least 4 bytes)
	_, err = DuidFromBytes([]byte{0, 3, 0xa})
	require.Error(t, err)

	// too short DUID_EN (must be at least 6 bytes)
	_, err = DuidFromBytes([]byte{0, 2, 0xa, 0xb, 0xc})
	require.Error(t, err)

	// too short DUID_LLT (must be at least 8 bytes)
	_, err = DuidFromBytes([]byte{0, 1, 0xa, 0xb, 0xc, 0xd, 0xe})
	require.Error(t, err)

	// too short DUID_UUID (must be at least 18 bytes)
	_, err = DuidFromBytes([]byte{0, 4, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf})
	require.Error(t, err)
}

func TestDuidLLTFromBytes(t *testing.T) {
	buf := []byte{
		0, 1, // DUID_LLT
		0, 1, // HwTypeEthernet
		0x01, 0x02, 0x03, 0x04, // time
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, // link-layer addr
	}
	duid, err := DuidFromBytes(buf)
	require.NoError(t, err)
	require.Equal(t, 14, duid.Length())
	require.Equal(t, DUID_LLT, duid.Type)
	require.Equal(t, uint32(0x01020304), duid.Time)
	require.Equal(t, iana.HWTypeEthernet, duid.HwType)
	require.Equal(t, net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, duid.LinkLayerAddr)
}

func TestDuidLLFromBytes(t *testing.T) {
	buf := []byte{
		0, 3, // DUID_LL
		0, 1, // HwTypeEthernet
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, // link-layer addr
	}
	duid, err := DuidFromBytes(buf)
	require.NoError(t, err)
	require.Equal(t, 10, duid.Length())
	require.Equal(t, DUID_LL, duid.Type)
	require.Equal(t, iana.HWTypeEthernet, duid.HwType)
	require.Equal(t, net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, duid.LinkLayerAddr)
}

func TestDuidUuidFromBytes(t *testing.T) {
	buf := []byte{
		0x00, 0x04, // DUID_UUID
	}
	uuid := []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00, 0x08}
	buf = append(buf, uuid...)
	duid, err := DuidFromBytes(buf)
	require.NoError(t, err)
	require.Equal(t, 18, duid.Length())
	require.Equal(t, DUID_UUID, duid.Type)
	require.Equal(t, uuid, duid.Uuid)
}

func TestDuidLLTToBytes(t *testing.T) {
	expected := []byte{
		0, 1, // DUID_LLT
		0, 1, // HwTypeEthernet
		0x01, 0x02, 0x03, 0x04, // time
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, // link-layer addr
	}
	duid := Duid{
		Type:          DUID_LLT,
		HwType:        iana.HWTypeEthernet,
		Time:          uint32(0x01020304),
		LinkLayerAddr: []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
	}
	toBytes := duid.ToBytes()
	require.Equal(t, expected, toBytes)
}

func TestDuidUuidToBytes(t *testing.T) {
	uuid := []byte{0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00, 0x08, 0x00, 0x09}
	expected := []byte{00, 04}
	expected = append(expected, uuid...)
	duid := Duid{
		Type: DUID_UUID,
		Uuid: uuid,
	}
	toBytes := duid.ToBytes()
	if !bytes.Equal(toBytes, expected) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, toBytes)
	}
}

func TestOpaqueDuid(t *testing.T) {
	duid := []byte("\x00\x0a\x00\x03\x00\x01\x4c\x5e\x0c\x43\xbf\x39")
	d, err := DuidFromBytes(duid)
	if err != nil {
		t.Fatalf("DuidFromBytes: unexpected error: %v", err)
	}
	if got, want := d.Length(), len(duid); got != want {
		t.Errorf("Length: unexpected result: got %d, want %d", got, want)
	}
	if got, want := d.ToBytes(), duid; !bytes.Equal(got, want) {
		t.Fatalf("ToBytes: unexpected result: got %x, want %x", got, want)
	}
}

func TestDuidEqual(t *testing.T) {
	d := Duid{
		Type:          DUID_LL,
		HwType:        iana.HWTypeEthernet,
		LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
	}
	o := Duid{
		Type:          DUID_LL,
		HwType:        iana.HWTypeEthernet,
		LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
	}
	require.True(t, d.Equal(o))
}

func TestDuidEqualNotEqual(t *testing.T) {
	d := Duid{
		Type:          DUID_LL,
		HwType:        iana.HWTypeEthernet,
		LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
	}
	o := Duid{
		Type:          DUID_LL,
		HwType:        iana.HWTypeEthernet,
		LinkLayerAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0x00},
	}
	require.False(t, d.Equal(o))
}

func TestDuidFromBytesCreatesCopy(t *testing.T) {
	// Test that DuidFromBytes creates a copy of the underlying buffer

	items := []struct {
		duidType  string
		duidBytes []byte
	}{
		{"DUID-LL", []byte{0, 3, 0, 1, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
		{"DUID-LL", []byte{0, 3, 0, 1, 0xbb, 0xcc, 0xaa, 0xee, 0xff, 0xdd}},
		{"DUID-LLT", []byte{0, 1, 0, 1, 0x01, 0x02, 0x03, 0x04, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
		{"DUID-LLT", []byte{0, 1, 0, 1, 0x05, 0x06, 0x07, 0x08, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}},
		{"DUID-EN", []byte{0, 2, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}},
		{"DUID-EN", []byte{0, 2, 0x0e, 0x0f, 0x01, 0x02, 0x3, 0x4, 0x5}},
		{"DUID-UUID", []byte{0x00, 0x04, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00, 0x08, 0x00, 0x09}},
		{"DUID-UUID", []byte{0x00, 0x04, 0x02, 0x08, 0x07, 0x00, 0x04, 0x03, 0x00, 0x05, 0x00, 0x05, 0x03, 0x04, 0x01, 0x08, 0x07, 0x09}},
		{"Unknown", []byte("\x00\x0a\x00\x03\x00\x01\x4c\x5e\x0c\x43\xbf\x39")},
		{"Unknown", []byte("\x00\x0a\x00\x04\x02\x03\x4e\x5f\x03\x33\xfb\x93")},
		{"Unknown", []byte("\x00\x0a\x00\x03\x00\x01\x4c\x5e\x0c\x43\xbf\x39")},
		{"Unknown", []byte("\x00\x0a\x00\x04\x02\x03\x4e\x5f\x03\x33\xfb\x93")},
	}

	// Compute max DUID length from the above items
	maxLen := 0
	for _, item := range items {
		l := len(item.duidBytes)
		if l > maxLen {
			maxLen = l
		}
	}

	// Create shared buffer
	buf := make([]byte, maxLen)

	results := make([]*Duid, len(items))

	// For each item, copy it into buf, then parse it from buf and store the result
	for i, item := range items {
		copy(buf, item.duidBytes)
		duid, err := DuidFromBytes(buf[:len(item.duidBytes)])
		require.NoError(t, err)
		results[i] = duid
	}

	// Require that all the parsed DUID values still match the original byte values
	for i := range items {
		t.Run(fmt.Sprintf("item %d", i+1), func(t *testing.T) {
			require.Equal(t, items[i].duidBytes, results[i].ToBytes())
			require.Equal(t, items[i].duidType, results[i].Type.String())
		})
	}
}
