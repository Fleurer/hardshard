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

func NewPacketIO(r io.Reader, w io.Writer) *PacketIO {
	p := &PacketIO{}
	p.r = r
	p.w = w
	p.sequence = 0
	return p
}

func NewPacketIOByConn(conn net.Conn) *PacketIO {
	r := bufio.NewReaderSize(conn, 1024*8)
	w := conn
	return NewPacketIO(r, w)
}

func (p *PacketIO) ReadPacket() ([]byte, error) {
	// [length: byte[3]][sequence_id: byte[1]][playload: byte[length]]
	header := []byte{0, 0, 0, 0}

	_, err := io.ReadFull(p.r, header)
	if err != nil {
		return nil, ErrBadConn
	}

	length := uint32(header[0]) + uint32(header[1])<<8 + uint32(header[2])<<16
	if length < 1 {
		return nil, fmt.Errorf("invalid payload length %d", length)
	}

	sequence := uint8(header[3])
	if p.sequence != sequence {
		return nil, fmt.Errorf("invalid sequence %d != %d", sequence, p.sequence)
	}

	p.sequence++

	// TODO: reuse the buffer ?
	payload := make([]byte, length)
	_, err = io.ReadFull(p.r, payload)
	if err != nil {
		return nil, ErrBadConn
	}

	if length < MAX_PACKET_PAYLOAD_LENGTH {
		return payload, nil
	} else {
		nextPayload, err := p.ReadPacket()
		if err != nil {
			return nil, err
		}
		return append(payload, nextPayload...), nil
	}
}

func (p *PacketIO) WriteErrorPacket(err error) {
}
