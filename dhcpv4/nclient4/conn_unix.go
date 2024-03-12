//go:build go1.12 && (darwin || freebsd || netbsd || openbsd || dragonfly)
// +build go1.12
// +build darwin freebsd netbsd openbsd dragonfly

package nclient4

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/mdlayher/raw"
	"github.com/u-root/uio/uio"
)

const (
	bpfFilterBidirectional int = 1

	etherIPv4Proto uint16 = 0x0800
	ethHdrMinimum  int    = 14

	vlanTagLen int    = 4
	vlanMax    uint16 = 0x0FFF
	vlanTPID   uint16 = 0x8100
)

var rawConnectionConfig = &raw.Config{
	BPFDirection: bpfFilterBidirectional,
}

// NewRawUDPConn returns a UDP connection bound to the interface and port
// given based on a raw packet socket. All packets are broadcasted.
//
// The interface can be completely unconfigured.
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

// processVLANStack receives a buffer starting at the first TPID/EtherType field, and walks through
// the VLAN stack until either an unexpected VLAN is found, or if an IPv4 EtherType is found.
// The data from the provided buffer is consumed until the end of the Ethernet header
//
// processVLANStack returns true if the VLAN stack in the packet corresponds to the VLAN configuration, false otherwise
func processVLANStack(buf *uio.Lexer, vlans []uint16) bool {
	var currentVLAN uint16
	var vlanStackIsCorrect bool
	configuredVLANs := make([]uint16, len(vlans))
	copy(configuredVLANs, vlans)

	for {
		switch etherType := binary.BigEndian.Uint16(buf.Consume(2)); etherType {
		case vlanTPID:
			tci := binary.BigEndian.Uint16(buf.Consume(2))
			vlanID := tci & vlanMax // Mask first 4 bytes
			if len(configuredVLANs) != 0 {
				currentVLAN, configuredVLANs = configuredVLANs[0], configuredVLANs[1:]
				if vlanID != currentVLAN {
					// Packet VLAN tag does not match configured VLAN stack
					vlanStackIsCorrect = false
				}
			} else {
				// Packet VLAN stack is too long
				vlanStackIsCorrect = false
			}
		case etherIPv4Proto:
			if len(configuredVLANs) == 0 {
				// Packet VLAN stack has been correctly consumed
				vlanStackIsCorrect = true
			} else {
				// VLAN tags remaining in configured stack -> not a match
				vlanStackIsCorrect = false
			}
			return vlanStackIsCorrect
		default:
			vlanStackIsCorrect = false
			return vlanStackIsCorrect
		}
	}
}

func getEthernetPayload(pkt []byte, vlans []uint16) []byte {
	buf := uio.NewBigEndianBuffer(pkt)
	dstMac := buf.Consume(6)
	srcMac := buf.Consume(6)
	_, _ = dstMac, srcMac

	if len(vlans) > 0 {
		success := processVLANStack(buf, vlans)
		if !success {
			return nil
		}
	} else {
		etherType := binary.BigEndian.Uint16(buf.Consume(2))
		if etherType != etherIPv4Proto {
			return nil
		}
	}

	return buf.Data()
}

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

		// We're only interested in properly tagged packets
		pkt = getEthernetPayload(pkt[:n], upc.VLANs)
		if pkt == nil {
			continue
		}
		dhcpPkt, srcAddr := getUDP4pkt(pkt[:n], upc.boundAddr)
		if dhcpPkt == nil {
			continue
		}

		return copy(b, dhcpPkt), srcAddr, nil
	}
}

// createVLANTag returns the bytes of a 4-byte long VLAN tag, which can be inserted
// in an Ethernet frame header.
func createVLANTag(vlan uint16) []byte {
	vlanTag := make([]byte, vlanTagLen)
	// First 2 bytes are the TPID. Only support 802.1Q for now (even for QinQ, 802.1ad is rarely used)
	binary.BigEndian.PutUint16(vlanTag, vlanTPID)

	var pcp, dei, tci uint16
	// TCI - tag control information, 2 bytes. Format: | PCP (3 bits) | DEI (1 bit) | VLAN ID (12 bits) |
	pcp = 0x0 // 802.1p priority level - 3 bits, valid values range from 0x0 to 0x7. 0x0 - best effort
	dei = 0x0 // drop eligible indicator - 1 bit, valid values are 0x0 or 0x1. 0x0 - not drop eligible
	tci |= pcp << 13
	tci |= dei << 12
	tci |= vlan
	binary.BigEndian.PutUint16(vlanTag[2:], tci)

	return vlanTag
}

// addEthernetHdr returns the supplied packet (in bytes) with an
// added Ethernet header with the specified EtherType.
func addEthernetHdr(b []byte, dstMac, srcMac net.HardwareAddr, etherProto uint16, vlans []uint16) []byte {
	ethHdrLen := ethHdrMinimum
	if len(vlans) > 0 {
		ethHdrLen += len(vlans) * vlanTagLen
	}
	b = append(make([]byte, ethHdrLen), b...)
	offset := 0
	copy(b, dstMac)
	offset += len(dstMac)
	copy(b[offset:], srcMac)
	offset += len(srcMac)
	for _, vlan := range vlans {
		copy(b[offset:], createVLANTag(vlan))
		offset += vlanTagLen
	}

	binary.BigEndian.PutUint16(b[offset:], etherProto)

	return b
}

// WriteTo implements net.PacketConn.WriteTo and broadcasts all packets at the
// raw socket level.
//
// WriteTo wraps the given packet in the appropriate UDP, IP and Ethernet header
// before sending it on the packet conn. Since the Ethernet encapsulation is done
// on the application's side, this implementation does not work well with VLAN
// tagging and such.
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
