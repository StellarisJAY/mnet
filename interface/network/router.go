package network

type HandlerContext interface {
	GetConnection() Connection // get connection from the context
	GetPacket() Packet         // get data packet from context

	SetAttribute(name string, value interface{})  // set attribute in handler
	GetAttribute(name string) (interface{}, bool) // get an attribute

	Send(response Packet) // send a response packet, this method will stop the handler chain
	CloseError(err error) // handler error exit, this method will end handler execution
	IsClosed() bool
	GetError() error
}

type Handler interface {
	PreHandle(ctx HandlerContext)
	Handle(ctx HandlerContext)
	PostHandle(ctx HandlerContext)
	HandleError(err error, ctx HandlerContext)
}

type Router interface {
	StartWorkers()
	Close()
	Execute(ctx HandlerContext)
	Submit(ctx HandlerContext)
	Register(typeCode byte, handler Handler)
}
