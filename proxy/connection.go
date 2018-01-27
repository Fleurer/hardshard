package proxy

import (
	"fmt"
	"net"

	"github.com/Fleurer/hardshard/mysql"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("connection")

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
		packet, err := c.packetio.ReadPacket()
		if err != nil {
			log.Warnf("connection.Run() readPacket error: %s", err.Error())
			return
		}

		err = c.handlePacket(packet)
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

func (c *Connection) handlePacket(packet []byte) error {
	return nil
}
