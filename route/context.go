package route

import "github.com/StellarisJAY/mnet/interface/network"

type HandlerContext struct {
	conn   network.Connection
	packet network.Packet

	attributes map[string]interface{}
	err        error
	closed     bool
}

func MakeHandlerContext(conn network.Connection, pack network.Packet) *HandlerContext {
	return &HandlerContext{
		conn:       conn,
		packet:     pack,
		attributes: make(map[string]interface{}),
		err:        nil,
		closed:     false,
	}
}

func (h *HandlerContext) GetConnection() network.Connection {
	return h.conn
}

func (h *HandlerContext) GetPacket() network.Packet {
	return h.packet
}

func (h *HandlerContext) SetAttribute(name string, value interface{}) {
	h.attributes[name] = value
}

func (h *HandlerContext) GetAttribute(name string) (interface{}, bool) {
	v, ok := h.attributes[name]
	return v, ok
}

func (h *HandlerContext) Send(response network.Packet) {
	h.conn.SendBuffered(response)
	h.closed = true
}

func (h *HandlerContext) CloseError(err error) {
	h.err = err
	h.closed = true
}

func (h *HandlerContext) GetError() error {
	return h.err
}

func (h *HandlerContext) IsClosed() bool {
	return h.closed
}
