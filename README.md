# dhcp
[![Build Status](https://travis-ci.org/insomniacslk/dhcp.svg?branch=master)](https://travis-ci.org/insomniacslk/dhcp)
[![codecov](https://codecov.io/gh/insomniacslk/dhcp/branch/master/graph/badge.svg)](https://codecov.io/gh/insomniacslk/dhcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/insomniacslk/dhcp)](https://goreportcard.com/report/github.com/insomniacslk/dhcp)

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
go get -u github.com/insomniacslk/dhcp/dhcpv{4,6}
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
	// A default Solicit packet will be used during the "conversation",
	// which can be manipulated by using modifiers.
	conversation, err := client.Exchange("eth0")
	
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


## DHCPv6 packet crafting and manipulation

```
package main

import (
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/iana"
)

func main() {
	// In this example we create and manipulate a DHCPv6 solicit packet
	// and encapsulate it in a relay packet. To to this, we use
	// `dhcpv6.DHCPv6Message` and `dhcpv6.DHCPv6Relay`, two structures
	// that implement the `dhcpv6.DHCPv6` interface.
	// Then print the wire-format representation of the packet.

	// Create the DHCPv6 Solicit first, using the interface "eth0"
	// to get the MAC address
	msg, err := dhcpv6.NewSolicitForInterface("eth0")
	if err != nil {
		log.Fatal(err)
	}

	// In this example I want to redact the MAC address of my
	// network interface, so instead of replacing it manually,
	// I will show how to use modifiers for the purpose.
	// A Modifier is simply a function that can be applied on
	// a DHCPv6 object to manipulate it. Here we use it to
	// replace the MAC address with a dummy one.
	// Modifiers can be passed to many functions, for example
	// to constructors, `Exchange()`, `Solicit()`, etc. Check
	// the source code to know where to use them.
	// Existing modifiers are implemented in dhcpv6/modifiers.go .
	mac, err := net.ParseMAC("00:fa:ce:b0:0c:00")
	if err != nil {
		log.Fatal(err)
	}
	duid := dhcpv6.Duid{
		Type:          dhcpv6.DUID_LLT,
		HwType:        iana.HwTypeEthernet,
		Time:          dhcpv6.GetTime(),
		LinkLayerAddr: mac,
	}
	// As suggested above, an alternative is to call
	// dhcpv6.NewSolicitForInterface("eth0", dhcpv6.WithCLientID(duid))
	msg = dhcpv6.WithClientID(duid)(msg)

	// Now encapsulate the message in a DHCPv6 relay.
	// As per RFC3315, the link-address and peer-address have
	// to be set by the relay agent. We use dummy values here.
	linkAddr := net.ParseIP("2001:0db8::1")
	peerAddr := net.ParseIP("2001:0db8::2")
	relay, err := dhcpv6.EncapsulateRelay(msg, dhcpv6.MessageTypeRelayForward, linkAddr, peerAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Print a verbose representation of the relay packet, that will also
	// show a short representation of the inner Solicit message.
	// To print a detailed summary of the inner packet, extract it
	// first from the relay using `relay.GetInnerMessage()`.
	log.Print(relay.Summary())

	// And finally, print the bytes that would be sent on the wire
	log.Print(relay.ToBytes())

	// Note: there are many more functions in the library, check them
	// out in the source code. For example, if you want to decode a
	// byte stream into a DHCPv6 message or relay, you can use
	// `dhcpv6.FromBytes`.
}
```

The output (slightly modified for readability) is
```
$ go run main.go
2018/11/08 13:56:31 DHCPv6Relay
  messageType=RELAY-FORW
  hopcount=0
  linkaddr=2001:db8::1
  peeraddr=2001:db8::2
  options=[OptRelayMsg{relaymsg=DHCPv6Message(messageType=SOLICIT transactionID=0x9e0242, 4 options)}]

2018/11/08 13:56:31 [12 0 32 1 13 184 0 0 0 0 0 0 0 0 0 0 0 1 32 1 13 184
                     0 0 0 0 0 0 0 0 0 0 0 2 0 9 0 52 1 158 2 66 0 1 0 14
                     0 1 0 1 35 118 253 15 0 250 206 176 12 0 0 6 0 4 0 23
                     0 24 0 8 0 2 0 0 0 3 0 12 250 206 176 12 0 0 14 16 0
                     0 21 24]
```

## DHCPv6 server

A DHCPv6 server requires the user to implement a request handler. Basically the
user has to provide the logic to answer to each packet. The library offers a few
facilities to forge response packets, e.g. `NewAdvertiseFromSolicit`,
`NewReplyFromDHCPv6Message` and so on. Look at the source code to see what's
available.

An example server that will print (but not reply to) the client's request is
shown below:

```
package main

import (
        "log"
        "net"

        "github.com/insomniacslk/dhcp/dhcpv6"
)

func handler(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
        // this function will just print the received DHCPv6 message, without replying
        log.Print(m.Summary())
}

func main() {
        laddr := net.UDPAddr{
                IP:   net.ParseIP("::1"),
                Port: dhcpv6.DefaultServerPort,
        }
        server := dhcpv6.NewServer(laddr, handler)

        defer server.Close()
        if err := server.ActivateAndServe(); err != nil {
                log.Panic(err)
        }
}
```

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
* CoreDHCP, a fast, multithreaded, modular and extensible DHCP server, https://github.com/coredhcp/coredhcp
