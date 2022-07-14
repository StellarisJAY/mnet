package network

// Packet data packet that holds message data
type Packet interface {
	Length() uint32 // Packet Length in bytes, header + Data
	Type() byte     // Packet type code
	ID() uint32     // Packet ID
	Data() []byte   // Packet Data
}
