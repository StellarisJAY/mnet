package network

import "io"

type Protocol interface {
	Decode(reader io.Reader) (Packet, error) // Protocol's Decoder
	Encode(packet Packet) ([]byte, error)    // Protocol's Encoder

	HeaderLength() uint32 // Protocol's HeaderLength in bytes
	Code() byte           // Protocol's unique ID Code

	Handle(conn Connection, packet Packet)           // HandleFunc of this protocol
	HandleWithWorker(conn Connection, packet Packet) // use protocol's worker pool to handle
}
