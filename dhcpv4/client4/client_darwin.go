package client4

//import "net"
//import "golang.org/x/sys/unix"
import "fmt"

// //https://stackoverflow.com/questions/17169298/af-packet-on-osx

// _INET, SOCK_RAW, IPPROTO_IP);
// And then you need to tell the system, that you want to provide your own IP header:

// int yes = 1;
// setsockopt(soc, IPPROTO_IP, IP_HDRINCL, &yes, sizeof(yes));
//
// Now you can send raw IP packets (e.g. IP header + UDP header +
//payload data) to the socket for sending, however, depending on your
//system, the system will perform some sanity checks and maybe
//override some fields in the header. E.g. it may not allow you to
//create malformed IP packets or prevent you from performing IP
//address spoofing. Therefor it may for example calculate the IPv4
//header checksum for you or automatically fill in the correct source
//address if your IP header uses 0.0.0.0 or :: as source
//address. Check the man page for ip(4) or for raw(7) on your target
//system. Apple doesn't ship programmer man pages for macOS any
//longer, but you can find them online.

// To quote from that man page:

// Unlike previous BSD releases, the program must set all the fields
//of the IP header, including the following:

//  ip->ip_v = IPVERSION;
//   ip->ip_hl = hlen >> 2;
//    ip->ip_id = 0;  /* 0 means kernel set appropriate value */
//     ip->ip_off = offset;
//      ip->ip_len = len;
//      Note that the ip_off and ip_len fields are in host byte order.

// If the header source address is set to INADDR_ANY, the kernel will
//choose an appropriate address.

// Note that ip_sum is not mentioned at all, so apparently you don't
//have to provide that one and the system will always calculate it for
//you.

// If you compare that to Linux raw(7):

// ┌───────────────────────────────────────────────────┐
// │IP Header fields modified on sending by IP_HDRINCL │
// ├──────────────────────┬────────────────────────────┤
// │IP Checksum           │ Always filled in           │
// ├──────────────────────┼────────────────────────────┤
// │Source Address        │ Filled in when zero        │
// ├──────────────────────┼────────────────────────────┤
// │Packet ID             │ Filled in when zero        │
// ├──────────────────────┼────────────────────────────┤
// │Total Length          │ Always filled in           │
// └──────────────────────┴────────────────────────────┘
//
// When receiving from a raw IP socket, you will either get all
//incoming IP packets that arrive at the host or just a subset of them
//(e.g. Windows does support raw sockets but won't ever let you send
//or receive TCP packets). You will receive the full packet, including
//all headers, so the first byte of every packet received is the first
//byte of the IP header.

// Some people here will ask why I use IPPROTO_IP and not
//IPPROTO_RAW. When using IPPROTO_RAW you don't have to set
//IP_HDRINCL:

// A protocol of IPPROTO_RAW implies enabled IP_HDRINCL and is able to
//send any IP protocol that is specified in the passed header.

// But you can only use IPPROTO_RAW for outgoing traffic:

// An IPPROTO_RAW socket is send only.

// On macOS you can use IPPROTO_IP and you will receive all IP packets
//but on Linux this may not work, hence the created a new socket
//PF_PACKET socket type. What should work on both systems is
//specifying a sub-protocol:

// int soc = socket(PF_INET, SOCK_RAW, IPPROTO_UDP); Of course, now
//you can only send/receive UDP packets over that socket. If you set
//IP_HDRINCL again, you need to provide a full IP header on send and
//you will receive a full IP header on receive. If you don't set it,
//you can just provide the UDP header on send and the system will add
//an IP header itself, that is, if the socket is connected and
//optionally bound, so the system knows which addresses to use in that
//header. For receiving that option plays no role, you always get the
//IP header for every UDP packet you receive on such a socket.

// In case people wonder why I use PF_INET and not AF_INET: PF means
//Protocol Family and AF means Address Family. Usually these are the
//same (e.g. AF_INET == PF_INET) so it won't matter what you use, but
//strictly speaking sockets should be creates with PF_ and the family
//in sockaddr structures should be set with AF_ as one day there might
//be a protocol that supports two kind of different addresses and then
//there will be AF_XXX1 and AF_XXX2 and neither one may be the same as
//PF_XXX.~

func makeListeningSocketWithCustomPort(ifname string, port int) (int, error) {
	// var  soc = unix.Socket(unix.PF_INET, unix.SOCK_RAW, unix.IPPROTO_IP)
	//      yes := uintptr(1)
	//      setsockopt(soc, unix.IPPROTO_IP, unix.IP_HDRINCL, &yes, sizeof(yes));
	//      return unix.Socket(unix.AF_PACKET, unix.SOCK_DGRAM, int(htons(unix.ETH_P_IP)))

	// fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_DGRAM, int(htons(unix.ETH_P_IP)))
	// 	if err != nil {
	// 		return fd, err
	// 	}
	// 	iface, err := net.InterfaceByName(ifname)
	// 	if err != nil {
	// 		return fd, err
	// 	}
	// 	llAddr := unix.SockaddrLinklayer{
	// 		Ifindex:  iface.Index,
	// 		Protocol: htons(unix.ETH_P_IP),
	// 	}
	// 	err = unix.Bind(fd, &llAddr)
	// 	return fd, err
	return -1, fmt.Errorf("not yet")
}
