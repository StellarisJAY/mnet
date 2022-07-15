package tcp

import (
	"github.com/StellarisJAY/mnet/common"
	"github.com/StellarisJAY/mnet/interface/network"
	"net"
)

type Connection struct {
	common.BaseConnection
	conn net.Conn
}

func MakeTcpConnection(conn net.Conn, id uint32, protocol network.Protocol, router network.Router, clientSide bool) *Connection {
	c := &Connection{
		common.MakeBaseConnection(protocol, id, router, conn, clientSide),
		conn,
	}
	return c
}

func (c *Connection) Send(packet network.Packet) error {
	protocol := c.BaseConnection.GetProtocol()
	encoded, err := protocol.Encode(packet)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(encoded)
	return err
}

func (c *Connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}
