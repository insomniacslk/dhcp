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

// withUDPLength overrides the UDP Length field written by frame, so a test can
// build a datagram whose UDP length disagrees with the IPv4 payload it rides in.
func withUDPLength(pkt []byte, udpLen int) []byte {
	binary.BigEndian.PutUint16(pkt[24:], uint16(udpLen))
	return pkt
}

// frameIHL6 builds a valid IPv4+UDP frame with a 24-byte header (IHL 6, four
// zero option bytes), so a test can check the UDP length is read past a
// variable-length IPv4 header rather than at a fixed offset.
func frameIHL6(dhcp []byte) []byte {
	pkt := make([]byte, 32+len(dhcp))
	pkt[0] = 0x46 // IPv4, IHL 6 -> 24-byte header
	binary.BigEndian.PutUint16(pkt[2:], uint16(32+len(dhcp)))
	pkt[8] = 64
	pkt[9] = byte(udpProtocolNumber)
	copy(pkt[12:16], []byte{1, 2, 3, 4})
	copy(pkt[16:20], []byte{255, 255, 255, 255})
	binary.BigEndian.PutUint16(pkt[24:], 67)                  // UDP source port
	binary.BigEndian.PutUint16(pkt[26:], 68)                  // UDP dest port
	binary.BigEndian.PutUint16(pkt[28:], uint16(8+len(dhcp))) // UDP length
	copy(pkt[32:], dhcp)
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

// TestReadFromValidatesUDPLength covers the UDP Length field, which the earlier
// IPv4 payload check does not look at. A length outside the header size and the
// IPv4 payload marks a malformed datagram and must be skipped; within those
// bounds, the IPv4 payload past it is a surplus area rather than dhcp data.
func TestReadFromValidatesUDPLength(t *testing.T) {
	dhcp := []byte{0xAA, 0xBB, 0xCC, 0xDD}

	// withSurplus carries four dhcp bytes and four surplus bytes in a 16-byte
	// IPv4 payload; a UDP length of 12 stops at the dhcp bytes.
	withSurplus := frame(20+8+len(dhcp)+4, []byte{0xAA, 0xBB, 0xCC, 0xDD, 0, 0, 0, 0})

	for _, tc := range []struct {
		name string
		pkt  []byte
		skip bool
		want []byte
	}{
		{"length 0, under the header", withUDPLength(frame(32, dhcp), 0), true, nil},
		{"length 7, under the header", withUDPLength(frame(32, dhcp), 7), true, nil},
		{"length 8, empty payload", frame(28, nil), false, nil},
		{"length past the ip payload", withUDPLength(frame(32, dhcp), 13), true, nil},
		{"length short, surplus dropped", withUDPLength(withSurplus, 8+len(dhcp)), false, dhcp},
		{"length equals the ip payload", frame(28+len(dhcp), dhcp), false, dhcp},
		{"ihl 6 header, valid length", frameIHL6(dhcp), false, dhcp},
	} {
		t.Run(tc.name, func(t *testing.T) {
			upc := &BroadcastRawUDPConn{
				PacketConn: &mockPacketConn{reads: [][]byte{tc.pkt}},
				boundAddr:  &net.UDPAddr{Port: 68},
			}
			b := make([]byte, 512)
			n, _, err := upc.ReadFrom(b)
			if tc.skip {
				if !errors.Is(err, errNoMorePackets) {
					t.Fatalf("want the frame skipped (next read %v), got n=%d err=%v", errNoMorePackets, n, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("want the frame accepted, got err=%v", err)
			}
			if !bytes.Equal(b[:n], tc.want) {
				t.Fatalf("dhcp payload = %x (n=%d), want %x", b[:n], n, tc.want)
			}
		})
	}
}

// TestReadFromSkipsBadUDPLengthThenReadsNext confirms a frame dropped for a bad
// UDP length does not stall the receive loop: the next valid frame is returned,
// not the dropped one. The two frames carry different payloads so a read that
// wrongly returns the dropped frame is caught.
func TestReadFromSkipsBadUDPLengthThenReadsNext(t *testing.T) {
	want := []byte{0x11, 0x22, 0x33, 0x44}
	upc := &BroadcastRawUDPConn{
		PacketConn: &mockPacketConn{reads: [][]byte{
			withUDPLength(frame(32, []byte{0xDE, 0xAD, 0xBE, 0xEF}), 99), // UDP length past the IPv4 payload: skipped
			frame(28+len(want), want),                                    // valid
		}},
		boundAddr: &net.UDPAddr{Port: 68},
	}
	b := make([]byte, 512)
	n, srcAddr, err := upc.ReadFrom(b)
	if err != nil {
		t.Fatalf("ReadFrom returned err=%v, want the valid frame", err)
	}
	if n != len(want) || !bytes.Equal(b[:n], want) {
		t.Fatalf("ReadFrom payload = %x (n=%d), want %x", b[:n], n, want)
	}
	if ua, ok := srcAddr.(*net.UDPAddr); !ok || ua.Port != 67 {
		t.Fatalf("ReadFrom srcAddr = %v, want UDP source port 67", srcAddr)
	}
}
