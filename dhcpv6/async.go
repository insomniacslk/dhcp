package dhcpv6

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// AsyncClient implements an asynchronous DHCPv6 client
type AsyncClient struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LocalAddr    net.Addr
	RemoteAddr   net.Addr
	IgnoreErrors bool

	connection   *net.UDPConn
	cancel       context.CancelFunc
	stopping     *sync.WaitGroup
	receiveQueue chan DHCPv6
	sendQueue    chan DHCPv6
	packetsLock  sync.Mutex
	packets      map[uint32](chan Response)
	errors       chan error
}

// NewAsyncClient creates an asynchronous client
func NewAsyncClient() *AsyncClient {
	return &AsyncClient{
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
	}
}

// OpenForInterface starts the client on the specified interface, replacing
// client LocalAddr with a link-local address of the given interface and
// standard DHCP port (546).
func (c *AsyncClient) OpenForInterface(ifname string, bufferSize int) error {
	addr, err := GetLinkLocalAddr(ifname)
	if err != nil {
		return err
	}
	c.LocalAddr = &net.UDPAddr{IP: *addr, Port: DefaultClientPort, Zone: ifname}
	return c.Open(bufferSize)
}

// Open starts the client
func (c *AsyncClient) Open(bufferSize int) error {
	var (
		addr *net.UDPAddr
		ok   bool
		err  error
	)

	if addr, ok = c.LocalAddr.(*net.UDPAddr); !ok {
		return fmt.Errorf("Invalid local address: %v not a net.UDPAddr", c.LocalAddr)
	}

	// prepare the socket to listen on for replies
	c.connection, err = net.ListenUDP("udp6", addr)
	if err != nil {
		return err
	}
	c.stopping = new(sync.WaitGroup)
	c.sendQueue = make(chan DHCPv6, bufferSize)
	c.receiveQueue = make(chan DHCPv6, bufferSize)
	c.packets = make(map[uint32](chan Response))
	c.packetsLock = sync.Mutex{}
	c.errors = make(chan error)

	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())
	go c.receiverLoop(ctx)
	go c.senderLoop(ctx)

	return nil
}

// Close stops the client
func (c *AsyncClient) Close() {
	// Wait for sender and receiver loops
	c.stopping.Add(2)
	c.cancel()
	c.stopping.Wait()

	close(c.sendQueue)
	close(c.receiveQueue)
	close(c.errors)

	c.connection.Close()
}

// Errors returns a channel where runtime errors are posted
func (c *AsyncClient) Errors() <-chan error {
	return c.errors
}

func (c *AsyncClient) addError(err error) {
	if !c.IgnoreErrors {
		c.errors <- err
	}
}

func (c *AsyncClient) receiverLoop(ctx context.Context) {
	defer func() { c.stopping.Done() }()
	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-c.receiveQueue:
			c.receive(packet)
		}
	}
}

func (c *AsyncClient) senderLoop(ctx context.Context) {
	defer func() { c.stopping.Done() }()
	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-c.sendQueue:
			c.send(packet)
		}
	}
}

func (c *AsyncClient) send(packet DHCPv6) {
	transactionID, err := GetTransactionID(packet)
	if err != nil {
		c.addError(fmt.Errorf("Warning: This should never happen, there is no transaction ID on %s", packet))
		return
	}
	c.packetsLock.Lock()
	f := c.packets[transactionID]
	c.packetsLock.Unlock()

	raddr, err := c.remoteAddr()
	if err != nil {
		f <- NewResponse(nil, err)
		return
	}

	c.connection.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	_, err = c.connection.WriteTo(packet.ToBytes(), raddr)
	if err != nil {
		f <- NewResponse(nil, err)
		return
	}

	c.receiveQueue <- packet
}

func (c *AsyncClient) receive(_ DHCPv6) {
	var (
		oobdata  = []byte{}
		received DHCPv6
	)

	c.connection.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	for {
		buffer := make([]byte, maxUDPReceivedPacketSize)
		n, _, _, _, err := c.connection.ReadMsgUDP(buffer, oobdata)
		if err != nil {
			if err, ok := err.(net.Error); !ok || !err.Timeout() {
				c.addError(fmt.Errorf("Error receiving the message: %s", err))
			}
			return
		}
		received, err = FromBytes(buffer[:n])
		if err != nil {
			// skip non-DHCP packets
			continue
		}
		break
	}

	transactionID, err := GetTransactionID(received)
	if err != nil {
		c.addError(fmt.Errorf("Unable to get a transactionID for %s: %s", received, err))
		return
	}

	c.packetsLock.Lock()
	if f, ok := c.packets[transactionID]; ok {
		delete(c.packets, transactionID)
		f <- NewResponse(received, nil)
	}
	c.packetsLock.Unlock()
}

func (c *AsyncClient) remoteAddr() (*net.UDPAddr, error) {
	if c.RemoteAddr == nil {
		return &net.UDPAddr{IP: AllDHCPRelayAgentsAndServers, Port: DefaultServerPort}, nil
	}

	if addr, ok := c.RemoteAddr.(*net.UDPAddr); ok {
		return addr, nil
	}
	return nil, fmt.Errorf("Invalid remote address: %v not a net.UDPAddr", c.RemoteAddr)
}

// Send inserts a message to the queue to be sent asynchronously.
// Returns a future which resolves to response and error.
func (c *AsyncClient) Send(message DHCPv6, modifiers ...Modifier) Future {
	for _, mod := range modifiers {
		message = mod(message)
	}

	transactionID, err := GetTransactionID(message)
	if err != nil {
		return NewFailureFuture(err)
	}

	f := NewFuture()
	c.packetsLock.Lock()
	c.packets[transactionID] = f
	c.packetsLock.Unlock()
	c.sendQueue <- message
	return f
}

// Exchange executes asynchronously a 4-way DHCPv6 request (SOLICIT,
// ADVERTISE, REQUEST, REPLY).
func (c *AsyncClient) Exchange(solicit DHCPv6, modifiers ...Modifier) Future {
	return c.Send(solicit).OnSuccess(func(advertise DHCPv6) Future {
		request, err := NewRequestFromAdvertise(advertise)
		if err != nil {
			return NewFailureFuture(err)
		}
		return c.Send(request, modifiers...)
	})
}
