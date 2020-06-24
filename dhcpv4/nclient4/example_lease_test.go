//this is an example for nclient4 with lease/release

package nclient4_test

import (
	"context"
	"fmt"
	"log"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
	"github.com/vishvananda/netlink"
)

//applyLease adding the assigned ip to the interface specified by ifname
func applyLease(lease *nclient4.Lease, ifname string) error {
	link, err := netlink.LinkByName(ifname)
	if err != nil {
		return err
	}
	prefixlen := 32
	if ipmask := lease.ACK.SubnetMask(); ipmask != nil {
		prefixlen, _ = ipmask.Size()

	}
	prefixstr := fmt.Sprintf("%v/%v", lease.ACK.YourIPAddr, prefixlen)
	naddr, err := netlink.ParseAddr(prefixstr)
	if err != nil {
		return err
	}
	err = netlink.AddrReplace(link, naddr)
	return err

}

func main() {
	ifname := "eth1.200"
	remoteid := "client-1"
	var idoptlist dhcpv4.OptionCodeList
	//specify option82 is part of client identification used by DHCPv4 server
	idoptlist.Add(dhcpv4.OptionRelayAgentInformation)
	clntOptions := []nclient4.ClientOpt{nclient4.WithClientIDOptions(idoptlist), nclient4.WithDebugLogger()}
	clnt, err := nclient4.New(ifname, clntOptions...)
	if err != nil {
		log.Fatalf("failed to create dhcpv4 client,%v", err)
	}
	//adding option82/remote-id option to discovery and request
	remoteidsubopt := dhcpv4.OptGeneric(dhcpv4.AgentRemoteIDSubOption, []byte(remoteid))
	option82 := dhcpv4.OptRelayAgentInfo(remoteidsubopt)
	_, lease, err := clnt.Request(context.Background(), dhcpv4.WithOption(option82))
	if err != nil {
		log.Fatal(err)
	}
	//print the lease
	log.Printf("Got lease:\n%+v", lease)
	//apply the lease
	applyLease(lease, ifname)
	//release the lease
	log.Print("Releasing lease...")
	err = clnt.Release(lease)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("done")
}
