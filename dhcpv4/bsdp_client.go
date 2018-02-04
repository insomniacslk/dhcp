// +build darwin

package dhcpv4

import (
	"fmt"
	"net"
	"syscall"
)

// BSDPExchange runs a full BSDP exchange (Inform[list], Ack, Inform[select],
// Ack). Returns a list of DHCPv4 structures representing the exchange.
func (c *Client) BSDPExchange(ifname string, d *DHCPv4) ([]DHCPv4, error) {
	conversation := make([]DHCPv4, 1)
	var err error

	// INFORM[LIST]
	if d == nil {
		d, err = NewInformListForInterface(ifname, ClientPort)
	}
	conversation[0] = *d

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return conversation, err
	}
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return conversation, err
	}
	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	if err != nil {
		return conversation, err
	}
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	if err != nil {
		return conversation, err
	}
	err = BindToInterface(fd, ifname)
	if err != nil {
		return conversation, err
	}

	daddr := syscall.SockaddrInet4{Port: ClientPort, Addr: [4]byte{255, 255, 255, 255}}
	packet, err := makeRawBroadcastPacket(d.ToBytes())
	if err != nil {
		return conversation, err
	}
	err = syscall.Sendto(fd, packet, 0, &daddr)
	if err != nil {
		return conversation, err
	}

	// ACK 1
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: ClientPort})
	if err != nil {
		return conversation, err
	}
	defer conn.Close()

	buf := make([]byte, maxUDPReceivedPacketSize)
	oobdata := []byte{} // ignoring oob data
	n, _, _, _, err := conn.ReadMsgUDP(buf, oobdata)
	ack1, err := FromBytes(buf[:n])
	if err != nil {
		return conversation, err
	}
	// TODO match the packet content
	// TODO check that the peer address matches the declared server IP and port
	conversation = append(conversation, *ack1)

	// Parse boot images sent back by server
	bootImages, err := ParseBootImageListFromAck(*ack1)
	fmt.Println(bootImages)
	if err != nil {
		return conversation, err
	}
	if len(bootImages) == 0 {
		return conversation, fmt.Errorf("Got no BootImages from server")
	}

	// INFORM[SELECT]
	request, err := InformSelectForAck(*ack1, ClientPort, bootImages[0])
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, *request)
	packet, err = makeRawBroadcastPacket(request.ToBytes())
	if err != nil {
		return conversation, err
	}
	err = syscall.Sendto(fd, packet, 0, &daddr)
	if err != nil {
		return conversation, err
	}

	// ACK 2
	buf = make([]byte, maxUDPReceivedPacketSize)
	n, _, _, _, err = conn.ReadMsgUDP(buf, oobdata)
	acknowledge, err := FromBytes(buf[:n])
	if err != nil {
		return conversation, err
	}
	// TODO match the packet content
	// TODO check that the peer address matches the declared server IP and port
	conversation = append(conversation, *acknowledge)

	return conversation, nil
}
