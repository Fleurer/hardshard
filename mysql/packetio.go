package mysql

// https://dev.mysql.com/doc/internals/en/mysql-packet.html
// https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_packets.html

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

const (
	MAX_PACKET_PAYLOAD_LENGTH = 1<<24 - 1
)

type PacketIO struct {
	r        io.Reader
	w        io.Writer
	sequence uint8
}

type Packet struct {
	data []byte
}

func NewPacketIO(r io.Reader, w io.Writer) *PacketIO {
	pio := &PacketIO{}
	pio.r = r
	pio.w = w
	pio.sequence = 0
	return pio
}

func NewPacketIOByConn(conn net.Conn) *PacketIO {
	r := bufio.NewReaderSize(conn, 1024*8)
	w := conn
	return NewPacketIO(r, w)
}

func (pio *PacketIO) ReadPacket() ([]byte, error) {
	// [length: byte[3]][sequence_id: byte[1]][playload: byte[length]]
	header := []byte{0, 0, 0, 0}

	_, err := io.ReadFull(pio.r, header)
	if err != nil {
		return nil, ErrBadConn
	}

	length := uint32(header[0]) + uint32(header[1])<<8 + uint32(header[2])<<16
	if length < 1 {
		return nil, fmt.Errorf("invalid payload length %d", length)
	}

	sequence := uint8(header[3])
	if pio.sequence != sequence {
		return nil, fmt.Errorf("invalid sequence %d != %d", sequence, pio.sequence)
	}

	pio.sequence++

	// TODO: reuse the buffer ?
	payload := make([]byte, length)
	_, err = io.ReadFull(pio.r, payload)
	if err != nil {
		return nil, ErrBadConn
	}

	if length < MAX_PACKET_PAYLOAD_LENGTH {
		return payload, nil
	} else {
		// If the payload is larger than or equal to 2**24-1 bytes the length is set to 2**24-1
		// (0xffffff) and a additional packets are sent with the rest of the payload until
		// the payload of a packet is less than 2**24-1 bytes.
		nextPayload, err := pio.ReadPacket()
		if err != nil {
			return nil, err
		}
		return append(payload, nextPayload...), nil
	}
}

func (pio *PacketIO) WriteErrorPacket(err error) {
}

func (pio *PacketIO) WritePacket(payload []byte) {
}
