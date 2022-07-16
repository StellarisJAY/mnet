package common

import (
	"github.com/StellarisJAY/mnet/interface/network"
	"github.com/StellarisJAY/mnet/util/pool"
	"sync"
)

type BaseClient struct {
	protocol      network.Protocol              // connection's protocol
	connPools     map[string]*pool.SharablePool // connection pools
	mutex         sync.Mutex                    // lock when initializing connection pool
	newConnection func(address string) network.Connection
}

func MakeBaseClient(protocol network.Protocol, newConnection func(address string) network.Connection) BaseClient {
	return BaseClient{
		protocol:      protocol,
		connPools:     make(map[string]*pool.SharablePool),
		mutex:         sync.Mutex{},
		newConnection: newConnection,
	}
}

func (c *BaseClient) Oneway(address string, packet network.Packet) error {
	connection := c.getConnection(address)
	defer c.returnConnection(address, connection)
	// send packet, discard response
	return connection.Send(packet)
}

func (c *BaseClient) Request(address string, packet network.Packet) (response network.Packet, err error) {
	// send and gets a channel
	wait, err := c.Future(address, packet)
	if err != nil {
		return nil, err
	}
	// wait for response
	response = <-wait
	return
}

func (c *BaseClient) Future(address string, packet network.Packet) (chan network.Packet, error) {
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

func (c *BaseClient) Async(address string, packet network.Packet, callback func(packet network.Packet, err error)) {
	future, err := c.Future(address, packet)
	// starts a goroutine to receive response and call callback
	go func(future chan network.Packet, err error) {
		if err != nil {
			callback(nil, err)
		} else {
			response := <-future
			callback(response, nil)
		}
	}(future, err)
}

// get a connection from the connection pool of this address
func (c *BaseClient) getConnection(address string) network.Connection {
	c.mutex.Lock()
	// create conn pool
	p, ok := c.connPools[address]
	if !ok {
		p = pool.NewSharablePool(20, 20, func() interface{} {
			return c.newConnection(address)
		})
		c.connPools[address] = p
	}
	c.mutex.Unlock()
	// get connection
	return p.Get().(network.Connection)
}

// return the connection to its pool after using
func (c *BaseClient) returnConnection(address string, conn network.Connection) {
	p, ok := c.connPools[address]
	if ok {
		p.Put(conn)
	}
}
