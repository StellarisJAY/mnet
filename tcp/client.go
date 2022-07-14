package tcp

import (
	"github.com/StellarisJAY/mnet/interface/network"
	"sync"
)

// Client in TCP
type Client struct {
	protocol  network.Protocol      // connection's protocol
	connPools map[string]*sync.Pool // connection pools
	mutex     sync.Mutex            // lock when initializing connection pool
}

func MakeTcpClient(protocol network.Protocol) *Client {
	return &Client{
		protocol:  protocol,
		connPools: make(map[string]*sync.Pool),
		mutex:     sync.Mutex{},
	}
}

func (c *Client) Oneway(address string, packet network.Packet) error {
	connection := c.getConnection(address)
	defer c.returnConnection(address, connection)
	// send packet, discard response
	return connection.Send(packet)
}

func (c *Client) Request(address string, packet network.Packet) (response network.Packet, err error) {
	// send and gets a channel
	wait, err := c.Future(address, packet)
	if err != nil {
		return nil, err
	}
	// wait for response
	response = <-wait
	return
}

func (c *Client) Future(address string, packet network.Packet) (chan network.Packet, error) {
	connection := c.getConnection(address)
	defer c.returnConnection(address, connection)
	// make response channel and put into connection's pending map
	wait := make(chan network.Packet)
	connection.AddPending(packet.ID(), wait)
	err := connection.Send(packet)
	if err != nil {
		return nil, err
	}
	return wait, nil
}

func (c *Client) Async(address string, packet network.Packet, callback func(packet network.Packet, err error)) {
	connection := c.getConnection(address)
	err := connection.Send(packet)
	wait := make(chan network.Packet)
	connection.AddPending(packet.ID(), wait)
	c.returnConnection(address, connection)
	if err != nil {
		callback(nil, err)
	}
	// starts a goroutine to receive response and call callback
	go func() {
		response := <-wait
		callback(response, nil)
	}()
}
