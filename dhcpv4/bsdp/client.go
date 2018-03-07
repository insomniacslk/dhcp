// +build darwin

package bsdp

import (
	"fmt"
	"net"
	"syscall"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Client is a BSDP-specific client suitable for performing BSDP exchanges.
type Client dhcpv4.Client

// Exchange runs a full BSDP exchange (Inform[list], Ack, Inform[select],
// Ack). Returns a list of DHCPv4 structures representing the exchange.
func (c *Client) Exchange(ifname string, informList *dhcpv4.DHCPv4) ([]dhcpv4.DHCPv4, error) {
	conversation := make([]dhcpv4.DHCPv4, 1)
	var err error

	// INFORM[LIST]
	if informList == nil {
		informList, err = NewInformListForInterface(ifname, dhcpv4.ClientPort)
		if err != nil {
			return conversation, err
		}
	}
	conversation[0] = *informList

	// TODO: deduplicate with code in dhcpv4/client.go
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
	err = dhcpv4.BindToInterface(fd, ifname)
	if err != nil {
		return conversation, err
	}

	bcast := [4]byte{}
	copy(bcast[:], net.IPv4bcast)
	daddr := syscall.SockaddrInet4{Port: dhcpv4.ClientPort, Addr: bcast}
	packet, err := dhcpv4.MakeRawBroadcastPacket(informList.ToBytes())
	if err != nil {
		return conversation, err
	}
	err = syscall.Sendto(fd, packet, 0, &daddr)
	if err != nil {
		return conversation, err
	}

	// ACK 1
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: dhcpv4.ClientPort})
	if err != nil {
		return conversation, err
	}
	defer conn.Close()

	buf := make([]byte, dhcpv4.MaxUDPReceivedPacketSize)
	oobdata := []byte{} // ignoring oob data
	n, _, _, _, err := conn.ReadMsgUDP(buf, oobdata)
	ack1, err := dhcpv4.FromBytes(buf[:n])
	if err != nil {
		return conversation, err
	}
	// TODO match the packet content
	// TODO check that the peer address matches the declared server IP and port
	conversation = append(conversation, *ack1)

	// Parse boot images sent back by server
	bootImages, err := ParseBootImageListFromAck(*ack1)
	if err != nil {
		return conversation, err
	}
	if len(bootImages) == 0 {
		return conversation, fmt.Errorf("Got no BootImages from server")
	}

	// INFORM[SELECT]
	informSelect, err := InformSelectForAck(*ack1, dhcpv4.ClientPort, bootImages[0])
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, *informSelect)
	packet, err = dhcpv4.MakeRawBroadcastPacket(informSelect.ToBytes())
	if err != nil {
		return conversation, err
	}
	err = syscall.Sendto(fd, packet, 0, &daddr)
	if err != nil {
		return conversation, err
	}

	// ACK 2
	buf = make([]byte, dhcpv4.MaxUDPReceivedPacketSize)
	n, _, _, _, err = conn.ReadMsgUDP(buf, oobdata)
	ack2, err := dhcpv4.FromBytes(buf[:n])
	if err != nil {
		return conversation, err
	}
	// TODO match the packet content
	// TODO check that the peer address matches the declared server IP and port
	conversation = append(conversation, *ack2)

	return conversation, nil
}
