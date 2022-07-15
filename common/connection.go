package common

import (
	"context"
	"github.com/StellarisJAY/mnet/interface/network"
	"github.com/StellarisJAY/mnet/route"
	"io"
	"log"
	"sync"
)

type BaseConnection struct {
	id         uint32
	protocol   network.Protocol
	sendBuffer chan network.Packet
	ctx        context.Context
	cancel     context.CancelFunc
	clientSide bool
	router     network.Router
	pending    sync.Map
	readWriter io.ReadWriteCloser
}

func MakeBaseConnection(protocol network.Protocol, id uint32, router network.Router, readWrite io.ReadWriteCloser, clientSide bool) BaseConnection {
	return BaseConnection{
		id:         id,
		protocol:   protocol,
		sendBuffer: make(chan network.Packet, 1<<10),
		clientSide: clientSide,
		router:     router,
		pending:    sync.Map{},
		readWriter: readWrite,
	}
}

func (c *BaseConnection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.readLoop(c.readWriter)
	go c.writeLoop(c.readWriter)

	<-c.ctx.Done()
	c.close(c.readWriter)
}

func (c *BaseConnection) Close() {
	c.cancel()
}

func (c *BaseConnection) ConnectionID() uint32 {
	return c.id
}

func (c *BaseConnection) Context() context.Context {
	return c.ctx
}

func (c *BaseConnection) RemoteAddr() string {
	panic("Call remote addr on BaseConnection is not allowed")
}

func (c *BaseConnection) Send(packet network.Packet) error {
	panic("Call Send on BaseConnection is not allowed")
}

func (c *BaseConnection) SendBuffered(packet network.Packet) {
	c.sendBuffer <- packet
}

func (c *BaseConnection) AddPending(id uint32, wait chan network.Packet) {
	c.pending.Store(id, wait)
}

func (c *BaseConnection) FinishPending(id uint32, packet network.Packet) {
	p, loaded := c.pending.LoadAndDelete(id)
	if loaded {
		wait := p.(chan network.Packet)
		wait <- packet
	}
}

func (c *BaseConnection) GetProtocol() network.Protocol {
	return c.protocol
}

func (c *BaseConnection) readLoop(reader io.Reader) {
	for {
		packet, err := c.protocol.Decode(reader)
		if err != nil {
			c.cancel()
			break
		}
		if c.clientSide {
			c.FinishPending(packet.ID(), packet)
		} else {
			ctx := route.MakeHandlerContext(c, packet)
			c.router.Submit(ctx)
		}
	}
}

func (c *BaseConnection) writeLoop(writer io.Writer) {
	for {
		select {
		case packet, ok := <-c.sendBuffer:
			// channel closed, break write loop
			if !ok {
				return
			}
			// encode packet and send
			encoded, err := c.protocol.Encode(packet)
			if err != nil {
				log.Println("encode error for packet: ", packet, " , error: ", err)
				continue
			}
			_, err = writer.Write(encoded)
			if err != nil {
				c.cancel()
				return
			}
		}
	}
}

func (c *BaseConnection) close(closer io.Closer) {
	// close channel
	close(c.sendBuffer)
	// close connection
	_ = closer.Close()
}
