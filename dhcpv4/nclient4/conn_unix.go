//go:build go1.12 && (darwin || freebsd || netbsd || openbsd || dragonfly)
// +build go1.12
// +build darwin freebsd netbsd openbsd dragonfly

package nclient4

import (
	"io"
	"net"

	"github.com/mdlayher/raw"
)

const (
	bpfFilterBidirectional int = 1
)

var rawConnectionConfig = &raw.Config{
	BPFDirection: bpfFilterBidirectional,
}

// NewRawUDPConn returns a UDP connection bound to the interface and port
// given based on a raw packet socket. All packets are broadcasted.
func NewRawUDPConn(iface string, port int, vlans ...uint16) (net.PacketConn, error) {
	ifc, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}

	var etherType uint16
	if len(vlans) > 0 {
		etherType = vlanTPID // The VLAN TPID field is located in the same offset as EtherType
	} else {
		etherType = etherIPv4Proto
	}

	rawConn, err := raw.ListenPacket(ifc, etherType, rawConnectionConfig)
	if err != nil {
		return nil, err
	}

	return NewBroadcastUDPConn(net.PacketConn(rawConn), &net.UDPAddr{Port: port}, vlans...), nil
}

type BroadcastRawUDPConn struct {
	// PacketConn is a raw network socket
	net.PacketConn

	boundAddr *net.UDPAddr
	// VLAN tags can be configured to make up for the shortcoming of the BSD implementation
	VLANs []uint16
}

// NewBroadcastUDPConn returns a PacketConn that marshals and unmarshals UDP
// packets, sending them to the broadcast MAC at on rawPacketConn.
// Supplied VLAN tags are inserted into the Ethernet frame before sending.
//
// Calls to ReadFrom will only return packets destined to boundAddr.
func NewBroadcastUDPConn(rawPacketConn net.PacketConn, boundAddr *net.UDPAddr, vlans ...uint16) net.PacketConn {
	return &BroadcastRawUDPConn{
		PacketConn: rawPacketConn,
		boundAddr:  boundAddr,
		VLANs:      vlans,
	}
}

// ReadFrom implements net.PacketConn.ReadFrom.
//
// ReadFrom reads raw Ethernet packets, parses the VLAN stack (if configured)
// and will try to match the IP+UDP destinations against upc.boundAddr.
//
// Any matching packets are returned via the given buffer.
func (upc *BroadcastRawUDPConn) ReadFrom(b []byte) (int, net.Addr, error) {
	ethHdrLen := ethHdrMinimum
	if len(upc.VLANs) > 0 {
		ethHdrLen += len(upc.VLANs) * vlanTagLen
	}
	ipHdrMaxLen := ipv4MaximumHeaderSize
	udpHdrLen := udpMinimumSize

	for {
		pkt := make([]byte, ethHdrLen+ipHdrMaxLen+udpHdrLen+len(b))
		n, _, err := upc.PacketConn.ReadFrom(pkt)
		if err != nil {
			return 0, nil, err
		}
		if n == 0 {
			return 0, nil, io.EOF
		}

		pkt = getEthernetPayload(pkt[:n], upc.VLANs)
		if pkt == nil {
			// VLAN stack does not match our configuration
			continue
		}
		dhcpPkt, srcAddr := getUDP4pkt(pkt[:n], upc.boundAddr)
		if dhcpPkt == nil {
			continue
		}

		return copy(b, dhcpPkt), srcAddr, nil
	}
}

// WriteTo implements net.PacketConn.WriteTo and broadcasts all packets at the
// raw socket level.
//
// WriteTo wraps the given packet in the appropriate UDP, IP and Ethernet header
// before sending it on the packet conn. Since the Ethernet encapsulation is done
// on the application's side, VLAN tagging also has to be handled in the application.
func (upc *BroadcastRawUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		return 0, ErrUDPAddrIsRequired
	}

	// Using the boundAddr is not quite right here, but it works.
	pkt := udp4pkt(b, udpAddr, upc.boundAddr)

	srcMac := upc.PacketConn.LocalAddr().(*raw.Addr).HardwareAddr
	pkt = addEthernetHdr(pkt, BroadcastMac, srcMac, etherIPv4Proto, upc.VLANs)

	// The `raw` packet connection does not take any address as an argument.
	return upc.PacketConn.WriteTo(pkt, nil)
}
