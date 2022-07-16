package tcp

import (
	"github.com/StellarisJAY/mnet/common"
	"github.com/StellarisJAY/mnet/interface/network"
	"log"
	"net"
	"sync/atomic"
)

// Client in TCP
type Client struct {
	common.BaseClient
	protocol network.Protocol // connection's protocol
}

var nextConnId uint32 = 0

func MakeClient(protocol network.Protocol) *Client {
	c := new(Client)
	c.protocol = protocol
	c.BaseClient = common.MakeBaseClient(protocol, c.newConnection)
	return c
}

// create a new connection to target address
func (c *Client) newConnection(address string) network.Connection {
	// tcp connect
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("create connection error, ", err)
		return nil
	}
	// make connection and start IO loop
	connection := MakeTcpConnection(conn, atomic.AddUint32(&nextConnId, 1), c.protocol, nil, true)
	go connection.Start()
	return connection
}
