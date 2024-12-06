//go:build go1.12
// +build go1.12

package nclient4

import (
	"errors"
	"net"

	"github.com/u-root/uio/uio"
)

var (
	// ErrUDPAddrIsRequired is an error used when a passed argument is not of type "*net.UDPAddr".
	ErrUDPAddrIsRequired = errors.New("must supply UDPAddr")
)

func udpMatch(addr *net.UDPAddr, bound *net.UDPAddr) bool {
	if bound == nil {
		return true
	}
	if bound.IP != nil && !bound.IP.Equal(addr.IP) {
		return false
	}
	return bound.Port == addr.Port
}

func getUDP4pkt(pkt []byte, boundAddr *net.UDPAddr) ([]byte, *net.UDPAddr) {
	buf := uio.NewBigEndianBuffer(pkt)

	ipHdr := ipv4(buf.Data())

	if !ipHdr.isValid(len(pkt)) {
		return nil, nil
	}

	ipHdr = ipv4(buf.Consume(int(ipHdr.headerLength())))

	if ipHdr.transportProtocol() != udpProtocolNumber {
		return nil, nil
	}

	if !buf.Has(udpMinimumSize) {
		return nil, nil
	}

	udpHdr := udp(buf.Consume(udpMinimumSize))

	addr := &net.UDPAddr{
		IP:   ipHdr.destinationAddress(),
		Port: int(udpHdr.destinationPort()),
	}
	if !udpMatch(addr, boundAddr) {
		return nil, nil
	}
	srcAddr := &net.UDPAddr{
		IP:   ipHdr.sourceAddress(),
		Port: int(udpHdr.sourcePort()),
	}
	// Extra padding after end of IP packet should be ignored,
	// if not dhcp option parsing will fail.
	dhcpLen := int(ipHdr.payloadLength()) - udpMinimumSize
	return buf.Consume(dhcpLen), srcAddr
}
