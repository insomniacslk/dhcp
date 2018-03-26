package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
)

var ver = flag.Int("v", 6, "IP version to use")
var infile = flag.String("r", "", "PCAP file to read from. If not specified, try to send an actual DHCP request")
var iface = flag.String("i", "eth0", "Network interface to send packets through")
var useEtherIP = flag.Bool("etherip", false, "Enables LayerTypeEtherIP instead of LayerTypeEthernet, use with linux-cooked PCAP files. (default: false)")
var debug = flag.Bool("debug", false, "Enable debug output (default: false)")
var live = flag.Bool("live", false, "Sniff DHCP packets from the network (default: false)")
var snaplen = flag.Int("s", 0, "Set the snaplen when using -live (default: 0)")
var count = flag.Int("c", 0, "Stop after <count> packets (default: 0)")
var unpack = flag.Bool("unpack", false, "Unpack inner DHCPv6 messages when parsing relay messages")
var to = flag.String("to", "", "Destination to send packets to. If empty, will use [ff02::1:2]:547")

// Clientv4 runs a DHCPv4 client and prints out a summary of the results.
func Clientv4() {
	client := dhcpv4.NewClient()
	conv, err := client.Exchange(*iface, nil)
	// don't exit immediately if there's an error, since `conv` will always
	// contain at least the SOLICIT message. So print it out first
	for _, m := range conv {
		log.Print(m.Summary())
	}
	if err != nil {
		log.Fatal(err)
	}
}

// Clientv6 runs a DHCPv6 client and prints out a summary of the results.
func Clientv6() {
	var (
		laddr, raddr net.UDPAddr
	)
	if *to == "" {
		llAddr, err := dhcpv6.GetLinkLocalAddr(*iface)
		if err != nil {
			log.Fatal(err)
		}
		laddr = net.UDPAddr{
			IP:   *llAddr,
			Port: 546,
			Zone: *iface,
		}
		raddr = net.UDPAddr{
			IP:   dhcpv6.AllDHCPRelayAgentsAndServers,
			Port: 547,
			Zone: *iface,
		}
	} else {
		laddr = net.UDPAddr{
			IP:   net.ParseIP("::"),
			Port: 546,
			Zone: *iface,
		}
		dstHost, dstPort, err := net.SplitHostPort(*to)
		if err != nil {
			log.Fatal(err)
		}
		dstPortInt, err := strconv.Atoi(dstPort)
		if err != nil {
			log.Fatal(err)
		}
		raddr = net.UDPAddr{
			IP:   net.ParseIP(dstHost),
			Port: dstPortInt,
			Zone: *iface, // this may clash with the scope passed in the dstHost, if any
		}
	}
	c := dhcpv6.Client{
		LocalAddr:  &laddr,
		RemoteAddr: &raddr,
	}
	conv, err := c.Exchange(*iface, nil)
	// don't exit immediately if there's an error, since `conv` will always
	// contain at least the SOLICIT message. So print it out first
	for _, m := range conv {
		fmt.Print(m.Summary())
	}
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	if *infile == "" && !*live {
		if *ver == 4 {
			Clientv4()
		} else {
			Clientv6()
		}
	} else {
		var (
			handle *pcap.Handle
			err    error
		)
		if *count < 0 {
			log.Fatal("count cannot be negative")
		}
		if *live {
			if *snaplen < 0 {
				log.Fatal("snaplen cannot be negative")
			}
			var slen int32
			slen = int32(*snaplen)
			if slen == 0 {
				// some libpcap versions don't support 0 as 'no snap len limit'.
				// Setting it to 262144 as per tcpdump's manual page
				slen = 262144
			}
			handle, err = pcap.OpenLive(*iface, slen, false /* promisc */, -1 /* timeout */)
		} else {
			handle, err = pcap.OpenOffline(*infile)
		}
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()
		var pcapFilter string
		if *ver == 6 {
			pcapFilter = "ip6 and udp portrange 546-547"
		} else {
			pcapFilter = "ip and udp portrange 67-68"
		}
		err = handle.SetBPFFilter(pcapFilter)
		if err != nil {
			log.Fatal(err)
		}
		var layerType gopacket.LayerType
		if *useEtherIP {
			layerType = layers.LayerTypeEtherIP
		} else {
			layerType = layers.LayerTypeEthernet
		}
		packetCount := 0
		for {
			packetCount++
			if *count != 0 && packetCount > *count {
				break
			}
			data, _, err := handle.ReadPacketData()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			pkt := gopacket.NewPacket(data, layerType, gopacket.Default)
			if *debug {
				fmt.Println(pkt)
			}
			if udpLayer := pkt.Layer(layers.LayerTypeUDP); udpLayer != nil {
				udp, _ := udpLayer.(*layers.UDP)
				if *debug {
					fmt.Println(udp.Payload)
				}
				if *ver == 4 {
					d, err := dhcpv4.FromBytes(udp.Payload)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(d.Summary())
				} else {
					var packet dhcpv6.DHCPv6
					d, err := dhcpv6.FromBytes(udp.Payload)
					if err != nil {
						log.Fatal(err)
					}
					packet = d
					if *unpack {
						if d.IsRelay() {
							inner, err := d.(*dhcpv6.DHCPv6Relay).GetInnerMessage()
							if err != nil {
								log.Fatal(err)
							}
							packet = inner
						}
					}
					fmt.Println(packet.Summary())
				}
			}
		}
	}
}
