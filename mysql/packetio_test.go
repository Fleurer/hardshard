package mysql

import (
	"bytes"
	"testing"

	"github.com/Fleurer/hardshard/mysql"
)

func TestNewPacketIO(t *testing.T) {
	r := bytes.NewBufferString("hello\n")
	w := bytes.NewBufferString("")

	mysql.NewPacketIO(r, w)
}

func TestPacketIOWithComQuit(t *testing.T) {
	comQuit := []byte{01, 00, 00, 00, 01}
	buf := bytes.NewBuffer(comQuit)
	p := mysql.NewPacketIO(buf, buf)
	payload, err := p.ReadPacket()
	if err != nil {
		t.Fatalf("err on p.ReadPacket: %s", err)
	}
	if !bytes.Equal(payload, []byte{1}) {
		t.Fatalf("payload not equal")
	}
}
