# dhcp

DHCPv4 and DHCPv6 decoding/encoding library with client and server code, written in Go.

# How to get the library

The library is split into several parts:
* `dhcpv6`: implementation of DHCPv6 packet, client and server
* `dhcpv4`: implementation of DHCPv4 packet, client and server
* `netboot`: network booting wrappers on top of `dhcpv6` and `dhcpv4`
* `iana`: several IANA constants, and helpers used by `dhcpv6` and `dhcpv4`
* `rfc1035label`: simple implementation of RFC1035 labels, used by `dhcpv6` and
  `dhcpv4`

You will probably only need `dhcpv6` and/or `dhcpv4` explicitly. The rest is
pulled in automatically if necessary.


So, to get `dhcpv6` and `dhpv4` just run:
```
go get github.com/insomniacslk/dhcp/dhcpv{4,6}
```


# Examples

The sections below will illustrate how to use the `dhcpv6` and `dhcpv4`
packages.

See more example code at https://github.com/insomniacslk/exdhcp


## DHCPv6 client

To run a DHCPv6 transaction on the interface "eth0":

```
package main

import (
	"log"

	"github.com/insomniacslk/dhcp/dhcpv6"
)


func main() {
	// NewClient sets up a new DHCPv6 client with default values
	// for read and write timeouts, for destination address and listening
	// address
	client := dhcpv6.NewClient()
	
	// Exchange runs a Solicit-Advertise-Request-Reply transaction on the
	// specified network interface, and returns a list of DHCPv6 packets
	// (a "conversation") and an error if any. Notice that Exchange may
	// return a non-empty packet list even if there is an error. This is
	// intended, because the transaction may fail at any point, and we
	// still want to know what packets were exchanged until then.
	// The `nil` argument indicates that we want to use a default Solicit
	// packet, instead of specifying a custom one ourselves.
	conversation, err := client.Exchange("eth0", nil)
	
	// Summary() prints a verbose representation of the exchanged packets.
	for _, packet := range conversation {
		log.Print(packet.Summary())
	}
	// error handling is done *after* printing, so we still print the
	// exchanged packets if any, as explained above.
	if err != nil {
		log.Fatal(err)
	}
}
```


## DHCPv6 packet parsing

TODO


## DHCPv6 server

TODO

## DHCPv4 client

TODO


## DHCPv4 packet parsing

TODO


## DHCPv4 server

TODO


# Public projects that use it

* Facebook's DHCP load balancer, `dhcplb`, https://github.com/facebookincubator/dhcplb
* Systemboot, a LinuxBoot distribution that runs as system firmware, https://github.com/systemboot/systemboot 
* Router7, a pure-Go router implementation for fiber7 connections, https://github.com/rtr7/router7
* Beats from ElasticSearch, https://github.com/elastic/beats
* Bender from Pinterest, a library for load-testing, https://github.com/pinterest/bender
* FBender from Facebook, a tool for load-testing based on Bender, https://github.com/facebookincubator/fbender
