package main

import (
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/server6"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/mdlayher/netx/eui64"
)

func handler(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
	// Print message.
	log.Print(m.Summary())

	msg, wrapRelays, err := dhcpv6.ServerPeelRelays(m)
	if err != nil {
		log.Printf("Invalid DHCPv6 message: %v", err)
		return
	}

	mac, err := dhcpv6.ExtractMAC(msg)
	if err != nil {
		log.Printf("Need MAC in request to service: %v", err)
		return
	}

	assignedIP, err := eui64.ParseMAC(net.ParseIP("fd00::"), mac)
	if err != nil {
		log.Printf("Could not generate ULA IP from MAC: %v", err)
		return
	}

	clientID := msg.Options.ClientID()
	if clientID == nil {
		log.Printf("No Client ID in request. Can only accept requests with Client ID.")
		return
	}

	var resp *dhcpv6.Message

	switch msg.MessageType {
	case dhcpv6.MessageTypeSolicit:
		if msg.GetOneOption(dhcpv6.OptionRapidCommit) == nil {
			resp, err = dhcpv6.NewReplyFromMessage(msg)
		} else {
			resp, err = dhcpv6.NewAdvertiseFromSolicit(msg)
		}
		if err != nil {
			log.Printf("Failed to create reply for %v: %v", m, err)
			return
		}

		// Was an IP assignment requested?
		if ianaRequest := msg.Options.OneIANA(); ianaRequest != nil {
			ianaRequest.Options.Add(&dhcpv6.OptIAAddress{IPv6Addr: assignedIP})
			resp.AddOption(ianaRequest)
		}

	case dhcpv6.MessageTypeRequest:
		resp, err = dhcpv6.NewReplyFromMessage(msg)
		if err != nil {
			log.Printf("Failed to create reply for %v: %v", m, err)
			return
		}

		for _, ianaReq := range msg.Options.IANA() {
			ips := ianaReq.Options.Addresses()
			for _, ip := range ips {
				if !ip.IPv6Addr.Equal(assignedIP) {
					ianaReq.Options.Add(&dhcpv6.OptStatusCode{
						StatusCode:    iana.StatusNotOnLink,
						StatusMessage: "IP address not supported",
					})
				}
			}
			resp.AddOption(ianaReq)
		}

	default:
		log.Printf("No handling implemented for message type for %v", msg)
	}

	// Add the RelayReply onion layers corresponding to the initially
	// peeled RelayForw messages.
	reply := wrapRelays(resp)

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		log.Printf("Failed to send response %v: %v", reply, err)
		return
	}
}

func main() {
	laddr := &net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: dhcpv6.DefaultServerPort,
	}
	server, err := server6.NewServer("", laddr, handler)
	if err != nil {
		log.Fatal(err)
	}

	server.Serve()
}
