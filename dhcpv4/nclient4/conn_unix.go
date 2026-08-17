// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.12 && (darwin || freebsd || linux || netbsd || openbsd || dragonfly)
// +build go1.12
// +build darwin freebsd linux netbsd openbsd dragonfly

package nclient4

import (
	"errors"
	"io"
	"net"

	"github.com/mdlayher/packet"
	"github.com/u-root/uio/uio"
	"golang.org/x/net/bpf"
	"golang.org/x/sys/unix"
)

var (
	// BroadcastMac is the broadcast MAC address.
	//
	// Any UDP packet sent to this address is broadcast on the subnet.
	BroadcastMac = net.HardwareAddr([]byte{255, 255, 255, 255, 255, 255})
)

var (
	// ErrUDPAddrIsRequired is an error used when a passed argument is not of type "*net.UDPAddr".
	ErrUDPAddrIsRequired = errors.New("must supply UDPAddr")

	// errPortOutOfRange is returned by NewRawUDPConn when port is not a uint16.
	errPortOutOfRange = errors.New("port out of range")
)

// NewRawUDPConn returns a UDP connection bound to the interface and port
// given based on a raw packet socket. All packets are broadcasted.
//
// The interface can be completely unconfigured.
func NewRawUDPConn(iface string, port int) (net.PacketConn, error) {
	if port < 0 || port > 0xffff {
		return nil, errPortOutOfRange
	}
	ifc, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}
	rawConn, err := listenFiltered(ifc, uint16(port))
	if err != nil {
		return nil, err
	}
	return NewBroadcastUDPConn(rawConn, &net.UDPAddr{Port: port}), nil
}

// listenPacket is used for testing purposes
var listenPacket = func(ifc *net.Interface, cfg *packet.Config) (net.PacketConn, error) {
	return packet.Listen(ifc, packet.Datagram, unix.ETH_P_IP, cfg)
}

// listenFiltered opens the raw socket with the prefilter already applied.
// Config.Filter takes effect before bind, so nothing is queued unfiltered; a
// kernel that refuses it still gets a plain socket, since ReadFrom revalidates.
func listenFiltered(ifc *net.Interface, clientPort uint16) (net.PacketConn, error) {
	var filterErr error
	if raw, err := bpf.Assemble(dhcpClientFilter(clientPort)); err != nil {
		filterErr = err
	} else if conn, err := listenPacket(ifc, &packet.Config{Filter: raw}); err != nil {
		filterErr = err
	} else {
		return conn, nil
	}

	conn, err := listenPacket(ifc, nil)
	if err != nil {
		// The filtered failure is usually the more specific one.
		return nil, errors.Join(filterErr, err)
	}
	return conn, nil
}

// dhcpClientFilter builds a classic-BPF program passing IPv4/UDP frames destined
// to clientPort. On the SOCK_DGRAM socket offset 0 is the IPv4 header, the same
// bytes ReadFrom parses; the filter only prefilters, ReadFrom still validates.
func dhcpClientFilter(clientPort uint16) []bpf.Instruction {
	const protocolUDP = 17
	return []bpf.Instruction{
		bpf.LoadAbsolute{Off: 0, Size: 1},                                     // version + IHL
		bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: 0x0f},                        // A = IHL
		bpf.JumpIf{Cond: bpf.JumpLessThan, Val: 5, SkipTrue: 5},               // IHL < 5: drop
		bpf.LoadAbsolute{Off: 9, Size: 1},                                     // IP protocol
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: protocolUDP, SkipTrue: 3},     // not UDP: drop
		bpf.LoadMemShift{Off: 0},                                              // X = IHL * 4
		bpf.LoadIndirect{Off: 2, Size: 2},                                     // UDP destination port
		bpf.JumpIf{Cond: bpf.JumpEqual, Val: uint32(clientPort), SkipTrue: 1}, // matches: accept
		bpf.RetConstant{Val: 0},                                               // drop
		bpf.RetConstant{Val: 262144},                                          // accept
	}
}

// BroadcastRawUDPConn uses a raw socket to send UDP packets to the broadcast
// MAC address.
type BroadcastRawUDPConn struct {
	// PacketConn is a raw DGRAM socket.
	net.PacketConn

	// boundAddr is the address this RawUDPConn is "bound" to.
	//
	// Calls to ReadFrom will only return packets destined to this address.
	boundAddr *net.UDPAddr
}

// NewBroadcastUDPConn returns a PacketConn that marshals and unmarshals UDP
// packets, sending them to the broadcast MAC at on rawPacketConn.
//
// Calls to ReadFrom will only return packets destined to boundAddr.
func NewBroadcastUDPConn(rawPacketConn net.PacketConn, boundAddr *net.UDPAddr) net.PacketConn {
	return &BroadcastRawUDPConn{
		PacketConn: rawPacketConn,
		boundAddr:  boundAddr,
	}
}

func udpMatch(addr *net.UDPAddr, bound *net.UDPAddr) bool {
	if bound == nil {
		return true
	}
	if bound.IP != nil && !bound.IP.Equal(addr.IP) {
		return false
	}
	return bound.Port == addr.Port
}

// ReadFrom implements net.PacketConn.ReadFrom.
//
// ReadFrom reads raw IP packets and will try to match them against
// upc.boundAddr. Any matching packets are returned via the given buffer.
func (upc *BroadcastRawUDPConn) ReadFrom(b []byte) (int, net.Addr, error) {
	ipHdrMaxLen := ipv4MaximumHeaderSize
	udpHdrLen := udpMinimumSize

	for {
		pkt := make([]byte, ipHdrMaxLen+udpHdrLen+len(b))
		n, _, err := upc.PacketConn.ReadFrom(pkt)
		if err != nil {
			return 0, nil, err
		}
		if n == 0 {
			return 0, nil, io.EOF
		}
		pkt = pkt[:n]
		buf := uio.NewBigEndianBuffer(pkt)

		ipHdr := ipv4(buf.Data())

		if !ipHdr.isValid(n) {
			continue
		}

		ipHdr = ipv4(buf.Consume(int(ipHdr.headerLength())))

		if ipHdr.transportProtocol() != udpProtocolNumber {
			continue
		}

		// The IPv4 payload must be large enough to hold a UDP header. Discarding
		// frames that declare fewer bytes keeps the DHCP length computed below
		// non-negative and avoids parsing trailing padding as a UDP header.
		ipPayloadLen := int(ipHdr.payloadLength())
		if ipPayloadLen < udpHdrLen {
			continue
		}

		if !buf.Has(udpHdrLen) {
			continue
		}

		udpHdr := udp(buf.Consume(udpHdrLen))

		addr := &net.UDPAddr{
			IP:   ipHdr.destinationAddress(),
			Port: int(udpHdr.destinationPort()),
		}
		if !udpMatch(addr, upc.boundAddr) {
			continue
		}
		srcAddr := &net.UDPAddr{
			IP:   ipHdr.sourceAddress(),
			Port: int(udpHdr.sourcePort()),
		}
		// Extra padding after the end of the IP payload is ignored; otherwise
		// dhcp option parsing would fail.
		dhcpLen := ipPayloadLen - udpHdrLen
		if !buf.Has(dhcpLen) {
			continue
		}
		return copy(b, buf.Consume(dhcpLen)), srcAddr, nil
	}
}

// WriteTo implements net.PacketConn.WriteTo and broadcasts all packets at the
// raw socket level.
//
// WriteTo wraps the given packet in the appropriate UDP and IP header before
// sending it on the packet conn.
func (upc *BroadcastRawUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		return 0, ErrUDPAddrIsRequired
	}

	// Using the boundAddr is not quite right here, but it works.
	pkt := udp4pkt(b, udpAddr, upc.boundAddr)

	// Broadcasting is not always right, but hell, what the ARP do I know.
	return upc.PacketConn.WriteTo(pkt, &packet.Addr{HardwareAddr: BroadcastMac})
}
