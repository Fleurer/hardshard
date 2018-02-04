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

func TestPacketIOWithComQuery1(t *testing.T) {
	// test sample from https://dev.mysql.com/doc/internals/en/example-one-mysql-packet.html
	data := []byte{
		0x2e, 0x00, 0x00, 0x00, 0x03, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x20, 0x22, 0x30, 0x31, 0x32,
		0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
		0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,
		0x35, 0x22,
	}
	buf := bytes.NewBuffer(data)
	p := mysql.NewPacketIO(buf, buf)
	payload, err := p.ReadPacket()
	if err != nil {
		t.Fatalf("err on p.ReadPacket: %s", err)
	}
	if len(payload) != 46 {
		t.Fatalf("payload len mismatch: %d", len(payload))
	}
	expectedPayload := []byte("\x03select \"012345678901234567890123456789012345\"")
	if !bytes.Equal(payload, expectedPayload) {
		t.Fatalf("payload not equal\npayload: %b\nexpected: %b", payload, expectedPayload)
	}
}
