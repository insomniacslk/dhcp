//go:build go1.12
// +build go1.12

package nclient4

import (
	"encoding/binary"
	"net"

	"github.com/u-root/uio/uio"
)

const (
	etherIPv4Proto uint16 = 0x0800
	ethHdrBaseLen  int    = 14

	vlanTagLen int    = 4
	vlanMax    uint16 = 0x0FFF
	vlanTPID   uint16 = 0x8100
)

var (
	// BroadcastMac is the broadcast MAC address.
	//
	// Any UDP packet sent to this address is broadcast on the subnet.
	BroadcastMac = net.HardwareAddr([]byte{255, 255, 255, 255, 255, 255})
)

// processVLANStack receives a buffer starting at the first TPID/EtherType field, and walks through
// the VLAN stack until either an unexpected VLAN is found, or if an IPv4 EtherType is found.
// The data from the provided buffer is consumed until the end of the Ethernet header.
//
// processVLANStack returns true if the VLAN stack in the packet corresponds to the VLAN configuration, false otherwise.
func processVLANStack(buf *uio.Lexer, vlans []uint16) bool {
	var currentVLAN uint16
	var vlanStackIsCorrect bool = true
	configuredVLANs := make([]uint16, len(vlans))
	copy(configuredVLANs, vlans)

	for {
		switch etherType := binary.BigEndian.Uint16(buf.Consume(2)); etherType {
		case vlanTPID:
			tci := binary.BigEndian.Uint16(buf.Consume(2))
			vlanID := tci & vlanMax // Mask first 4 bits
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
				// Packet VLAN stack has been consumed, return result
				return vlanStackIsCorrect
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

// getEthernetPayload processes an Ethernet header, verifies the
// VLAN tags contained in it and returns the payload as a byte slice.
//
// If the VLAN tag stack does not match the VLAN configuration,
// nil is returned (since the packet is not meant for us).
// In case the EtherType does not match the IPv4 proto value,
// nil is returned too (since the packet could not be DHCPv4).
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

// createVLANTag returns the bytes of a 4-byte long 802.1Q VLAN tag,
// which can be inserted in an Ethernet frame header.
func createVLANTag(vlan uint16) []byte {
	vlanTag := make([]byte, vlanTagLen)
	// First 2 bytes are the TPID. Only support 802.1Q for now (even for QinQ, 802.1ad is rarely used)
	binary.BigEndian.PutUint16(vlanTag, vlanTPID)

	var pcp, dei, tci uint16
	// TCI - tag control information, 2 bytes. Format: | PCP (3 bits) | DEI (1 bit) | VLAN ID (12 bits) |
	pcp = 0x0        // 802.1p priority level - 3 bits, valid values range from 0x0 to 0x7. 0x0 - best effort
	dei = 0x0        // drop eligible indicator - 1 bit, valid values are 0x0 or 0x1. 0x0 - not drop eligible
	tci |= pcp << 13 // 16-3 = 13 offset
	tci |= dei << 12 // 13-1 = 12 offset
	tci |= vlan      // VLAN ID (VID) is 12 bits
	binary.BigEndian.PutUint16(vlanTag[2:], tci)

	return vlanTag
}

// addEthernetHdr returns the supplied packet (in bytes) with an
// added Ethernet header with the specified EtherType.
func addEthernetHdr(b []byte, dstMac, srcMac net.HardwareAddr, etherProto uint16, vlans []uint16) []byte {
	ethHdrLen := ethHdrBaseLen
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
