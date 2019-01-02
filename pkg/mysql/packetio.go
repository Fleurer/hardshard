package mysql

// Protocol Basics: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_packets.html
// OK_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_ok_packet.html
// ERR_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_err_packet.html
// EOF_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_eof_packet.html

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

const (
	MAX_PACKET_PAYLOAD_LENGTH = 1<<24 - 1 // 16mb
)

type PacketIO struct {
	r        io.Reader
	w        io.Writer
	Sequence uint8
}

func NewPacketIO(r io.Reader, w io.Writer) *PacketIO {
	pio := &PacketIO{}
	pio.r = r
	pio.w = w
	pio.Sequence = 0
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

	sequence := uint8(header[3])
	if pio.Sequence != sequence {
		return nil, fmt.Errorf("invalid sequence: packet sequence(%d) != pio.Sequence(%d)", sequence, pio.Sequence)
	}

	pio.Sequence++

	// 恰好 16Mb 的 packet 后面会追加一个长度为 0 的 packet
	if length == 0 {
		return []byte{}, nil
	}

	// TODO: reuse the buffer ?
	payload := make([]byte, length)
	_, err = io.ReadFull(pio.r, payload)
	if err != nil {
		return nil, ErrBadConn
	}

	if length < MAX_PACKET_PAYLOAD_LENGTH {
		return payload, nil
	} else {
		// https://dev.mysql.com/doc/internals/en/sending-more-than-16mbyte.html
		// If the payload is larger than or equal to 2**24-1 bytes the length is set to 2**24-1
		// (0xffffff) and a additional packets are sent with the rest of the payload until
		// the payload of a packet is less than 2**24-1 bytes.
		nextPayload, err := pio.ReadPacket()
		if err != nil {
			return nil, err
		}
		// TODO: 这里性能不好，埋指标计数 slow path
		return append(payload, nextPayload...), nil
	}
}

func (pio *PacketIO) WritePacket(payload []byte) error {
	for len(payload) > 0 {
		length := len(payload)
		if length >= MAX_PACKET_PAYLOAD_LENGTH {
			length = MAX_PACKET_PAYLOAD_LENGTH
		}

		header := []byte{byte(length), byte(length >> 8), byte(length >> 16), pio.Sequence}
		n, err := pio.w.Write(header)
		if err != nil || n != 4 {
			return ErrBadConn
		}
		chunk := payload[0:length]
		n, err = pio.w.Write(chunk)
		if err != nil || n != len(chunk) {
			return ErrBadConn
		}
		pio.Sequence++

		// 如果 payload 恰好等于 16MB，后面追加一个长度为 0 的 packet
		// https://github.com/Qihoo360/Atlas/blob/128b0544cefc800366f70e534c5130f35574721c/src/network-mysqld.c#L364
		if len(payload) == MAX_PACKET_PAYLOAD_LENGTH {
			header := []byte{0, 0, 0, pio.Sequence}
			n, err := pio.w.Write(header)
			if err != nil || n != 4 {
				return ErrBadConn
			}
			pio.Sequence++
			return nil
		}

		payload = payload[length:]
	}
	return nil
}
