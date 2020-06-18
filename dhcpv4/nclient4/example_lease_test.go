//this is an example for nclient4 with lease

package nclient4_test

import (
	"context"
	"log"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
)

func Example_dHCPv4ClientLease() {
	ifname := "eth0"
	remoteID := "client-1"
	var idoptlist dhcpv4.OptionCodeList
	//specify option82 is part of client identification used by DHCPv4 server
	idoptlist.Add(dhcpv4.OptionRelayAgentInformation)
	clntOptions := []nclient4.ClientOpt{nclient4.WithClientIDOptions(idoptlist), nclient4.WithDebugLogger()}
	clnt, err := nclient4.New(ifname, clntOptions...)
	if err != nil {
		log.Fatalf("failed to create dhcpv4 client,%v", err)
	}
	//adding option82/remote-id option to discovery and request
	remoteIDSubOpt := dhcpv4.OptGeneric(dhcpv4.AgentRemoteIDSubOption, []byte(remoteID))
	option82 := dhcpv4.OptRelayAgentInfo(remoteIDSubOpt)
	_, _, err = clnt.RequestSavingLease(context.Background(), dhcpv4.WithOption(option82))
	if err != nil {
		log.Fatal(err)
	}
	//print the lease
	log.Printf("Got lease:\n%v", clnt.GetLease())
	//release the lease
	log.Print("Releasing lease...")
	err = clnt.Release()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("done")
}
