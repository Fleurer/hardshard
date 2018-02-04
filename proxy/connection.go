package proxy

import (
	"fmt"
	"net"

	"github.com/Fleurer/hardshard/mysql"
	"github.com/siddontang/go-log/log"
)

type Connection struct {
	conn     net.Conn
	packetio *mysql.PacketIO

	isClosed bool
}

func NewConnection(conn net.Conn) *Connection {
	c := &Connection{}
	c.conn = conn
	c.packetio = mysql.NewPacketIOByConn(conn)
	c.isClosed = false
	return c
}

func (c *Connection) Run() {
	for {
		payload, err := c.packetio.ReadPacket()
		if err != nil {
			log.Warn("connection.Run() readPacket error: %s", err.Error())
			return
		}

		err = c.handlePacket(payload)
		if err != nil {
			fmt.Errorf("handlePacket error: %s", err.Error())
			c.packetio.WriteErrorPacket(err)
		}

		if c.isClosed {
			return
		}
	}
}

func (c *Connection) Close() error {
	if c.isClosed {
		return nil
	}

	err := c.conn.Close()
	if err != nil {
		return err
	}

	c.isClosed = true
	return nil
}

func (c *Connection) handlePacket(payload []byte) error {
	return nil
}
