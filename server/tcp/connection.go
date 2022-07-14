package tcp

import (
	"context"
	"fmt"
	"github.com/StellarisJAY/mnet/interface/network"
	"log"
	"net"
)

type Connection struct {
	conn net.Conn
	id   uint32

	ctx        context.Context
	cancel     context.CancelFunc
	protocol   network.Protocol
	sendBuffer chan network.Packet
}

func MakeTcpConnection(conn net.Conn, id uint32, protocol network.Protocol) *Connection {
	return &Connection{
		conn:       conn,
		id:         id,
		ctx:        nil,
		protocol:   protocol,
		sendBuffer: make(chan network.Packet, 1024),
	}
}

func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.readLoop()
	go c.writeLoop()

	<-c.ctx.Done()
	c.close()
}

func (c *Connection) Close() {
	c.cancel()
}

func (c *Connection) ConnectionID() uint32 {
	return c.id
}

func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Connection) Send(packet network.Packet) error {
	encoded, err := c.protocol.Encode(packet)
	if err != nil {
		return fmt.Errorf("encode error: %v", err)
	}
	_, err = c.conn.Write(encoded)
	return err
}

func (c *Connection) SendBuffered(packet network.Packet) {
	c.sendBuffer <- packet
}

func (c *Connection) readLoop() {
	for {
		packet, err := c.protocol.Decode(c.conn)
		if err != nil {
			break
		}
		c.protocol.HandleWithWorker(c, packet)
	}
}

func (c *Connection) writeLoop() {
	for {
		select {
		case packet, ok := <-c.sendBuffer:
			// channel closed, break write loop
			if !ok {
				break
			}
			// encode packet and send
			encoded, err := c.protocol.Encode(packet)
			if err != nil {
				log.Println("encode error for packet: ", packet, " , error: ", err)
				continue
			}
			_, err = c.conn.Write(encoded)
			if err != nil {
				break
			}
		}
	}
}

func (c *Connection) close() {
	// close channel
	close(c.sendBuffer)
	// close connection
	_ = c.conn.Close()
}
