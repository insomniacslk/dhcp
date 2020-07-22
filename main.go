package main

import (
	_ "github.com/insomniacslk/dhcp/dhcpv4"
	_ "github.com/insomniacslk/dhcp/dhcpv4/server4"
	_ "github.com/insomniacslk/dhcp/dhcpv4/ztpv4"
	_ "github.com/insomniacslk/dhcp/dhcpv6"
	_ "github.com/insomniacslk/dhcp/dhcpv6/async"
	_ "github.com/insomniacslk/dhcp/dhcpv6/client6"
	_ "github.com/insomniacslk/dhcp/dhcpv6/nclient6"
	_ "github.com/insomniacslk/dhcp/dhcpv6/server6"
	_ "github.com/insomniacslk/dhcp/dhcpv6/ztpv6"
	_ "github.com/insomniacslk/dhcp/iana"
	_ "github.com/insomniacslk/dhcp/interfaces"
	_ "github.com/insomniacslk/dhcp/rfc1035label"
)

// The only purpose is to import all packages so that "go build" or "go test ./..."
// triggered updating go.mod/go.sum
func main() {
}
