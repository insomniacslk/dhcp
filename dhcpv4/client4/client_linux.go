package client4

import (
	"net"

	"golang.org/x/sys/unix"
)

func makeListeningSocketWithCustomPort(ifname string, port int) (int, error) {
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_DGRAM, int(htons(unix.ETH_P_IP)))
	if err != nil {
		return fd, err
	}
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return fd, err
	}
	llAddr := unix.SockaddrLinklayer{
		Ifindex:  iface.Index,
		Protocol: htons(unix.ETH_P_IP),
	}
	err = unix.Bind(fd, &llAddr)
	return fd, err
}
