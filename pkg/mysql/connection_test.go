package mysql

import (
	"bytes"
	"io/ioutil"
	"net"
	"testing"
)

func setupConnnection() (*Connection, net.Conn) {
	server, client := net.Pipe()
	return NewConnection(server), client
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
