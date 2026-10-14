// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.12 && windows

package nclient4

import (
	"errors"
	"fmt"
	"net"
	"time"
)

var (
	// BroadcastMac is the broadcast MAC address.
	//
	// Any UDP packet sent to this address is broadcast on the subnet.
	BroadcastMac = net.HardwareAddr([]byte{255, 255, 255, 255, 255, 255})
)

var (
	// ErrUDPAddrIsRequired is an error used when a passed argument is not of type "*net.UDPAddr".
	ErrUDPAddrIsRequired = errors.New("must supply UDPAddr")
)

// WindowsUDPConn wraps a standard UDP connection for Windows
type WindowsUDPConn struct {
	conn      *net.UDPConn
	boundAddr *net.UDPAddr
}

// NewRawUDPConn returns a UDP connection bound to the port.
// On Windows, we cannot bind to a specific interface, so we listen on all interfaces.
//
// The interface parameter is ignored on Windows.
func NewRawUDPConn(iface string, port int) (net.PacketConn, error) {
	// Verify interface exists (for error reporting)
	_, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, fmt.Errorf("interface %s not found: %v", iface, err)
	}

	// Use standard UDP socket on Windows - listen on all interfaces
	addr := &net.UDPAddr{IP: net.IPv4zero, Port: port}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %v", port, err)
	}

	return &WindowsUDPConn{
		conn:      conn,
		boundAddr: addr,
	}, nil
}

// NewBroadcastUDPConn returns a PacketConn that can send and receive UDP packets.
// On Windows, this wraps the provided PacketConn.
func NewBroadcastUDPConn(rawPacketConn net.PacketConn, boundAddr *net.UDPAddr) net.PacketConn {
	return &BroadcastRawUDPConn{
		PacketConn: rawPacketConn,
		boundAddr:  boundAddr,
	}
}

// BroadcastRawUDPConn wraps a PacketConn for Windows compatibility
type BroadcastRawUDPConn struct {
	net.PacketConn
	boundAddr *net.UDPAddr
}

// ReadFrom implements net.PacketConn.ReadFrom
func (w *WindowsUDPConn) ReadFrom(b []byte) (int, net.Addr, error) {
	n, addr, err := w.conn.ReadFromUDP(b)
	if err != nil {
		return 0, nil, err
	}
	return n, addr, nil
}

// WriteTo implements net.PacketConn.WriteTo
func (w *WindowsUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		return 0, ErrUDPAddrIsRequired
	}
	return w.conn.WriteTo(b, udpAddr)
}

// Close implements net.PacketConn.Close
func (w *WindowsUDPConn) Close() error {
	return w.conn.Close()
}

// LocalAddr implements net.PacketConn.LocalAddr
func (w *WindowsUDPConn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

// SetDeadline implements net.PacketConn.SetDeadline
func (w *WindowsUDPConn) SetDeadline(t time.Time) error {
	return w.conn.SetDeadline(t)
}

// SetReadDeadline implements net.PacketConn.SetReadDeadline
func (w *WindowsUDPConn) SetReadDeadline(t time.Time) error {
	return w.conn.SetReadDeadline(t)
}

// SetWriteDeadline implements net.PacketConn.SetWriteDeadline
func (w *WindowsUDPConn) SetWriteDeadline(t time.Time) error {
	return w.conn.SetWriteDeadline(t)
}
