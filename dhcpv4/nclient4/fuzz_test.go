// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.18 && (darwin || freebsd || linux || netbsd || openbsd || dragonfly)

package nclient4

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

// Model constants mirror the production IPv4/UDP layout but stay independent
// of the package constants, so a wrong change to one of those is caught by the
// fuzzer rather than copied into the oracle.
const (
	modelIPv4MinHeader          = 20
	modelIPv4MaxHeader          = 60
	modelUDPHeader              = 8
	modelUDPProtocol            = 17
	modelSrcIPOffset            = 12
	modelIPv4AddrSize           = 4
	modelFlagsFragOffset        = 6
	modelIPv4MoreFragments      = 0x2000
	modelIPv4FragmentOffsetMask = 0x1fff
)

// fuzzFrame builds an IPv4+UDP frame with the given total-length and UDP length
// fields. Both are parameters because the receive path bounds the payload at
// each, so a seed stopping short at one but not the other cannot be built if
// either is derived. Kept local so the target builds standalone under OSS-Fuzz.
func fuzzFrame(ipTotalLen, udpLen int, dhcp []byte) []byte {
	pkt := make([]byte, 28+len(dhcp))
	pkt[0] = 0x45 // IPv4, IHL 5
	binary.BigEndian.PutUint16(pkt[2:], uint16(ipTotalLen))
	pkt[9] = byte(modelUDPProtocol)
	binary.BigEndian.PutUint16(pkt[22:], 68) // UDP destination port
	binary.BigEndian.PutUint16(pkt[24:], uint16(udpLen))
	copy(pkt[28:], dhcp)
	return pkt
}

// fuzzFragment builds an otherwise valid frame with the given flags and fragment
// offset; the fuzzer is unlikely to reach one by chance.
func fuzzFragment(flagsFragOffset uint16, dhcp []byte) []byte {
	pkt := fuzzFrame(28+len(dhcp), modelUDPHeader+len(dhcp), dhcp)
	binary.BigEndian.PutUint16(pkt[modelFlagsFragOffset:], flagsFragOffset)
	return pkt
}

// fuzzConn hands ReadFrom the fuzz frame once, then reports EOF so the receive
// loop stops.
type fuzzConn struct {
	frame []byte
	done  bool
}

func (c *fuzzConn) ReadFrom(p []byte) (int, net.Addr, error) {
	if c.done {
		return 0, nil, io.EOF
	}
	c.done = true
	return copy(p, c.frame), &net.UDPAddr{}, nil
}
func (c *fuzzConn) WriteTo([]byte, net.Addr) (int, error) { return 0, nil }
func (c *fuzzConn) Close() error                          { return nil }
func (c *fuzzConn) LocalAddr() net.Addr                   { return &net.UDPAddr{} }
func (c *fuzzConn) SetDeadline(time.Time) error           { return nil }
func (c *fuzzConn) SetReadDeadline(time.Time) error       { return nil }
func (c *fuzzConn) SetWriteDeadline(time.Time) error      { return nil }

// expectedReadFrom recomputes what ReadFrom must return for the bytes delivered,
// so the target pins the exact length, payload, and source, not just the frame
// type. It encodes the whole contract: fragments are rejected, the IPv4 payload
// must hold a UDP header, UDP Length must fall inside it, and the dhcp bytes end
// at UDP Length, dropping the surplus area and anything past IP Total Length.
func expectedReadFrom(delivered []byte, bufLen int) (accept bool, wantN int, payload []byte, srcIP net.IP, srcPort int) {
	if len(delivered) < modelIPv4MinHeader {
		return false, 0, nil, nil, 0
	}
	hlen := int(delivered[0]&0x0f) * 4
	tlen := int(binary.BigEndian.Uint16(delivered[2:4]))
	// isValid: header within [20, total length], total length within the frame.
	if hlen < modelIPv4MinHeader || hlen > tlen || tlen > len(delivered) {
		return false, 0, nil, nil, 0
	}
	if delivered[0]>>4 != 4 { // IPv4 version
		return false, 0, nil, nil, 0
	}
	fragField := binary.BigEndian.Uint16(delivered[modelFlagsFragOffset : modelFlagsFragOffset+2])
	if fragField&(modelIPv4MoreFragments|modelIPv4FragmentOffsetMask) != 0 {
		return false, 0, nil, nil, 0
	}
	if delivered[9] != byte(modelUDPProtocol) {
		return false, 0, nil, nil, 0
	}
	ipPayloadLen := tlen - hlen
	if ipPayloadLen < modelUDPHeader { // room for a UDP header
		return false, 0, nil, nil, 0
	}
	// tlen <= len(delivered) and ipPayloadLen >= 8 together put the UDP header
	// within the delivered bytes.
	udpLen := int(binary.BigEndian.Uint16(delivered[hlen+4 : hlen+6]))
	if udpLen < modelUDPHeader || udpLen > ipPayloadLen {
		return false, 0, nil, nil, 0
	}
	if binary.BigEndian.Uint16(delivered[hlen+2:hlen+4]) != 68 { // destination port
		return false, 0, nil, nil, 0
	}
	// hlen+udpLen <= hlen+ipPayloadLen == tlen <= len(delivered), so the dhcp
	// bytes are all present.
	payload = delivered[hlen+modelUDPHeader : hlen+udpLen]
	wantN = udpLen - modelUDPHeader
	if wantN > bufLen {
		wantN = bufLen
	}
	srcIP = net.IP(delivered[modelSrcIPOffset : modelSrcIPOffset+modelIPv4AddrSize])
	srcPort = int(binary.BigEndian.Uint16(delivered[hlen : hlen+2]))
	return true, wantN, payload, srcIP, srcPort
}

// FuzzBroadcastRawUDPConnReadFrom feeds arbitrary bytes as one raw frame into
// ReadFrom with a fuzzed buffer size and checks the result against an
// independent model: a rejected frame must surface as an error, and an accepted
// frame must return exactly the payload, length, and source the model derives.
func FuzzBroadcastRawUDPConnReadFrom(f *testing.F) {
	dhcp := []byte("hello world dhcp payload")
	trailing := []byte("trailing data past total")
	seeds := [][]byte{
		fuzzFrame(20, modelUDPHeader, nil),                      // IPv4 payload 0, below the UDP header (#583)
		fuzzFrame(27, modelUDPHeader, nil),                      // IPv4 payload 7, below the UDP header (#583)
		fuzzFrame(28, modelUDPHeader, []byte{1, 2, 3}),          // IPv4 payload 8, empty datagram
		fuzzFrame(28, 0, nil),                                   // UDP length below the header (#589)
		fuzzFrame(28, 7, nil),                                   // UDP length below the header (#589)
		fuzzFrame(28, 11, []byte{1, 2, 3}),                      // UDP length past the IPv4 payload (#589)
		fuzzFrame(28+len(dhcp), modelUDPHeader+len(dhcp), dhcp), // a valid reply
		fuzzFrame(32, 12, trailing),                             // total length stops mid-frame (#292)
		fuzzFrame(36, 12, []byte{1, 2, 3, 4, 5, 6, 7, 8}),       // surplus area after the datagram (#589)
		fuzzFragment(modelIPv4MoreFragments, dhcp),              // first fragment (#591)
		fuzzFragment(1, dhcp),                                   // non-first fragment (#591)
		fuzzFragment(modelIPv4MoreFragments|1, dhcp),            // middle fragment (#591)
		fuzzFragment(0x4000, dhcp),                              // don't-fragment is not fragmentation
		{0x45, 0, 0, 40},                                        // header claims 40 bytes, frame truncated (#455)
		{0x4f, 0, 0, 60},                                        // IHL 15 header with nothing after it (#507)
		{0x45},                                                  // truncated IPv4 header
		{},                                                      // empty frame
	}
	for _, s := range seeds {
		f.Add(s, uint16(512))
	}
	valid := fuzzFrame(28+len(dhcp), modelUDPHeader+len(dhcp), dhcp)
	f.Add(valid, uint16(0)) // a valid reply into an empty buffer
	f.Add(valid, uint16(4)) // ... and one too small to hold it

	f.Fuzz(func(t *testing.T, raw []byte, bufLen uint16) {
		b := make([]byte, int(bufLen%2048))
		upc := &BroadcastRawUDPConn{
			PacketConn: &fuzzConn{frame: raw},
			boundAddr:  &net.UDPAddr{Port: 68},
		}
		n, addr, err := upc.ReadFrom(b)

		// The socket copies raw into a fixed-size receive buffer, so ReadFrom only
		// ever parses this prefix. That mirrors what the receive path does today
		// rather than a guarantee: a datagram larger than the caller's buffer is
		// dropped, not truncated to it.
		maxRecv := modelIPv4MaxHeader + modelUDPHeader + len(b)
		delivered := raw
		if len(delivered) > maxRecv {
			delivered = delivered[:maxRecv]
		}
		accept, wantN, payload, srcIP, srcPort := expectedReadFrom(delivered, len(b))

		if !accept {
			// A rejected frame drops out of the receive loop and surfaces as the
			// unwrapped io.EOF the mock returns once exhausted, with no bytes and no
			// address.
			if err == nil {
				t.Fatalf("accepted a frame the model rejects: n=%d delivered=%x", n, delivered)
			}
			if !errors.Is(err, io.EOF) {
				t.Fatalf("rejected frame: err=%v, want io.EOF; delivered=%x", err, delivered)
			}
			if n != 0 {
				t.Fatalf("rejected frame: n=%d, want 0; delivered=%x", n, delivered)
			}
			if addr != nil {
				t.Fatalf("rejected frame: addr=%v, want nil; delivered=%x", addr, delivered)
			}
			return
		}
		if err != nil {
			t.Fatalf("rejected a frame the model accepts: err=%v delivered=%x", err, delivered)
		}
		if n != wantN || n < 0 || n > len(b) {
			t.Fatalf("ReadFrom n=%d, want %d (buffer %d)", n, wantN, len(b))
		}
		if !bytes.Equal(b[:n], payload[:wantN]) {
			t.Fatalf("ReadFrom payload=%x, want %x", b[:n], payload[:wantN])
		}
		ua, ok := addr.(*net.UDPAddr)
		if !ok || ua == nil {
			t.Fatalf("ReadFrom addr=%v, want *net.UDPAddr", addr)
		}
		if !ua.IP.Equal(srcIP) {
			t.Fatalf("ReadFrom source IP=%v, want %v", ua.IP, srcIP)
		}
		if ua.Port != srcPort {
			t.Fatalf("ReadFrom source port=%d, want %d", ua.Port, srcPort)
		}
	})
}
