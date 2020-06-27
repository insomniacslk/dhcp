// this tests nclient4 with lease and release

package nclient4

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/hugelgupf/socketpair"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

type testLeaseKey struct {
	mac       net.HardwareAddr
	idOptions dhcpv4.Options
}

func (lk testLeaseKey) compare(b testLeaseKey) bool {
	if !bytes.Equal(lk.idOptions.ToBytes(), b.idOptions.ToBytes()) {
		return false
	}
	for i := 0; i < 6; i++ {
		if lk.mac[i] != b.mac[i] {
			return false
		}
	}
	return true
}

//this represents one test case
type testServerLease struct {
	key          *testLeaseKey
	assignedAddr net.IP
	ShouldFail   bool //expected result
}

type testServerLeaseList struct {
	list            []*testServerLease
	clientIDOptions dhcpv4.OptionCodeList
}

func newtestServerLeaseList(l dhcpv4.OptionCodeList) *testServerLeaseList {
	r := &testServerLeaseList{}
	r.clientIDOptions = l
	return r
}

func (sll testServerLeaseList) get(k *testLeaseKey) *testServerLease {
	for i := range sll.list {
		if sll.list[i].key.compare(*k) {
			return sll.list[i]
		}
	}
	return nil
}

func (sll *testServerLeaseList) getKey(m *dhcpv4.DHCPv4) *testLeaseKey {
	key := &testLeaseKey{}
	key.mac = m.ClientHWAddr
	key.idOptions = make(dhcpv4.Options)
	for _, optioncode := range sll.clientIDOptions {
		v := m.Options.Get(optioncode)
		key.idOptions.Update(dhcpv4.OptGeneric(optioncode, v))
	}
	return key

}

//use following setting to handle DORA
//server-id: 1.2.3.4
//subnet-mask: /24
//return address from sll.list
func (sll *testServerLeaseList) testLeaseDORAHandle(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) error {
	reply, err := dhcpv4.NewReplyFromRequest(m)
	if err != nil {
		return fmt.Errorf("NewReplyFromRequest failed: %v", err)
	}
	svrIP := net.ParseIP("1.2.3.4")
	reply.UpdateOption(dhcpv4.OptServerIdentifier(svrIP.To4()))
	reply.UpdateOption(dhcpv4.OptSubnetMask(net.IPv4Mask(255, 255, 255, 0)))
	//build lease key
	key := sll.getKey(m)
	clease := sll.get(key)
	if clease == nil {
		return fmt.Errorf("unable to find the lease")
	}
	reply.YourIPAddr = clease.assignedAddr
	switch mt := m.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
	case dhcpv4.MessageTypeRequest:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))

	default:
		return fmt.Errorf("Unhandled message type: %v", mt)
	}

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		return fmt.Errorf("Cannot reply to client: %v", err)
	}
	return nil
}

//return check list for options must and may in the release msg according to RFC2131,section 4.4.1
func (sll *testServerLeaseList) getCheckList() (mustHaveOpts, mayHaveOpts map[uint8]bool) {
	mustHaveOpts = make(map[uint8]bool)
	mayHaveOpts = make(map[uint8]bool)
	mustHaveOpts[dhcpv4.OptionDHCPMessageType.Code()] = false
	mustHaveOpts[dhcpv4.OptionServerIdentifier.Code()] = false

	for _, o := range sll.clientIDOptions {
		mustHaveOpts[o.Code()] = false
	}
	mayHaveOpts[dhcpv4.OptionClassIdentifier.Code()] = false
	mayHaveOpts[dhcpv4.OptionMessage.Code()] = false
	return

}

//check request message according to RFC2131, section 4.4.1
func (sll *testServerLeaseList) testLeaseReleaseHandle(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) error {

	if m.HopCount != 0 {
		return fmt.Errorf("hop count is %v, should be 0", m.HopCount)
	}
	if m.NumSeconds != 0 {
		return fmt.Errorf("seconds is %v, should be 0", m.NumSeconds)
	}
	if m.Flags != 0 {
		return fmt.Errorf("flags is %v, should be 0", m.Flags)
	}
	key := sll.getKey(m)
	clease := sll.get(key)
	if clease == nil {
		return fmt.Errorf("can't find the lease")
	}
	if !m.ClientIPAddr.Equal(clease.assignedAddr) {
		return fmt.Errorf("client IP is %v, expecting %v", m.ClientIPAddr, clease.assignedAddr)
	}
	if !m.YourIPAddr.Equal(net.ParseIP("0.0.0.0")) {
		return fmt.Errorf("your IP is %v, expect 0", m.YourIPAddr)
	}
	if !m.GatewayIPAddr.Equal(net.ParseIP("0.0.0.0")) {
		return fmt.Errorf("gateway IP is %v, expect 0", m.GatewayIPAddr)
	}
	mustlist, maylist := sll.getCheckList()
	for o := range m.Options {
		foundInMust := false
		foundInMay := false
		if _, ok := mustlist[o]; ok {
			mustlist[o] = true
			foundInMust = true
		}
		if _, ok := maylist[o]; ok {
			foundInMay = true
		}
		if !foundInMay && !foundInMust {
			return fmt.Errorf("option %v is not allowed in DHCP release msg", o)
		}
	}
	for o, got := range mustlist {
		if !got {
			return fmt.Errorf("option %v is missing in  DHCP release msg", o)
		}
	}

	if !net.IP(m.Options.Get(dhcpv4.OptionServerIdentifier)).Equal(net.ParseIP("1.2.3.4")) {
		return fmt.Errorf("release misses servier id option=1.2.3.4")
	}
	return nil
}

func (sll *testServerLeaseList) handle(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	if m == nil {
		log.Fatal("Packet is nil!")
	}
	if m.OpCode != dhcpv4.OpcodeBootRequest {
		log.Fatal("Not a BootRequest!")
	}
	var err error
	switch m.MessageType() {
	case dhcpv4.MessageTypeDiscover, dhcpv4.MessageTypeRequest:
		err = sll.testLeaseDORAHandle(conn, peer, m)
		if err != nil {
			log.Printf("svr failed to handle DORA,%v", err)
		}

	case dhcpv4.MessageTypeRelease:
		err = sll.testLeaseReleaseHandle(conn, peer, m)
		if err != nil {
			log.Printf("svr failed to handle release,%v", err)
		}
	default:
		log.Printf("svr got unexpeceted message type %v", m.MessageType())
	}
}

func (sll *testServerLeaseList) runTest(t *testing.T) {
	for _, l := range sll.list {
		t.Logf("running lease test case for mac %v", l.key.mac)
		// Fake PacketConn connection.
		//note can't reuse conn between different clients, because there is currently
		//no way to stop a client's reciev loop
		clientRawConn, serverRawConn, err := socketpair.PacketSocketPair()
		if err != nil {
			panic(err)
		}
		clientConn := NewBroadcastUDPConn(clientRawConn, &net.UDPAddr{Port: ClientPort})
		serverConn := NewBroadcastUDPConn(serverRawConn, &net.UDPAddr{Port: ServerPort})
		s, err := server4.NewServer("", nil, sll.handle, server4.WithConn(serverConn))
		if err != nil {
			t.Fatal(err)
		}
		go func() {
			_ = s.Serve()
		}()
		clnt, err := testCreateClientWithServerLease(clientConn, l)
		if err != nil {
			t.Fatal(err)
		}
		modList := []dhcpv4.Modifier{}
		for op, val := range l.key.idOptions {
			modList = append(modList, dhcpv4.WithOption(dhcpv4.OptGeneric(dhcpv4.GenericOptionCode(op), val)))
		}
		chkerr := func(err error, shouldfail bool, t *testing.T) bool {
			if err != nil {
				if !shouldfail {
					t.Fatal(err)
				} else {
					t.Logf("case failed as expected,%v", err)
					return false
				}
			}
			return true
		}

		_, lease, err := clnt.Request(context.Background(), modList...)
		keepgoing := chkerr(err, l.ShouldFail, t)
		if keepgoing {
			err = clnt.Release(lease)
			chkerr(err, l.ShouldFail, t)
		}
	}

}

func testCreateClientWithServerLease(conn net.PacketConn, sl *testServerLease) (*Client, error) {
	clntModList := []ClientOpt{WithRetry(1), WithTimeout(2 * time.Second)}
	clntModList = append(clntModList, WithHWAddr(sl.key.mac))
	var idoptlist dhcpv4.OptionCodeList
	for op := range sl.key.idOptions {
		idoptlist.Add(dhcpv4.GenericOptionCode(op))
	}
	clntModList = append(clntModList, WithClientIDOptions(idoptlist))
	clnt, err := NewWithConn(conn, sl.key.mac, clntModList...)
	if err != nil {
		return nil, fmt.Errorf("failed to create dhcpv4 client,%v", err)
	}
	return clnt, nil
}

func TestLease(t *testing.T) {
	//test data set 1
	var idoptlist dhcpv4.OptionCodeList
	idoptlist.Add(dhcpv4.OptionClientIdentifier)
	sll := newtestServerLeaseList(idoptlist)
	sll.list = []*testServerLease{
		&testServerLease{
			assignedAddr: net.ParseIP("192.168.1.1"),
			key: &testLeaseKey{
				mac: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 1, 1},
				idOptions: dhcpv4.Options{
					uint8(dhcpv4.OptionClientIdentifier): []byte("client-1"),
				},
			},
		},

		&testServerLease{
			assignedAddr: net.ParseIP("192.168.1.2"),
			key: &testLeaseKey{
				mac: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 1, 2},
				idOptions: dhcpv4.Options{
					uint8(dhcpv4.OptionClientIdentifier): []byte("client-2"),
				},
			},
		},
		//negative case
		&testServerLease{
			assignedAddr: net.ParseIP("192.168.2.2"),
			key: &testLeaseKey{
				mac:       net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 1, 3},
				idOptions: dhcpv4.Options{},
			},
			ShouldFail: true,
		},
	}
	sll.runTest(t)
	//test data set 2
	idoptlist = dhcpv4.OptionCodeList{}
	sll = newtestServerLeaseList(idoptlist)
	sll.list = []*testServerLease{
		&testServerLease{
			assignedAddr: net.ParseIP("192.168.2.1"),
			key: &testLeaseKey{
				mac:       net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 2, 1},
				idOptions: dhcpv4.Options{},
			},
		},

		&testServerLease{
			assignedAddr: net.ParseIP("192.168.2.2"),
			key: &testLeaseKey{
				mac:       net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 2, 2},
				idOptions: dhcpv4.Options{},
			},
		},
		//negative case
		&testServerLease{
			assignedAddr: net.ParseIP("192.168.2.2"),
			key: &testLeaseKey{
				mac: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 1, 3},
				idOptions: dhcpv4.Options{
					uint8(dhcpv4.OptionClientIdentifier): []byte("client-fake"),
				},
			},
			ShouldFail: true,
		},
	}
	sll.runTest(t)
}
