package network

import "context"

type Connection interface {
	Start() // start connection's IO loop
	Close() // close connection

	Context() context.Context // get the Context of this connection
	RemoteAddr() string       // get the remote address of this connection
	ConnectionID() uint32     // get connection's ID

	Send(packet Packet) error   // send a packet through connection, this function will block until written
	SendBuffered(packet Packet) // send a packet non-blocking

	AddPending(id uint32, wait chan Packet)
	FinishPending(id uint32, packet Packet)
}
