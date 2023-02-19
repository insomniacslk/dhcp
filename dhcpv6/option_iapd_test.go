package dhcpv6

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseMessageWithIAPD(t *testing.T) {
	data := []byte{
		0, 25, // IAPD option code
		0, 41, // length
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
		0, 0, 0, 2, // IAPrefix preferredLifetime
		0, 0, 0, 4, // IAPrefix validLifetime
		36,                                             // IAPrefix prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix
	}
	var got MessageOptions
	if err := got.FromBytes(data); err != nil {
		t.Errorf("FromBytes = %v", err)
	}

	want := &OptIAPD{
		IaId: [4]byte{1, 0, 0, 0},
		T1:   1 * time.Second,
		T2:   2 * time.Second,
		Options: PDOptions{Options: Options{&OptIAPrefix{
			PreferredLifetime: 2 * time.Second,
			ValidLifetime:     4 * time.Second,
			Prefix: &net.IPNet{
				Mask: net.CIDRMask(36, 128),
				IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			},
			Options: PrefixOptions{Options: Options{}},
		}}},
	}
	if gotIAPD := got.OneIAPD(); !reflect.DeepEqual(gotIAPD, want) {
		t.Errorf("OneIAPD = %v, want %v", gotIAPD, want)
	}
}

func TestOptIAPDParseOptIAPD(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
		0xaa, 0xbb, 0xcc, 0xdd, // IAPrefix preferredLifetime
		0xee, 0xff, 0x00, 0x11, // IAPrefix validLifetime
		36,                                             // IAPrefix prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix
	}
	var opt OptIAPD
	err := opt.FromBytes(data)
	require.NoError(t, err)
	require.Equal(t, OptionIAPD, opt.Code())
	require.Equal(t, [4]byte{1, 0, 0, 0}, opt.IaId)
	require.Equal(t, time.Second, opt.T1)
	require.Equal(t, 2*time.Second, opt.T2)
}

func TestOptIAPDParseOptIAPDInvalidLength(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		// truncated from here
	}
	var opt OptIAPD
	err := opt.FromBytes(data)
	require.Error(t, err)
}

func TestOptIAPDParseOptIAPDInvalidOptions(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
		0xaa, 0xbb, 0xcc, 0xdd, // IAPrefix preferredLifetime
		0xee, 0xff, 0x00, 0x11, // IAPrefix validLifetime
		36,                                          // IAPrefix prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // IAPrefix ipv6Prefix missing last byte
	}
	var opt OptIAPD
	err := opt.FromBytes(data)
	require.Error(t, err)
}

func TestOptIAPDToBytes(t *testing.T) {
	oaddr := OptIAPrefix{
		PreferredLifetime: 0xaabbccdd * time.Second,
		ValidLifetime:     0xeeff0011 * time.Second,
		Prefix: &net.IPNet{
			Mask: net.CIDRMask(36, 128),
			IP:   net.IPv6loopback,
		},
	}
	opt := OptIAPD{
		IaId:    [4]byte{1, 2, 3, 4},
		T1:      12345 * time.Second,
		T2:      54321 * time.Second,
		Options: PDOptions{[]Option{&oaddr}},
	}

	expected := []byte{
		1, 2, 3, 4, // IA ID
		0, 0, 0x30, 0x39, // T1 = 12345
		0, 0, 0xd4, 0x31, // T2 = 54321
		0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
		0xaa, 0xbb, 0xcc, 0xdd, // IAPrefix preferredLifetime
		0xee, 0xff, 0x00, 0x11, // IAPrefix validLifetime
		36,                                             // IAPrefix prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix
	}
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptIAPDString(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
		0xaa, 0xbb, 0xcc, 0xdd, // IAPrefix preferredLifetime
		0xee, 0xff, 0x00, 0x11, // IAPrefix validLifetime
		36,                                             // IAPrefix prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix
	}
	var opt OptIAPD
	err := opt.FromBytes(data)
	require.NoError(t, err)

	str := opt.String()
	require.Contains(
		t, str,
		"IAID=0x01000000",
		"String() should return the IAID",
	)
	require.Contains(
		t, str,
		"T1=1s T2=2s",
		"String() should return the T1/T2 options",
	)
	require.Contains(
		t, str,
		"Options={",
		"String() should return a list of options",
	)
}
