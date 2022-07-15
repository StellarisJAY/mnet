package tcp

import (
	"encoding/binary"
	"fmt"
	"github.com/StellarisJAY/mnet/interface/network"
	"io"
)

/*
LengthFieldBasedDecode
A common way to solve TCP sticky packet problem using LENGTH field in protocol's data packet.
To use this function, you must provide the following arguments.

io.Reader is the source of data, this is normally a TCP connection.
HEADER_LENGTH is the length of header part in your protocol's data packet.
OFFSET is length field's offset in data packet's header, starting at 0.
BYTES is the size of your length field.
network.Packet is an empty Packet to put the parsed data, so it must be pointer type
*/
func LengthFieldBasedDecode(reader io.Reader, headerLength uint32, offset uint32, bytes uint32, packet network.Packet) (network.Packet, error) {
	// check arguments
	if offset >= headerLength || offset < 0 || offset+bytes >= headerLength {
		return nil, fmt.Errorf("length field offset out of range")
	}
	// make a buffer to hold header bytes
	header := make([]byte, headerLength)
	_, err := io.ReadFull(reader, header)
	if err != nil {
		return nil, fmt.Errorf("read header error %v", err)
	}
	// parse length field according to given args
	length := binary.BigEndian.Uint32(header[offset : offset+bytes])
	if length < 0 || length < headerLength {
		return nil, fmt.Errorf("broken header packet error")
	}
	packet.SetLength(length)
	packet.SetHeader(header)
	// read data using length field value
	if length > 0 {
		data := make([]byte, length-headerLength)
		_, err := io.ReadFull(reader, data)
		if err != nil {
			return nil, fmt.Errorf("read data error %v", err)
		}
		packet.SetData(data)
	}
	return packet, nil
}
