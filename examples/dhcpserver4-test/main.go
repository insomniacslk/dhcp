// Package main provides a simple DHCPv4 server example for testing.
package main

import (
	"flag"
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

var (
	ifname     = flag.String("i", "Ethernet", "Interface name to listen on")
	serverIP   = flag.String("s", "192.168.1.1", "Server IP address")
	rangeStart = flag.String("start", "192.168.1.100", "Start of IP range")
	rangeEnd   = flag.String("end", "192.168.1.200", "End of IP range")
	mask       = flag.String("mask", "255.255.255.0", "Subnet mask")
	router     = flag.String("router", "192.168.1.1", "Default router")
	dns        = flag.String("dns", "8.8.8.8", "DNS server")
	leaseTime  = flag.Int("lease", 3600, "Lease time in seconds")
)

// Simple IP allocator (not production-ready, just for testing)
var nextIP net.IP

func handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	log.Printf("Received %s from %s", m.MessageType(), peer)
	log.Printf("  Client MAC: %s", m.ClientHWAddr)
	log.Printf("  Transaction ID: %v", m.TransactionID)

	var resp *dhcpv4.DHCPv4
	var err error

	switch m.MessageType() {
	case dhcpv4.MessageTypeDiscover:
		// Allocate an IP (simple incrementing, not production-ready)
		if nextIP == nil {
			nextIP = net.ParseIP(*rangeStart).To4()
		}
		allocatedIP := make(net.IP, 4)
		copy(allocatedIP, nextIP)

		resp, err = dhcpv4.NewReplyFromRequest(m,
			dhcpv4.WithMessageType(dhcpv4.MessageTypeOffer),
			dhcpv4.WithServerIP(net.ParseIP(*serverIP)),
			dhcpv4.WithYourIP(allocatedIP),
			dhcpv4.WithNetmask(net.IPMask(net.ParseIP(*mask).To4())),
			dhcpv4.WithRouter(net.ParseIP(*router)),
			dhcpv4.WithDNS(net.ParseIP(*dns)),
			dhcpv4.WithLeaseTime(uint32(*leaseTime)),
		)
		if err != nil {
			log.Printf("Error creating OFFER: %v", err)
			return
		}
		log.Printf("Sending OFFER with IP %s", allocatedIP)

	case dhcpv4.MessageTypeRequest:
		// For simplicity, just ACK whatever was requested
		requestedIP := m.RequestedIPAddress()
		if requestedIP == nil {
			requestedIP = m.YourIPAddr
		}

		resp, err = dhcpv4.NewReplyFromRequest(m,
			dhcpv4.WithMessageType(dhcpv4.MessageTypeAck),
			dhcpv4.WithServerIP(net.ParseIP(*serverIP)),
			dhcpv4.WithYourIP(requestedIP),
			dhcpv4.WithNetmask(net.IPMask(net.ParseIP(*mask).To4())),
			dhcpv4.WithRouter(net.ParseIP(*router)),
			dhcpv4.WithDNS(net.ParseIP(*dns)),
			dhcpv4.WithLeaseTime(uint32(*leaseTime)),
		)
		if err != nil {
			log.Printf("Error creating ACK: %v", err)
			return
		}
		log.Printf("Sending ACK for IP %s", requestedIP)

		// Increment next IP for next client
		nextIP[3]++
		if nextIP[3] > net.ParseIP(*rangeEnd).To4()[3] {
			nextIP = net.ParseIP(*rangeStart).To4()
		}

	default:
		log.Printf("Ignoring message type %s", m.MessageType())
		return
	}

	if resp != nil {
		log.Printf("Response: %s", resp.Summary())
		if _, err := conn.WriteTo(resp.ToBytes(), peer); err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
}

func main() {
	flag.Parse()

	log.Printf("Starting DHCPv4 server on interface %s", *ifname)
	log.Printf("Server IP: %s", *serverIP)
	log.Printf("IP Range: %s - %s", *rangeStart, *rangeEnd)
	log.Printf("Subnet Mask: %s", *mask)
	log.Printf("Router: %s", *router)
	log.Printf("DNS: %s", *dns)
	log.Printf("Lease Time: %d seconds", *leaseTime)

	laddr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 67,
	}

	server, err := server4.NewServer(*ifname, &laddr, handler,
		server4.WithSummaryLogger())
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Printf("Server listening on %s:67", *ifname)
	if err := server.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
