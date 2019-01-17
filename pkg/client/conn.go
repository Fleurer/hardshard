package client

import "net"

type Conn struct {
	conn net.Conn

	addr     string
	user     string
	password string
	db       string

	capability uint32
	status     uint16
	collation  CollationId
	salt       []byte
}

func (c *Conn) Connect(addr string, user string, password string, db string) error {
	return nil
}

func (c *Conn) readInitialHandshake() error {
	return nil
}
