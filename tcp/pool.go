package tcp

import (
	"github.com/StellarisJAY/mnet/interface/network"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

var nextConnId uint32 = 0

// create a new connection to target address
func (c *Client) newConnection(address string) network.Connection {
	// tcp connect
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("create connection error, ", err)
		return nil
	}
	// make connection and start IO loop
	connection := MakeTcpConnection(conn, atomic.AddUint32(&nextConnId, 1), c.protocol, true)
	go connection.Start()
	return connection
}

// get a connection from the connection pool of this address
func (c *Client) getConnection(address string) network.Connection {
	c.mutex.Lock()
	// create conn pool
	pool, ok := c.connPools[address]
	if !ok {
		pool = &sync.Pool{New: func() interface{} {
			return c.newConnection(address)
		}}
		c.connPools[address] = pool
	}
	c.mutex.Unlock()
	// get connection
	return pool.Get().(network.Connection)
}

// return the connection to its pool after using
func (c *Client) returnConnection(address string, conn network.Connection) {
	p, ok := c.connPools[address]
	if ok {
		p.Put(conn)
	}
}
