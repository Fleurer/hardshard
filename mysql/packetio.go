package mysql

import (
	"bufio"
	"io"
)

const (
	MAX_PACKET_LEN = 2**24 - 1
)

type PacketIO struct {
	r *bufio.Reader
	w io.Writer
}

func NewPacketIO(r io.Reader, w io.Writer) PacketIO {
	p := PacketIO{}
	p.r = bufio.NewReaderSize(r, 1024)
	p.w = w
	return p
}
