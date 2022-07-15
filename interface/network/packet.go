package network

// Packet data packet that holds message data
type Packet interface {
	Length() uint32 // Packet Length in bytes, header + Data
	Type() byte     // Packet type code
	ID() uint32     // Packet ID
	Data() []byte   // Packet Data
	Header() []byte // the original Header slice

	SetLength(length uint32)
	SetType(typeCode byte)
	SetID(id uint32)
	SetData(data []byte)
	SetHeader(header []byte)
}
