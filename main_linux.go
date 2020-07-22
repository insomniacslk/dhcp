package main

import (
	_ "github.com/insomniacslk/dhcp/dhcpv4/async"
	_ "github.com/insomniacslk/dhcp/dhcpv4/bsdp"
	_ "github.com/insomniacslk/dhcp/dhcpv4/client4"
	_ "github.com/insomniacslk/dhcp/dhcpv4/nclient4"
	_ "github.com/insomniacslk/dhcp/netboot"
)

// Imports packages that can only be compiled on linux
