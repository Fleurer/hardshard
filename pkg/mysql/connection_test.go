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
	return NewConnection(server), client
}

func TestWriteOk(t *testing.T) {
	fmt.Printf("blah!")
	conn, client := setupConnnection()
	defer conn.Close()
	defer client.Close()
	conn.writeOK(0, 233, 234)
	buf, err := ioutil.ReadAll(client)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expectedBuf := []byte{0, 233, 233, 0, 0, 0, 0}
	if !bytes.Equal(buf, expectedBuf) {
		t.Fatalf("bad result: %s, expected: %s", buf, expectedBuf)
	}
}
