package network

type Client interface {
	// Oneway send a packet without response
	Oneway(address string, packet Packet) error
	// Request send a packet and wait for response
	Request(address string, packet Packet) (response Packet, err error)
	// Future send a packet and get a channel to receive response
	Future(address string, packet Packet) (chan Packet, error)
	// Async send a packet and do nothing, client will call callback automatically when receiving response
	Async(address string, packet Packet, callback func(packet Packet, err error))
}
