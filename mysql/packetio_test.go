package mysql

import (
	"bytes"
	"testing"
)

func TestNewPacketIO(t *testing.T) {
	r := bytes.NewBufferString("hello\n")
	w := bytes.NewBufferString("")

	NewPacketIO(r, w)
}

func TestPacketIOWithComQuit(t *testing.T) {
	comQuit := []byte{01, 00, 00, 00, 01}
	buf := bytes.NewBuffer(comQuit)
	p := NewPacketIO(buf, buf)
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
	p := NewPacketIO(buf, buf)
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

func TestWritePacket1(t *testing.T) {
	buf := bytes.NewBufferString("")
	pio := NewPacketIO(buf, buf)
	payload := []byte("\x03select \"012345678901234567890123456789012345\"")
	err := pio.WritePacket(payload)
	if err != nil {
		t.Fatalf("err on WritePacket: %s", err)
	}
	if !bytes.Equal(buf.Bytes()[0:4], []byte{46, 0, 0, 0}) {
		t.Fatalf("invalid header: %v", buf.Bytes()[0:4])
	}
	if !bytes.Equal(buf.Bytes()[4:], payload) {
		t.Fatalf("invalid payload", buf.Bytes()[4:])
	}
}

func TestWritePacket2(t *testing.T) {
	buf := bytes.NewBufferString("")
	pio := NewPacketIO(buf, buf)
	payload := make([]byte, MAX_PACKET_PAYLOAD_LENGTH*2+1)
	for i, _ := range payload {
		payload[i] = '6'
	}
	err := pio.WritePacket(payload)
	if err != nil {
		t.Fatalf("err on WritePacket: %s", err)
	}
	if buf.Len() != MAX_PACKET_PAYLOAD_LENGTH*2+1+12 {
		t.Fatalf("mismatch len(buf): %v", buf.Len())
	}
	bbuf := buf.Bytes()[0:4]
	if !bytes.Equal(bbuf, []byte{255, 255, 255, 0}) {
		t.Fatalf("invalid header: %v", bbuf)
	}
	bbuf = buf.Bytes()[MAX_PACKET_PAYLOAD_LENGTH+4 : MAX_PACKET_PAYLOAD_LENGTH+8]
	if !bytes.Equal(bbuf, []byte{255, 255, 255, 1}) {
		t.Fatalf("invalid header: %v", bbuf)
	}
	bbuf = buf.Bytes()[MAX_PACKET_PAYLOAD_LENGTH*2+8 : MAX_PACKET_PAYLOAD_LENGTH*2+12]
	if !bytes.Equal(bbuf, []byte{1, 0, 0, 2}) {
		t.Fatalf("invalid header: %v", bbuf)
	}
}

func TestWritePacket3(t *testing.T) {
	buf := bytes.NewBufferString("")
	pio := NewPacketIO(buf, buf)
	payload := make([]byte, MAX_PACKET_PAYLOAD_LENGTH*2)
	for i, _ := range payload {
		payload[i] = '6'
	}
	err := pio.WritePacket(payload)
	if err != nil {
		t.Fatalf("err on WritePacket: %s", err)
	}
	if buf.Len() != MAX_PACKET_PAYLOAD_LENGTH*2+8 {
		t.Fatalf("mismatch len(buf): %v", buf.Len())
	}
	if !bytes.Equal(buf.Bytes()[0:4], []byte{255, 255, 255, 0}) {
		t.Fatalf("invalid header: %v", buf.Bytes()[0:4])
	}
}
