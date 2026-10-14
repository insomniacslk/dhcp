// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build integration && linux
// +build integration,linux

package nclient4

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/mdlayher/packet"
	"golang.org/x/sys/unix"
)

// TestClientFilterOnRealSocket brings up a veth pair and checks the kernel drops
// non-DHCP frames while a real reply still arrives, covering what the VM test
// cannot: the SO_ATTACH_FILTER call and the offset base AF_PACKET presents. It
// goes through NewRawUDPConn and reads the raw socket underneath, so a filter
// that never attached cannot hide behind ReadFrom's own checks. Needs root; the
// CI runs it under sudo.
func TestClientFilterOnRealSocket(t *testing.T) {
	// PID-suffixed so an interrupted run leaves no name the next one collides
	// with. IFNAMSIZ caps these at 15 characters, which a 7-digit pid fits.
	cli := fmt.Sprintf("dhcpcli%d", os.Getpid())
	srv := fmt.Sprintf("dhcpsrv%d", os.Getpid())
	ipLink(t, "add", cli, "type", "veth", "peer", "name", srv)
	t.Cleanup(func() { exec.Command("ip", "link", "del", cli).Run() })
	ipLink(t, "set", cli, "up")
	ipLink(t, "set", srv, "up")

	conn, err := NewRawUDPConn(cli, 68)
	if err != nil {
		t.Fatalf("NewRawUDPConn(%s, 68): %v", cli, err)
	}
	defer conn.Close()
	upc, ok := conn.(*BroadcastRawUDPConn)
	if !ok {
		t.Fatalf("NewRawUDPConn returned %T, want *BroadcastRawUDPConn", conn)
	}
	client := upc.PacketConn

	sender := listen(t, srv)
	defer sender.Close()

	dhcp := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	valid := frame(28+len(dhcp), dhcp)
	// Distinct payloads: if a dropped frame ever reaches the socket, the payload
	// check below sees the wrong bytes instead of matching the valid reply by luck.
	wrongPort := frame(28+len(dhcp), []byte{0x11, 0x22, 0x33, 0x44})
	binary.BigEndian.PutUint16(wrongPort[22:], 67)
	notUDP := frame(28+len(dhcp), []byte{0x55, 0x66, 0x77, 0x88})
	notUDP[9] = 6 // TCP

	bcast := &packet.Addr{HardwareAddr: net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}
	for _, f := range [][]byte{wrongPort, notUDP, valid} {
		if _, err := sender.WriteTo(f, bcast); err != nil {
			t.Fatalf("send: %v", err)
		}
	}

	// The dropped frames never reach the socket, so the first frame read must be
	// the valid reply. Its dhcp payload sits at offset 28 (20-byte IPv4 header +
	// 8-byte UDP header); the tail may carry Ethernet padding, so compare a slice.
	if err := client.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1500)
	n, _, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatalf("read: %v (the valid reply was dropped by the filter)", err)
	}
	// The frame that arrived must be the valid reply, not a rejected candidate
	// that slipped past an inverted filter: UDP (protocol 17), destination port
	// 68, and the valid frame's own distinct payload.
	if n < 28+len(dhcp) {
		t.Fatalf("read %d bytes, want at least %d", n, 28+len(dhcp))
	}
	if buf[9] != 17 {
		t.Fatalf("received IP protocol %d, want 17 (UDP)", buf[9])
	}
	if dport := binary.BigEndian.Uint16(buf[22:24]); dport != 68 {
		t.Fatalf("received UDP destination port %d, want 68", dport)
	}
	if !bytes.Equal(buf[28:28+len(dhcp)], dhcp) {
		t.Fatalf("received payload %x, want the valid reply %x", buf[28:28+len(dhcp)], dhcp)
	}

	// Nothing else should arrive: the wrong-port and non-UDP frames were
	// dropped. Require a deadline timeout specifically; a socket or device
	// error would otherwise be mistaken for successful filtering.
	if err := client.SetReadDeadline(time.Now().Add(300 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	n2, _, err := client.ReadFrom(buf)
	if err == nil {
		t.Fatalf("a filtered frame reached userspace: %d bytes", n2)
	}
	var netErr net.Error
	if !errors.As(err, &netErr) || !netErr.Timeout() {
		t.Fatalf("second read error = %v, want a timeout", err)
	}
}

func ipLink(t *testing.T, args ...string) {
	t.Helper()
	if out, err := exec.Command("ip", append([]string{"link"}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("ip link %v: %v (%s)", args, err, out)
	}
}

func listen(t *testing.T, iface string) *packet.Conn {
	t.Helper()
	ifc, err := net.InterfaceByName(iface)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := packet.Listen(ifc, packet.Datagram, unix.ETH_P_IP, nil)
	if err != nil {
		t.Fatalf("listen on %s: %v", iface, err)
	}
	return conn
}
