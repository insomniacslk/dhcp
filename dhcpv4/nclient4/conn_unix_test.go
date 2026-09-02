// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.12 && (darwin || freebsd || linux || netbsd || openbsd || dragonfly)

package nclient4

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/mdlayher/packet"
	"golang.org/x/net/bpf"
)

var errNoMorePackets = errors.New("no more packets")

// mockPacketConn replays a preset sequence of packets, then returns
// errNoMorePackets.
type mockPacketConn struct {
	reads [][]byte
	i     int
}

func (m *mockPacketConn) ReadFrom(p []byte) (int, net.Addr, error) {
	if m.i >= len(m.reads) {
		return 0, nil, errNoMorePackets
	}
	n := copy(p, m.reads[m.i])
	m.i++
	return n, &net.UDPAddr{}, nil
}
func (m *mockPacketConn) WriteTo(p []byte, _ net.Addr) (int, error) { return len(p), nil }
func (m *mockPacketConn) Close() error                              { return nil }
func (m *mockPacketConn) LocalAddr() net.Addr                       { return &net.UDPAddr{} }
func (m *mockPacketConn) SetDeadline(_ time.Time) error             { return nil }
func (m *mockPacketConn) SetReadDeadline(_ time.Time) error         { return nil }
func (m *mockPacketConn) SetWriteDeadline(_ time.Time) error        { return nil }

// frame builds an IPv4+UDP frame. ipTotalLen is written into the IPv4
// total-length field; dhcp is placed after the UDP header. The frame always
// carries a full 8-byte UDP header, so a small ipTotalLen models a header that
// declares fewer payload bytes than the frame physically contains.
func frame(ipTotalLen int, dhcp []byte) []byte {
	pkt := make([]byte, 28+len(dhcp))
	pkt[0] = 0x45 // IPv4, IHL 5 -> header length 20
	binary.BigEndian.PutUint16(pkt[2:], uint16(ipTotalLen))
	pkt[8] = 64 // TTL
	pkt[9] = byte(udpProtocolNumber)
	copy(pkt[12:16], []byte{1, 2, 3, 4}) // source address
	copy(pkt[16:20], []byte{255, 255, 255, 255})
	binary.BigEndian.PutUint16(pkt[20:], 67) // UDP source port
	binary.BigEndian.PutUint16(pkt[22:], 68) // UDP dest port (matches boundAddr)
	binary.BigEndian.PutUint16(pkt[24:], uint16(8+len(dhcp)))
	copy(pkt[28:], dhcp)
	return pkt
}

// TestReadFromShortIPv4Payload covers the IPv4 payload-length boundary around a
// UDP header. A frame that declares fewer than eight payload bytes must be
// skipped; before the fix such a frame produced a negative DHCP length that
// reached buf.Consume with a negative slice bound.
func TestReadFromShortIPv4Payload(t *testing.T) {
	for _, tc := range []struct {
		name       string
		ipTotalLen int
		skip       bool
	}{
		{"payload 0, dhcpLen -8", 20, true},
		{"payload 7, dhcpLen -1", 27, true},
		{"payload 8, dhcpLen 0", 28, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			upc := &BroadcastRawUDPConn{
				PacketConn: &mockPacketConn{reads: [][]byte{frame(tc.ipTotalLen, nil)}},
				boundAddr:  &net.UDPAddr{Port: 68},
			}
			n, _, err := upc.ReadFrom(make([]byte, 512))
			if tc.skip {
				if !errors.Is(err, errNoMorePackets) {
					t.Fatalf("want the frame skipped (next read %v), got n=%d err=%v", errNoMorePackets, n, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("want the frame accepted, got err=%v", err)
			}
			if n != 0 {
				t.Fatalf("want 0 dhcp bytes, got %d", n)
			}
		})
	}
}

// TestReadFromSkipsShortFrameThenReadsNext makes sure a skipped frame does not
// stop the receive loop: a following valid frame is still returned with its
// payload and source address intact. This guards against turning the skip into
// a hard error.
func TestReadFromSkipsShortFrameThenReadsNext(t *testing.T) {
	dhcp := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	upc := &BroadcastRawUDPConn{
		PacketConn: &mockPacketConn{reads: [][]byte{
			frame(27, nil),            // declares 7 payload bytes: skipped
			frame(28+len(dhcp), dhcp), // valid: 8-byte UDP header + 4 dhcp bytes
		}},
		boundAddr: &net.UDPAddr{Port: 68},
	}
	b := make([]byte, 512)
	n, srcAddr, err := upc.ReadFrom(b)
	if err != nil {
		t.Fatalf("ReadFrom returned err=%v, want the valid frame", err)
	}
	if n != len(dhcp) || !bytes.Equal(b[:n], dhcp) {
		t.Fatalf("ReadFrom payload = %x (n=%d), want %x", b[:n], n, dhcp)
	}
	if ua, ok := srcAddr.(*net.UDPAddr); !ok || ua.Port != 67 {
		t.Fatalf("ReadFrom srcAddr = %v, want UDP source port 67", srcAddr)
	}
}

// TestDHCPClientFilter runs the kernel prefilter program in the userspace BPF VM
// to check its logic: accept an IPv4/UDP frame destined to the client port, drop
// everything else. The VM sees the frame starting at the IPv4 header, which is
// what the SOCK_DGRAM socket presents to the kernel-attached filter.
func TestDHCPClientFilter(t *testing.T) {
	vm, err := bpf.NewVM(dhcpClientFilter(68))
	if err != nil {
		t.Fatal(err)
	}
	accepted := func(pkt []byte) bool {
		out, err := vm.Run(pkt)
		if err != nil {
			t.Fatalf("vm run: %v", err)
		}
		return out > 0
	}

	// ipOpt builds a frame whose IPv4 header carries options (IHL > 5), so the UDP
	// port is not at a fixed offset and the filter must track IHL to find it.
	ipOpt := func(ihl int, dport uint16) []byte {
		hdrLen := ihl * 4
		pkt := make([]byte, hdrLen+8)
		pkt[0] = 0x40 | byte(ihl)
		pkt[9] = byte(udpProtocolNumber)
		binary.BigEndian.PutUint16(pkt[hdrLen+2:], dport)
		return pkt
	}
	// tweak returns a valid port-68 frame with one big-endian 16-bit field set.
	tweak := func(off int, v uint16) []byte {
		pkt := frame(28+4, []byte{1, 2, 3, 4})
		binary.BigEndian.PutUint16(pkt[off:], v)
		return pkt
	}
	notUDP := frame(28+4, []byte{1, 2, 3, 4})
	notUDP[9] = 6 // TCP

	for _, tc := range []struct {
		name string
		pkt  []byte
		want bool
	}{
		{"dhcp reply to port 68", frame(28+4, []byte{1, 2, 3, 4}), true},
		{"ip options ihl 6 to port 68", ipOpt(6, 68), true},
		{"ip options ihl 15 to port 68", ipOpt(15, 68), true},
		{"ihl 4 is malformed", ipOpt(4, 68), false},
		{"udp to port 67", tweak(22, 67), false},
		{"ip options ihl 6 to port 67", ipOpt(6, 67), false},
		{"not udp", notUDP, false},
		{"truncated before the port", []byte{0x45, 0, 0, 20}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := accepted(tc.pkt); got != tc.want {
				t.Fatalf("filter accepted = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestListenFiltered covers the open path: the socket carries the filter from
// the start, and a refused filter still leaves the client with a working socket.
func TestListenFiltered(t *testing.T) {
	saved := listenPacket
	t.Cleanup(func() { listenPacket = saved })

	// record captures each config opened with; refuse fails the filtered attempt.
	var cfgs []*packet.Config
	record := func(refuse bool) {
		cfgs = nil
		listenPacket = func(_ *net.Interface, cfg *packet.Config) (net.PacketConn, error) {
			cfgs = append(cfgs, cfg)
			if refuse && cfg != nil {
				return nil, errors.New("attach refused")
			}
			return &mockPacketConn{}, nil
		}
	}

	t.Run("filter applied when the socket is opened", func(t *testing.T) {
		record(false)
		if _, err := listenFiltered(&net.Interface{}, 68); err != nil {
			t.Fatalf("listenFiltered: %v", err)
		}
		if len(cfgs) != 1 {
			t.Fatalf("opened the socket %d times, want 1", len(cfgs))
		}
		if cfgs[0] == nil || len(cfgs[0].Filter) == 0 {
			t.Fatal("the socket was opened without a filter, so traffic can be captured before it applies")
		}
	})

	t.Run("unfiltered fallback when the filter is refused", func(t *testing.T) {
		record(true)
		conn, err := listenFiltered(&net.Interface{}, 68)
		if err != nil {
			t.Fatalf("listenFiltered: %v, want the unfiltered fallback", err)
		}
		if conn == nil {
			t.Fatal("no socket returned")
		}
		if len(cfgs) != 2 {
			t.Fatalf("opened the socket %d times, want 2 (filtered, then plain)", len(cfgs))
		}
		if cfgs[1] != nil {
			t.Fatal("the retry still carried a filter")
		}
	})

	t.Run("both failures are reported", func(t *testing.T) {
		filterErr := errors.New("attach refused")
		plainErr := errors.New("no such device")
		listenPacket = func(_ *net.Interface, cfg *packet.Config) (net.PacketConn, error) {
			if cfg != nil {
				return nil, filterErr
			}
			return nil, plainErr
		}
		_, err := listenFiltered(&net.Interface{}, 68)
		if !errors.Is(err, filterErr) || !errors.Is(err, plainErr) {
			t.Fatalf("error = %v, want both %v and %v", err, filterErr, plainErr)
		}
	})
}

// TestNewRawUDPConnRejectsBadPort checks a port outside uint16 is rejected before
// any socket is opened, rather than silently truncated into the filter.
func TestNewRawUDPConnRejectsBadPort(t *testing.T) {
	// A nonexistent interface fails regardless, so match the specific
	// port-range error to prove the guard ran rather than InterfaceByName.
	for _, port := range []int{-1, 0x10000} {
		if _, err := NewRawUDPConn("nonexistent0", port); !errors.Is(err, errPortOutOfRange) {
			t.Errorf("NewRawUDPConn(port %d) error = %v, want errPortOutOfRange", port, err)
		}
	}
}
