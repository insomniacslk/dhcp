//This is lease support for nclient4

package nclient4

import (
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

//Lease contains a DHCPv4 lease after DORA.
//note: Lease doesn't include binding interface name
type Lease struct {
	Offer        *dhcpv4.DHCPv4
	ACK          *dhcpv4.DHCPv4
	CreationTime time.Time
	IDOptions    dhcpv4.Options //DHCPv4 options to identify the client like client-id, option82/remote-id
}

// WithClientIDOptions configures a list of DHCPv4 option code that DHCP server
// uses to identify client, beside the MAC address.
func WithClientIDOptions(cidl dhcpv4.OptionCodeList) ClientOpt {
	return func(c *Client) (err error) {
		c.clientIDOptions = cidl
		return
	}
}

//Release send DHCPv4 release messsage to server, based on specified lease.
//release is sent as unicast per RFC2131, section 4.4.4.
//Note: some DHCP server requries of using assigned IP address as source IP,
//use nclient4.WithUnicast to create client for such case.
func (c *Client) Release(lease *Lease) error {
	if lease == nil {
		return fmt.Errorf("lease is nil")
	}
	modList := []dhcpv4.Modifier{
		dhcpv4.WithMessageType(dhcpv4.MessageTypeRelease),
		dhcpv4.WithClientIP(lease.ACK.YourIPAddr),
		dhcpv4.WithHwAddr(lease.ACK.ClientHWAddr),
		dhcpv4.WithBroadcast(false),
		dhcpv4.WithOption(dhcpv4.OptServerIdentifier(lease.ACK.ServerIdentifier())),
	}

	req, err := dhcpv4.New(modList...)
	if err != nil {
		return err
	}
	//This is to make sure use same client identification options used during
	//DORA, so that DHCP server could identify the required lease
	for t := range lease.IDOptions {
		req.UpdateOption(
			dhcpv4.OptGeneric(dhcpv4.GenericOptionCode(t),
				lease.IDOptions.Get(dhcpv4.GenericOptionCode(t))),
		)
	}
	_, err = c.conn.WriteTo(req.ToBytes(), &net.UDPAddr{IP: lease.ACK.Options.Get(dhcpv4.OptionServerIdentifier), Port: ServerPort})
	if err == nil {
		c.logger.PrintMessage("sent message:", req)
	}
	return err
}
