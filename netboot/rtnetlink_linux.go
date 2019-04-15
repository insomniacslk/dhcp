package netboot

import (
	"encoding/binary"
	"net"

	"github.com/jsimonetti/rtnetlink"
	"golang.org/x/sys/unix"
)

// RTNL is a rtnetlink object with a high-level interface.
type RTNL struct {
	conn *rtnetlink.Conn
}

func (r *RTNL) init() error {
	if r.conn != nil {
		return nil
	}
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return err
	}
	r.conn = conn
	return nil
}

// Close closes the netlink connection. Must be called to avoid leaks!
func (r *RTNL) Close() {
	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}
}

// GetLinkState returns the operational state for the given interface index.
func (r *RTNL) GetLinkState(iface int) (rtnetlink.OperationalState, error) {
	if err := r.init(); err != nil {
		return 0, err
	}
	msg, err := r.conn.Link.Get(uint32(iface))
	if err != nil {
		return 0, err
	}
	return msg.Attributes.OperationalState, nil
}

// SetLinkState sets the operational state up or down for the given interface
// index.
func (r *RTNL) SetLinkState(iface int, up bool) error {
	if err := r.init(); err != nil {
		return err
	}
	var state uint32
	if up {
		state = unix.IFF_UP
	}
	msg := rtnetlink.LinkMessage{
		Family: unix.AF_UNSPEC,
		Type:   unix.ARPHRD_NETROM,
		Index:  uint32(iface),
		Flags:  state,
		Change: unix.IFF_UP,
	}
	if err := r.conn.Link.Set(&msg); err != nil {
		return err
	}
	return nil
}

func getFamily(ip net.IP) int {
	if ip.To4() != nil {
		return unix.AF_INET
	}
	return unix.AF_INET6
}

// SetAddr sets the interface address.
func (r *RTNL) SetAddr(iface int, a net.IPNet) error {
	if err := r.init(); err != nil {
		return err
	}
	ones, _ := a.Mask.Size()
	msg := rtnetlink.AddressMessage{
		Family:       uint8(getFamily(a.IP)),
		PrefixLength: uint8(ones),
		// TODO detect the right scope to set, or get it as input argument
		Scope: unix.RT_SCOPE_UNIVERSE,
		Index: uint32(iface),
		Attributes: rtnetlink.AddressAttributes{
			Address: a.IP,
			Local:   a.IP,
		},
	}
	if a.IP.To4() != nil {
		// Broadcast is only required for IPv4
		ip := make(net.IP, net.IPv4len)
		binary.BigEndian.PutUint32(
			ip,
			binary.BigEndian.Uint32(a.IP.To4())|
				^binary.BigEndian.Uint32(net.IP(a.Mask).To4()))
		msg.Attributes.Broadcast = ip
	}
	if err := r.conn.Address.New(&msg); err != nil {
		return err
	}
	return nil
}

// RouteDel deletes a route to the given destination
func (r *RTNL) RouteDel(dst net.IP) error {
	if err := r.init(); err != nil {
		return err
	}
	msg := rtnetlink.RouteMessage{
		Family: uint8(getFamily(dst)),
		Table:  unix.RT_TABLE_MAIN,
		// TODO make this configurable?
		Protocol: unix.RTPROT_UNSPEC,
		// TODO make this configurable?
		Scope: unix.RT_SCOPE_NOWHERE,
		Type:  unix.RTN_UNSPEC,
		Attributes: rtnetlink.RouteAttributes{
			Dst: dst,
		},
	}
	if err := r.conn.Route.Delete(&msg); err != nil {
		return err
	}
	return nil
}

// RouteAdd adds a route to dst, from src (if set), via gw.
func (r *RTNL) RouteAdd(iface int, dst, src net.IPNet, gw net.IP) error {
	if err := r.init(); err != nil {
		return err
	}
	dstLen, _ := dst.Mask.Size()
	srcLen, _ := src.Mask.Size()
	msg := rtnetlink.RouteMessage{
		Family: uint8(getFamily(dst.IP)),
		Table:  unix.RT_TABLE_MAIN,
		// TODO make this configurable?
		Protocol: unix.RTPROT_BOOT,
		// TODO make this configurable?
		Scope:     unix.RT_SCOPE_UNIVERSE,
		Type:      unix.RTN_UNICAST,
		DstLength: uint8(dstLen),
		SrcLength: uint8(srcLen),
		Attributes: rtnetlink.RouteAttributes{
			Dst:      dst.IP,
			Src:      src.IP,
			Gateway:  gw,
			OutIface: uint32(iface),
		},
	}
	if err := r.conn.Route.Add(&msg); err != nil {
		return err
	}
	return nil

}
