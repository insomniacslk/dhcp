package dhcpv6

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func Test4RDParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []*Opt4RD
	}{
		{
			buf: []byte{
				0, 97, // 4RD option code
				0, 28, // length
				0, 98, // 4RD Map Rule option
				0, 24, // length
				16,             // prefix4-length
				16,             // prefix6-length
				8,              // ea-len
				0,              // WKPAuthorized
				192, 168, 0, 1, // rule-ipv4-prefix
				0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // rule-ipv6-prefix
			},
			want: []*Opt4RD{
				&Opt4RD{
					FourRDOptions: FourRDOptions{Options: Options{
						&Opt4RDMapRule{
							Prefix4: net.IPNet{
								IP:   net.IP{192, 168, 0, 1},
								Mask: net.CIDRMask(16, 32),
							},
							Prefix6: net.IPNet{
								IP:   net.ParseIP("fe80::"),
								Mask: net.CIDRMask(16, 128),
							},
							EABitsLength: 8,
						},
					}},
				},
			},
		},
		{
			buf: []byte{
				0, 97, // 4RD option code
				0, 28, // length
				0, 98, // 4RD Map Rule option
				0, 24, // length
				16,             // prefix4-length
				16,             // prefix6-length
				8,              // ea-len
				0,              // WKPAuthorized
				192, 168, 0, 1, // rule-ipv4-prefix
				0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // rule-ipv6-prefix

				0, 97, // 4RD
				0, 8, // length
				0, 99, // 4RD non map rule
				0, 4, // length
				0x80, 0x00, 0x05, 0xd4,
			},
			want: []*Opt4RD{
				&Opt4RD{
					FourRDOptions: FourRDOptions{Options: Options{
						&Opt4RDMapRule{
							Prefix4: net.IPNet{
								IP:   net.IP{192, 168, 0, 1},
								Mask: net.CIDRMask(16, 32),
							},
							Prefix6: net.IPNet{
								IP:   net.ParseIP("fe80::"),
								Mask: net.CIDRMask(16, 128),
							},
							EABitsLength: 8,
						},
					}},
				},
				&Opt4RD{
					FourRDOptions: FourRDOptions{Options: Options{
						&Opt4RDNonMapRule{
							HubAndSpoke: true,
							DomainPMTU:  1492,
						},
					}},
				},
			},
		},
		{
			buf:  []byte{0, 97, 0, 1, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
		{
			// Allowed, because the RFC doesn't really specify that
			// it can't be empty. RFC doesn't really specify
			// anything, frustratingly.
			buf: []byte{
				0, 97, // 4RD option code
				0, 0, // length
			},
			want: []*Opt4RD{&Opt4RD{FourRDOptions: FourRDOptions{Options: Options{}}}},
		},
		{
			buf: []byte{
				0, 97, // 4RD option code
				0, 6, // length
				0, 98, // 4RD Map Rule option
				0, 4, // length
				16, // prefix4-length
				16, // prefix6-length
				8,  // ea-len
				0,  // WKPAuthorized
				// Missing
			},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.FourRD(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FourRD = %v, want %v", got, tt.want)
			}
			if len(tt.want) >= 1 {
				var b MessageOptions
				for _, frd := range tt.want {
					b.Add(frd)
				}
				got := b.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func Test4RDMapRuleParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []*Opt4RDMapRule
	}{
		{
			buf: []byte{
				0, 98, // 4RD Map Rule option
				0, 24, // length
				16,             // prefix4-length
				16,             // prefix6-length
				8,              // ea-len
				0,              // WKPAuthorized
				192, 168, 0, 1, // rule-ipv4-prefix
				0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // rule-ipv6-prefix
			},
			want: []*Opt4RDMapRule{
				&Opt4RDMapRule{
					Prefix4: net.IPNet{
						IP:   net.IP{192, 168, 0, 1},
						Mask: net.CIDRMask(16, 32),
					},
					Prefix6: net.IPNet{
						IP:   net.ParseIP("fe80::"),
						Mask: net.CIDRMask(16, 128),
					},
					EABitsLength: 8,
				},
			},
		},
		{
			buf: []byte{
				0, 98, // 4RD Map Rule option
				0, 24, // length
				16,             // prefix4-length
				16,             // prefix6-length
				8,              // ea-len
				1 << 7,         // WKPAuthorized
				192, 168, 0, 1, // rule-ipv4-prefix
				0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // rule-ipv6-prefix
			},
			want: []*Opt4RDMapRule{
				&Opt4RDMapRule{
					Prefix4: net.IPNet{
						IP:   net.IP{192, 168, 0, 1},
						Mask: net.CIDRMask(16, 32),
					},
					Prefix6: net.IPNet{
						IP:   net.ParseIP("fe80::"),
						Mask: net.CIDRMask(16, 128),
					},
					EABitsLength:  8,
					WKPAuthorized: true,
				},
			},
		},
		{
			buf:  []byte{0, 98, 0, 1, 0},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 98, // 4RD Map Rule option
				0, 4, // length
				16, // prefix4-length
				16, // prefix6-length
				8,  // ea-len
				0,  // WKPAuthorized
				// Missing
			},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var frdo FourRDOptions
			if err := frdo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := frdo.MapRules(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapRules = %v, want %v", got, tt.want)
			}
			if len(tt.want) >= 1 {
				var b FourRDOptions
				for _, frd := range tt.want {
					b.Add(frd)
				}
				got := b.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func Test4RDNonMapRuleParseAndGetter(t *testing.T) {
	trafficClassOne := uint8(1)
	for i, tt := range []struct {
		buf  []byte
		err  error
		want *Opt4RDNonMapRule
	}{
		{
			buf: []byte{
				0, 99, // 4RD Non Map Rule option
				0, 4, // length
				0x80, 0, 0x05, 0xd4,
			},
			want: &Opt4RDNonMapRule{
				HubAndSpoke: true,
				DomainPMTU:  1492,
			},
		},
		{
			buf: []byte{
				0, 99, // 4RD Non Map Rule option
				0, 4, // length
				0, 0, 0x05, 0xd4,
			},
			want: &Opt4RDNonMapRule{
				DomainPMTU: 1492,
			},
		},
		{
			buf: []byte{
				0, 99, // 4RD Non Map Rule option
				0, 4, // length
				0x1, 0x01, 0x05, 0xd4,
			},
			want: &Opt4RDNonMapRule{
				TrafficClass: &trafficClassOne,
				DomainPMTU:   1492,
			},
		},
		{
			buf:  []byte{0, 99, 0, 1, 0},
			want: nil,
			err:  uio.ErrBufferTooShort,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var frdo FourRDOptions
			if err := frdo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := frdo.NonMapRule(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NonMapRule = %v, want %v", got, tt.want)
			}
			if tt.want != nil {
				var b FourRDOptions
				b.Add(tt.want)
				got := b.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOpt4RDNonMapRuleString(t *testing.T) {
	var tClass uint8 = 120
	opt := Opt4RDNonMapRule{
		HubAndSpoke:  true,
		TrafficClass: &tClass,
		DomainPMTU:   9000,
	}

	str := opt.String()

	require.Contains(t, str, "HubAndSpoke=true",
		"String() should contain the HubAndSpoke flag value")
	require.Contains(t, str, "TrafficClass=120",
		"String() should contain the TrafficClass flag value")
	require.Contains(t, str, "DomainPMTU=9000",
		"String() should contain the domain PMTU")
}

func TestOpt4RDMapRuleToBytes(t *testing.T) {
	opt := Opt4RDMapRule{
		EABitsLength:  32,
		WKPAuthorized: true,
	}

	expected := append([]byte{
		0,    // v4 prefix length
		0,    // v6 prefix length
		32,   // EA-bits
		0x80, // WKPs authorized
	}, bytes.Repeat([]byte{0x00}, 4+16)...)
	require.Equal(t, expected, opt.ToBytes())
}

func TestOpt4RDMapRuleString(t *testing.T) {
	opt := Opt4RDMapRule{
		Prefix4: net.IPNet{
			IP:   net.IPv4(100, 64, 0, 238),
			Mask: net.CIDRMask(24, 32),
		},
		Prefix6: net.IPNet{
			IP:   net.ParseIP("2001:db8::1234:5678:0:aabb"),
			Mask: net.CIDRMask(80, 128),
		},
		EABitsLength:  32,
		WKPAuthorized: true,
	}

	str := opt.String()
	require.Contains(t, str, "WKPAuthorized=true", "String() should write the flag values")
	require.Contains(t, str, "Prefix6=2001:db8::1234:5678:0:aabb/80",
		"String() should include the IPv6 prefix")
	require.Contains(t, str, "Prefix4=100.64.0.238/24",
		"String() should include the IPv4 prefix")
	require.Contains(t, str, "EA-Bits=32", "String() should include the value for EA-Bits")
}
