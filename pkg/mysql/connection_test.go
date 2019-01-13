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
