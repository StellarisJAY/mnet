package network

import "io"

type DecodeFunc func(reader io.Reader) Packet
type EncodeFunc func(packet Packet) []byte

type Protocol interface {
	Decoder() DecodeFunc // Protocol's Decoder
	Encoder() EncodeFunc // Protocol's Encoder

	HeaderLength() uint32 // Protocol's HeaderLength in bytes
	Code() uint32         // Protocol's unique ID Code

	Handle(conn Connection, packet Packet) // HandleFunc of this protocol
}
