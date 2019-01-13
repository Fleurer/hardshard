package mysql

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"testing"
)

func setupConnnection() (*Connection, net.Conn) {
	server, client := net.Pipe()
	conn := NewConnection(server)
	conn.salt = []byte("salt1salt2salt3salt4")
	return conn, client
}

func TestWriteOk(t *testing.T) {
	conn, client := setupConnnection()
	defer client.Close()
	go func() {
		conn.writeOK(5, 233, 234)
		conn.Close()
	}()
	buf, err := ioutil.ReadAll(client)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expectedBuf := []byte{7, 0, 0, 0, 0, 233, 234, 5, 0, 0, 0}
	if !bytes.Equal(buf, expectedBuf) {
		t.Fatalf("bad result: %v, expected: %v", buf, expectedBuf)
	}
}

func TestWriteEOF(t *testing.T) {
	conn, client := setupConnnection()
	defer client.Close()
	go func() {
		conn.writeEOF(123, 124)
		conn.Close()
	}()
	buf, err := ioutil.ReadAll(client)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expectedBuf := []byte{
		5, 0, 0, 0,
		254, 123, 0, 124, 0,
	}
	if !bytes.Equal(buf, expectedBuf) {
		t.Fatalf("bad result: %v, expected: %v", buf, expectedBuf)
	}
}

func TestWriteError(t *testing.T) {
	conn, client := setupConnnection()
	defer client.Close()
	go func() {
		m := NewMySqlError(ER_NO_TABLES_USED, "No tables used")
		conn.writeError(m)
		conn.Close()
	}()
	buf, err := ioutil.ReadAll(client)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expectedBuf := []byte{
		23, 0, 0, 0,
		0xff, 0x48, 0x04,
		0x23, 0x48, 0x59, 0x30, 0x30, 0x30, 0x4e, 0x6f, 0x20,
		0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x20, 0x75, 0x73, 0x65, 0x64,
	}
	if !bytes.Equal(buf, expectedBuf) {
		t.Fatalf("bad result: %v, expected: %v", buf, expectedBuf)
	}
}

func TestWriteError2(t *testing.T) {
	conn, client := setupConnnection()
	defer client.Close()
	go func() {
		conn.writeError(fmt.Errorf("oh fuck"))
		conn.Close()
	}()
	buf, err := ioutil.ReadAll(client)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expectedBuf := []byte{
		16, 0, 0, 0,
		255, 81, 4, 35, 72, 89, 48, 48, 48, 111, 104, 32, 102, 117, 99, 107,
	}
	if !bytes.Equal(buf, expectedBuf) {
		t.Fatalf("bad result: %v, expected: %v", buf, expectedBuf)
	}
}

func TestWriteInitialHandshake(t *testing.T) {
	conn, client := setupConnnection()
	defer client.Close()
	go func() {
		conn.writeInitialHandshake()
		conn.Close()
	}()
	buf, err := ioutil.ReadAll(client)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expectedBuf := []byte{
		66, 0, 0, 0, 10, 53, 46, 53,
		46, 51, 49, 45, 104, 97, 114, 100,
		115, 104, 97, 114, 100, 45, 48, 46,
		49, 0, 21, 39, 0, 0, 115, 97,
		108, 116, 49, 115, 97, 108, 0, 8,
		130, 33, 2, 0, 24, 0, 20, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 116, 50, 115, 97, 108, 116, 51,
		115, 97, 108, 116, 52, 0,
	}
	if !bytes.Equal(buf, expectedBuf) {
		t.Fatalf("bad result: %v, expected: %v", buf, expectedBuf)
	}
}

func TestReadHandshakeResponse(t *testing.T) {
	conn, client := setupConnnection()
	defer client.Close()
	go func() {
		buf := []byte{
			0x79, 0x00, 0x00, 0x00, 0x0D, 0xA2, 0x3A, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x2D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x75, 0x75, 0x75, 0x75, 0x75, 0x00, 0x14, 0xAF, 0x62, 0xA6, 0x6F, 0x39, 0xA8, 0x43, 0x33, 0xC4, 0x0C, 0xD9, 0x70, 0x48, 0xAC, 0x1D, 0x2C, 0x8E, 0xB0, 0xA6, 0x77, 0x64, 0x62, 0x32, 0x33, 0x33, 0x00, 0x00, 0x36, 0x0F, 0x5F, 0x63, 0x6C, 0x69, 0x65, 0x6E, 0x74, 0x5F, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6F, 0x6E, 0x05, 0x30, 0x2E, 0x39, 0x2E, 0x33, 0x04, 0x5F, 0x70, 0x69, 0x64, 0x05, 0x34, 0x38, 0x32, 0x35, 0x33, 0x0C, 0x5F, 0x63, 0x6C, 0x69, 0x65, 0x6E, 0x74, 0x5F, 0x6E, 0x61, 0x6D, 0x65, 0x07, 0x70, 0x79, 0x6D, 0x79, 0x73, 0x71, 0x6C,
		}
		client.Write(buf)
		client.Close()
	}()
	h, err := conn.readHandshakeResponse()
	if err != nil {
		t.Fatalf("readHandshakeResponse err: %s", err)
	}
	expectedUser := []byte("uuuuu\x00")
	if !bytes.Equal(h.user, expectedUser) {
		t.Fatalf("h.user mismatch, expected: %s, got: %s", expectedUser, h.user)
	}
	conn.Close()
}
