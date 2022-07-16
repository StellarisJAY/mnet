package quic

import (
	"context"
	"crypto/tls"
	"github.com/StellarisJAY/mnet/common"
	"github.com/StellarisJAY/mnet/interface/network"
	"github.com/lucas-clemente/quic-go"
	"log"
	"sync/atomic"
)

type Client struct {
	common.BaseClient
	tlsConfig *tls.Config
	protocol  network.Protocol
}

var nextConnId uint32 = 0

func (c *Client) newQuicConnection(address string) network.Connection {
	conn, err := quic.DialAddr(address, c.tlsConfig, nil)
	if err != nil {
		log.Println("Quic dial address error: ", err)
		return nil
	}
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Println("Quic open stream error ", err)
		return nil
	}
	connection := MakeQuicConnection(c.protocol, atomic.AddUint32(&nextConnId, 1), nil, stream, true)
	go connection.Start()
	return connection
}

func MakeClient(protocol network.Protocol, tlsConfig *tls.Config) *Client {
	c := new(Client)
	c.tlsConfig = tlsConfig
	c.BaseClient = common.MakeBaseClient(protocol, c.newQuicConnection)
	c.protocol = protocol
	return c
}
