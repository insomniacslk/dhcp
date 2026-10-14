// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || freebsd || linux || netbsd || openbsd || dragonfly

package nclient4

import (
	"encoding/binary"
	"net"
	"testing"
)

// TestUDP4PktEncodesCalculatedZeroChecksumAsAllOnes pins the wire field rather
// than recomputing it, so the model cannot drift with the encoder.
func TestUDP4PktEncodesCalculatedZeroChecksumAsAllOnes(t *testing.T) {
	// The one's-complement sum of this datagram is all ones, so the checksum
	// calculation yields zero.
	pkt := udp4pkt([]byte{0xff, 0x53},
		&net.UDPAddr{IP: net.IPv4bcast, Port: ServerPort},
		&net.UDPAddr{IP: net.IPv4zero, Port: ClientPort})

	const off = ipv4MinimumSize + udpchecksum
	if got := binary.BigEndian.Uint16(pkt[off:]); got != 0xffff {
		t.Fatalf("UDP checksum = %#04x, want 0xffff", got)
	}
}
