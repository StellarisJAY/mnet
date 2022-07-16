package quic

import (
	"github.com/StellarisJAY/mnet/common"
	"github.com/StellarisJAY/mnet/interface/network"
	"github.com/lucas-clemente/quic-go"
)

type Connection struct {
	common.BaseConnection
	stream quic.Stream
}

func MakeQuicConnection(protocol network.Protocol, id uint32, router network.Router, stream quic.Stream, clientSide bool) *Connection {
	return &Connection{common.MakeBaseConnection(protocol, id, router, stream, clientSide), stream}
}

func (c *Connection) Send(packet network.Packet) error {
	protocol := c.BaseConnection.GetProtocol()
	encoded, err := protocol.Encode(packet)
	if err != nil {
		return err
	}
	_, err = c.stream.Write(encoded)
	return err
}

func (c *Connection) RemoteAddr() string {
	return ""
}
