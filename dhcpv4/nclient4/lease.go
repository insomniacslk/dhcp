//This is lease managment for nclient4

package nclient4

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/vishvananda/netlink"
)

const (
	//default lease time if server doesn't return lease time option or return zero
	defaultLeaseTime = time.Hour
)

//DHCPv4ClientLease contains a DHCPv4 lease after DORA,
//could be used for creating a new Client with NewWithLease()
type DHCPv4ClientLease struct {
	IfName         string
	MACAddr        net.HardwareAddr
	ServerAddr     net.UDPAddr
	AssignedIP     net.IP
	AssignedIPMask net.IPMask
	CreationTime   time.Time
	LeaseDuration  time.Duration
	RenewInterval  time.Duration
	RebindInterval time.Duration
	IDOptions      dhcpv4.Options //DHCPv4 options to identify the client like client-id, option82/remote-id
	AckOptions     dhcpv4.Options //DHCPv4 options in ACK, could be used for applying lease

}

//return a string representation
func (lease DHCPv4ClientLease) String() string {
	const fmtstr = "%-35s\t%-35s\n"
	const timefmtstr = "01/02/2006 15:04:05.000"
	rstr := fmt.Sprintf(fmtstr, fmt.Sprintf("Interface:%v", lease.IfName), fmt.Sprintf("MAC:%v", lease.MACAddr))
	rstr += fmt.Sprintf(fmtstr, fmt.Sprintf("Svr:%v", lease.ServerAddr.IP), fmt.Sprintf("Created:%v", lease.CreationTime.Format(timefmtstr)))
	prefixlen, _ := lease.AssignedIPMask.Size()
	rstr += fmt.Sprintf(fmtstr, fmt.Sprintf("IP:%v/%v", lease.AssignedIP, prefixlen), fmt.Sprintf("Lease time:%v", lease.LeaseDuration))
	rstr += fmt.Sprintf(fmtstr, fmt.Sprintf("Renew interval:%v", lease.RenewInterval), fmt.Sprintf("Rebind interval:%v", lease.RebindInterval))
	rstr += fmt.Sprintf("Id options:\n%v", lease.IDOptions)
	rstr += fmt.Sprintf("ACK options:\n%v", lease.AckOptions)
	return rstr
}

// WithClientIDOptions configures a list of DHCPv4 option code that DHCP server
// uses to identify client, beside the MAC address.
func WithClientIDOptions(cidl dhcpv4.OptionCodeList) ClientOpt {
	return func(c *Client) (err error) {
		c.clientIDOptions = cidl
		return
	}
}

// WithApplyLeaseHandler specifies a handler function which is called when
// Client.ApplyLease() is called; without this, a default handler function is called.
// the default handler will add/remove the assigned address to/from the binding interface;
// bool parameter is true when lease is applied, false when lease is released
func WithApplyLeaseHandler(h func(DHCPv4ClientLease, bool) error) ClientOpt {
	return func(c *Client) (err error) {
		c.leaseApplyHandler = h
		return
	}
}

//default lease apply handler
//add/remove address to/from binding interface
func defaultLeaseApplyHandler(l DHCPv4ClientLease, enable bool) error {
	link, err := netlink.LinkByName(l.IfName)
	if err != nil {
		return err
	}
	plen, _ := l.AssignedIPMask.Size()
	prefixstr := fmt.Sprintf("%v/%v", l.AssignedIP, plen)
	naddr, err := netlink.ParseAddr(prefixstr)
	if err != nil {
		return err
	}
	if enable {
		err = netlink.AddrReplace(link, naddr)

	} else {
		err = netlink.AddrDel(link, naddr)
	}
	return err

}

//ApplyLease apply/unapply the lease, call the c.leaseApplyHandler
func (c *Client) ApplyLease(enable bool) error {
	if c.lease == nil {
		return fmt.Errorf("no lease to apply")
	}
	return c.leaseApplyHandler(c.GetLease(), enable)
}

//GetLease return the lease
func (c *Client) GetLease() (clease DHCPv4ClientLease) {
	clease = *c.lease
	clease.MACAddr = c.ifaceHWAddr
	clease.IfName = c.ifName
	clease.ServerAddr = *c.serverAddr
	return
}

// RequestSavingLease completes DORA handshake and store&apply the lease
//
// Note that modifiers will be applied *both* to Discover and Request packets.
func (c *Client) RequestSavingLease(ctx context.Context, modifiers ...dhcpv4.Modifier) (offer, ack *dhcpv4.DHCPv4, err error) {
	offer, err = c.DiscoverOffer(ctx, modifiers...)
	if err != nil {
		err = fmt.Errorf("unable to receive an offer: %w", err)
		return
	}

	// TODO(chrisko): should this be unicast to the server?
	request, err := dhcpv4.NewRequestFromOffer(offer, dhcpv4.PrependModifiers(modifiers,
		dhcpv4.WithOption(dhcpv4.OptMaxMessageSize(MaxMessageSize)))...)
	if err != nil {
		err = fmt.Errorf("unable to create a request: %w", err)
		return
	}

	ack, err = c.SendAndRead(ctx, c.serverAddr, request, nil)
	if err != nil {
		err = fmt.Errorf("got an error while processing the request: %w", err)
		return
	}
	//save lease
	c.lease = &DHCPv4ClientLease{}
	c.lease.AssignedIP = ack.YourIPAddr
	c.lease.AssignedIPMask = ack.SubnetMask()
	c.lease.CreationTime = time.Now()
	c.lease.LeaseDuration = ack.IPAddressLeaseTime(0)
	if c.lease.LeaseDuration == 0 {
		c.lease.LeaseDuration = defaultLeaseTime
		c.logger.Printf("warning: server doesn't include Lease Time option or it is zero seconds, setting lease time to default %v", c.lease.LeaseDuration)

	}
	c.lease.RenewInterval = ack.IPAddressRenewalTime(0)
	if c.lease.RenewInterval == 0 {
		//setting default to half of lease time based on RFC2131,section 4.4.5
		c.lease.RenewInterval = time.Duration(float64(c.lease.LeaseDuration) / 2)
		c.logger.Printf("warning: server doesn't include Renew Time option or it is zero seconds, setting lease time to default %v", c.lease.RenewInterval)

	}
	c.lease.RebindInterval = ack.IPAddressRebindingTime(0)
	if c.lease.RebindInterval == 0 {
		//setting default to 0.875 of lease time based on RFC2131,section 4.4.5
		c.lease.RebindInterval = time.Duration(float64(c.lease.LeaseDuration) * 0.875)
		c.logger.Printf("warning: server doesn't include Renew Time option or it is zero seconds, setting lease time to default %v", c.lease.RebindInterval)

	}
	c.lease.IDOptions = dhcpv4.Options{}
	for _, optioncode := range c.clientIDOptions {
		v := request.Options.Get(optioncode)
		c.lease.IDOptions.Update(dhcpv4.OptGeneric(optioncode, v))
	}
	c.lease.AckOptions = ack.Options
	//update server address
	c.serverAddr = &(net.UDPAddr{IP: ack.ServerIdentifier(), Port: 67})
	err = c.ApplyLease(true)
	return
}

//Release send DHCPv4 release messsage to server.
//release is sent as unicast per RFC2131, section 4.4.4.
//The lease need to be applied with c.ApplyLease(true) first before calling Release.
func (c *Client) Release() error {
	if c.lease == nil {
		return fmt.Errorf("There is no lease to release")
	}
	req, err := dhcpv4.New()
	if err != nil {
		return err
	}
	//This is to make sure use same client identification options used during
	//DORA, so that DHCP server could identify the required lease
	req.Options = c.lease.IDOptions

	req.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeRelease))
	req.ClientHWAddr = c.ifaceHWAddr
	req.ClientIPAddr = c.lease.AssignedIP
	req.UpdateOption(dhcpv4.OptServerIdentifier(c.serverAddr.IP))
	req.SetUnicast()
	luaddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", c.lease.AssignedIP, 68))
	if err != nil {
		return err
	}

	uniconn, err := net.DialUDP("udp4", luaddr, c.serverAddr)
	if err != nil {
		return err
	}
	_, err = uniconn.Write(req.ToBytes())
	if err != nil {
		return err
	}
	c.logger.PrintMessage("sent message:", req)
	return c.ApplyLease(false)
}

//NewWithLease return a Client with populated lease.
//this function could be used to release a saved lease.
func NewWithLease(clease DHCPv4ClientLease, opts ...ClientOpt) (*Client, error) {
	clntoptlist := []ClientOpt{
		WithServerAddr(&clease.ServerAddr),
		WithHWAddr(clease.MACAddr),
	}
	clntoptlist = append(clntoptlist, opts...)
	clnt, err := New(clease.IfName, clntoptlist...)
	if err != nil {
		return nil, err
	}
	clnt.ifName = clease.IfName
	clnt.lease = &clease
	for optioncode := range clease.IDOptions {
		clnt.clientIDOptions.Add(dhcpv4.GenericOptionCode(optioncode))
	}
	return clnt, nil

}
