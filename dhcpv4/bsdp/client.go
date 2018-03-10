// +build darwin

package bsdp

import (
	"errors"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Exchange runs a full BSDP exchange (Inform[list], Ack, Inform[select],
// Ack). Returns a list of DHCPv4 structures representing the exchange.
func Exchange(client *dhcpv4.Client, ifname string, informList *dhcpv4.DHCPv4) ([]dhcpv4.DHCPv4, error) {
	conversation := make([]dhcpv4.DHCPv4, 1)
	var err error

	// Get our file descriptor for the broadcast socket.
	fd, err := dhcpv4.MakeBroadcastSocket(ifname)
	if err != nil {
		return conversation, err
	}

	// INFORM[LIST]
	if informList == nil {
		informList, err = NewInformListForInterface(ifname, dhcpv4.ClientPort)
		if err != nil {
			return conversation, err
		}
	}
	conversation[0] = *informList

	// ACK[LIST]
	ackForList, err := dhcpv4.SendReceive(client, fd, informList)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, *ackForList)

	// Parse boot images sent back by server
	bootImages, err := ParseBootImageListFromAck(*ackForList)
	if err != nil {
		return conversation, err
	}
	if len(bootImages) == 0 {
		return conversation, errors.New("got no BootImages from server")
	}

	// INFORM[SELECT]
	informSelect, err := InformSelectForAck(*ackForList, dhcpv4.ClientPort, bootImages[0])
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, *informSelect)

	// ACK[SELECT]
	ackForSelect, err := dhcpv4.SendReceive(client, fd, informSelect)
	if err != nil {
		return conversation, err
	}
	return append(conversation, *ackForSelect), nil
}
