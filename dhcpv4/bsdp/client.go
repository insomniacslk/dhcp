package bsdp

import (
	"errors"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Client represents a BSDP client that can perform BSDP exchanges via the
// broadcast address.
type Client dhcpv4.Client

// NewClient constructs a new client with default read and write timeouts from
// dhcpv4.Client.
func NewClient() *Client {
	c := dhcpv4.NewClient()
	return &Client{
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	}
}

func castVendorOpt(ack *dhcpv4.DHCPv4) {
	opts := ack.Options()
	for i := 0; i < len(opts); i++ {
		if opts[i].Code() == dhcpv4.OptionVendorSpecificInformation {
			vendorOpt, err := ParseOptVendorSpecificInformation(opts[i].ToBytes())
			// Oh well, we tried
			if err != nil {
				return
			}
			opts[i] = vendorOpt
		}
	}
}

// Exchange runs a full BSDP exchange (Inform[list], Ack, Inform[select],
// Ack). Returns a list of DHCPv4 structures representing the exchange.
func (c *Client) Exchange(ifname string, informList *dhcpv4.DHCPv4) ([]dhcpv4.DHCPv4, error) {
	conversation := make([]dhcpv4.DHCPv4, 1)
	var err error

	// Get our file descriptor for the broadcast socket.
	sendFd, err := dhcpv4.MakeBroadcastSocket(ifname)
	if err != nil {
		return conversation, err
	}
	recvFd, err := dhcpv4.MakeListeningSocket(ifname)
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
	ackForList, err := dhcpv4.BroadcastSendReceive(sendFd, recvFd, informList, c.ReadTimeout, c.WriteTimeout, dhcpv4.MessageTypeAck)
	if err != nil {
		return conversation, err
	}

	// Rewrite vendor-specific option for pretty printing.
	castVendorOpt(ackForList)
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
	ackForSelect, err := dhcpv4.BroadcastSendReceive(sendFd, recvFd, informSelect, c.ReadTimeout, c.WriteTimeout, dhcpv4.MessageTypeAck)
	castVendorOpt(ackForSelect)
	if err != nil {
		return conversation, err
	}
	return append(conversation, *ackForSelect), nil
}
