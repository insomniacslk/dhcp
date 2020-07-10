package nclient4_test

import (
	"context"
	"log"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
)

func Example() {
	ifname := "enp0s10"
	clntid := "client-1"
	var idoptlist dhcpv4.OptionCodeList
	//specify option61 is part of client identification used by DHCPv4 server
	idoptlist.Add(dhcpv4.OptionClientIdentifier)
	clntOptions := []nclient4.ClientOpt{nclient4.WithClientIDOptions(idoptlist), nclient4.WithDebugLogger()}
	clnt, err := nclient4.New(ifname, clntOptions...)
	if err != nil {
		log.Fatalf("failed to create dhcpv4 client,%v", err)
	}
	//add option61 to discovery and request
	option61 := dhcpv4.OptClientIdentifier([]byte(clntid))
	lease, err := clnt.Request(context.Background(), dhcpv4.WithOption(option61))
	if err != nil {
		log.Fatal(err)
	}
	//print the lease
	log.Printf("Got lease:\n%+v", lease)
	//release the lease
	log.Print("Releasing lease...")
	err = clnt.Release(lease)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("done")
}
