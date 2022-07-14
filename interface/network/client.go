package network

type Client interface {
	Oneway(address string, packet Packet) error
	Request(address string, packet Packet) (response Packet, err error)
	Future(address string, packet Packet) (chan Packet, error)
	Async(address string, packet Packet, callback func(packet Packet, err error))
}
