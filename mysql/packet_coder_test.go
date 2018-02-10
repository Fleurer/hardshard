package mysql

import (
	"bytes"
	"fmt"
	"testing"
)

func setupPacketCoder() *PacketCoder {
	c := &PacketCoder{capabilities: CLIENT_PROTOCOL_41}
	return c
}

func TestEncodeOK(t *testing.T) {
	c := setupPacketCoder()
	payload := c.encodeOK(0, 233, 233)
	expectedPayload := []byte{0, 233, 233, 0, 0, 0, 0}
	if len(payload) != len(expectedPayload) {
		t.Fatalf("bad result %v, expected: %v", len(payload), len(expectedPayload))
	}
	if !bytes.Equal(payload, expectedPayload) {
		t.Fatalf("bad result %v, expected: %v", payload, expectedPayload)
	}
}

func TestEncodeError(t *testing.T) {
	// test example from https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_err_packet.html
	c := setupPacketCoder()
	m := NewMySqlError(ER_NO_TABLES_USED, "No tables used")
	payload := c.encodeError(m)
	expectedPayload := []byte{
		0xff, 0x48, 0x04,
		0x23, 0x48, 0x59, 0x30, 0x30, 0x30, 0x4e, 0x6f, 0x20,
		0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x20, 0x75, 0x73, 0x65, 0x64,
	}

	if !bytes.Equal(payload, expectedPayload) {
		t.Fatalf("bad result %v, expected: %v", payload, expectedPayload)
	}

	if len(payload) != len(m.Message)+9 {
		t.Fatalf("bad result %v, expected: %v", len(payload), len(m.Message))
	}
}

func TestEncodeError2(t *testing.T) {
	// test example from https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_err_packet.html
	c := setupPacketCoder()
	payload := c.encodeError(fmt.Errorf("oh fuck"))
	expectedPayload := []byte{
		255, 81, 4, 35, 72, 89, 48, 48, 48, 111, 104, 32, 102, 117, 99, 107,
	}

	if !bytes.Equal(payload, expectedPayload) {
		t.Fatalf("bad result %v, expected: %v", payload, expectedPayload)
	}
}

func TestEncodeEOF(t *testing.T) {
	// test example from https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_eof_packet.html
	c := setupPacketCoder()
	payload := c.encodeEOF(0, 2)
	expectedPayload := []byte{
		254, 0, 0, 2, 0,
	}

	if !bytes.Equal(payload, expectedPayload) {
		t.Fatalf("bad result %v, expected: %v", payload, expectedPayload)
	}
}
