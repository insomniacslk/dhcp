package netboot

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv6"
)

// RequestNetbootv6 sends a netboot request via DHCPv6 and returns the exchanged packets. Additional modifiers
// can be passed to manipulate both solicit and advertise packets.
func RequestNetbootv6(ifname string, timeout time.Duration, retries int, modifiers ...dhcpv6.Modifier) ([]dhcpv6.DHCPv6, error) {
	var (
		conversation []dhcpv6.DHCPv6
	)
	modifiers = append(modifiers, dhcpv6.WithNetboot)
	delay := 2 * time.Second
	for i := 0; i <= retries; i++ {
		log.Printf("sending request, attempt #%d", i+1)
		solicit, err := dhcpv6.NewSolicitForInterface(ifname, modifiers...)
		if err != nil {
			return nil, fmt.Errorf("failed to create SOLICIT for interface %s: %v", ifname, err)
		}

		client := dhcpv6.NewClient()
		client.ReadTimeout = timeout
		conversation, err = client.Exchange(ifname, solicit, modifiers...)
		if err != nil {
			log.Printf("Client.Exchange failed: %v", err)
			log.Printf("sleeping %v before retrying", delay)
			if i >= retries {
				// don't wait at the end of the last attempt
				break
			}
			time.Sleep(delay)
			// TODO add random splay
			delay = delay * 2
			continue
		}
		break
	}
	return conversation, nil
}

// ConversationToNetconf extracts network configuration and boot file URL from a
// DHCPv6 4-way conversation and returns them, or an error if any.
func ConversationToNetconf(conversation []dhcpv6.DHCPv6) (*NetConf, string, error) {
	var reply dhcpv6.DHCPv6
	for _, m := range conversation {
		// look for a REPLY
		if m.Type() == dhcpv6.REPLY {
			reply = m
			break
		}
	}
	if reply == nil {
		return nil, "", errors.New("no REPLY received")
	}
	netconf, err := GetNetConfFromPacketv6(reply.(*dhcpv6.DHCPv6Message))
	if err != nil {
		return nil, "", fmt.Errorf("cannot get netconf from packet: %v", err)
	}
	// look for boot file
	var (
		opt      dhcpv6.Option
		bootfile string
	)
	opt = reply.GetOneOption(dhcpv6.OPT_BOOTFILE_URL)
	if opt == nil {
		log.Printf("no bootfile URL option found in REPLY, looking for it in ADVERTISE")
		// as a fallback, look for bootfile URL in the advertise
		var advertise dhcpv6.DHCPv6
		for _, m := range conversation {
			// look for an ADVERTISE
			if m.Type() == dhcpv6.ADVERTISE {
				advertise = m
				break
			}
		}
		if advertise == nil {
			return nil, "", errors.New("no ADVERTISE found")
		}
		opt = advertise.GetOneOption(dhcpv6.OPT_BOOTFILE_URL)
		if opt == nil {
			return nil, "", errors.New("no bootfile URL option found in ADVERTISE")
		}
	}
	obf := opt.(*dhcpv6.OptBootFileURL)
	bootfile = string(obf.BootFileURL)
	return netconf, bootfile, nil
}
