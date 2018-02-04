package mysql

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

const (
	MAX_PACKET_LEN = 0xFFFFFF // 16MB
)

type PacketIO struct {
	r        io.Reader
	w        io.Writer
	sequence int64
}

func NewPacketIO(r io.Reader, w io.Writer) *PacketIO {
	p := &PacketIO{}
	p.r = r
	p.w = w
	p.sequence = 0
	return p
}

func NewPacketIOByConn(conn net.Conn) *PacketIO {
	r := bufio.NewReaderSize(conn, 1024)
	w := conn
	return NewPacketIO(r, w)
}

func (p *PacketIO) ReadPacket() ([]byte, error) {
	buf := make([]byte, 512)
	n, _ := p.r.Read(buf)

	fmt.Printf("%d %s", n, buf[:n])

	return buf[:n], nil
}

func (p *PacketIO) WriteErrorPacket(err error) {
}
