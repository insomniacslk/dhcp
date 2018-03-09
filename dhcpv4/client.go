package dhcpv4

import (
	"encoding/binary"
	"net"
	"syscall"
	"time"

	"golang.org/x/net/ipv4"
)

const (
	MaxUDPReceivedPacketSize = 8192 // arbitrary size. Theoretically could be up to 65kb
)

// just triggering the golint complainer..
type Client struct {
	Network string
	Dialer  *net.Dialer
	Timeout time.Duration
}

func MakeRawBroadcastPacket(payload []byte) ([]byte, error) {
	udp := make([]byte, 8)
	binary.BigEndian.PutUint16(udp[:2], ClientPort)
	binary.BigEndian.PutUint16(udp[2:4], ServerPort)
	binary.BigEndian.PutUint16(udp[4:6], uint16(8+len(payload)))
	binary.BigEndian.PutUint16(udp[6:8], 0) // try to offload the checksum

	h := ipv4.Header{
		Version:  4,
		Len:      20,
		TotalLen: 20 + len(udp) + len(payload),
		TTL:      64,
		Protocol: 17, // UDP
		Dst:      net.IPv4bcast,
		Src:      net.IPv4zero,
	}
	ret, err := h.Marshal()
	if err != nil {
		return nil, err
	}
	ret = append(ret, udp...)
	ret = append(ret, payload...)
	return ret, nil
}

// Run a full DORA transaction: Discovery, Offer, Request, Acknowledge, over
// UDP. Does not retry in case of failures.
// Returns a list of DHCPv4 structures representing the exchange. It can contain
// up to four elements, ordered as Discovery, Offer, Request and Acknowledge.
// In case of errors, an error is returned, and the list of DHCPv4 objects will
// be shorted than 4, containing all the sent and received DHCPv4 messages.
func (c *Client) Exchange(ifname string, d *DHCPv4) ([]DHCPv4, error) {
	conversation := make([]DHCPv4, 1)
	var err error

	// Discovery
	if d == nil {
		d, err = NewDiscoveryForInterface(ifname)
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
	packet, err := MakeRawBroadcastPacket(d.ToBytes())
	if err != nil {
		return conversation, err
	}
	err = syscall.Sendto(fd, packet, 0, &daddr)
	if err != nil {
		return conversation, err
	}

	// Offer
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: ClientPort})
	if err != nil {
		return conversation, err
	}
	defer conn.Close()

	buf := make([]byte, MaxUDPReceivedPacketSize)
	oobdata := []byte{} // ignoring oob data
	n, _, _, _, err := conn.ReadMsgUDP(buf, oobdata)
	offer, err := FromBytes(buf[:n])
	if err != nil {
		return conversation, err
	}
	// TODO match the packet content
	// TODO check that the peer address matches the declared server IP and port
	conversation = append(conversation, *offer)

	// Request
	request, err := RequestFromOffer(*offer)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, *request)
	packet, err = MakeRawBroadcastPacket(request.ToBytes())
	if err != nil {
		return conversation, err
	}
	err = syscall.Sendto(fd, packet, 0, &daddr)
	if err != nil {
		return conversation, err
	}

	// Acknowledge
	buf = make([]byte, MaxUDPReceivedPacketSize)
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
