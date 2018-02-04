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
	cmd = payload[0]
	body = payload[1:]

	switch cmd {
	case mysql.COM_QUIT:
	case mysql.COM_QUERY:
	case mysql.COM_INIT_DB:
	case mysql.COM_FIELD_LIST:
	case mysql.COM_STMT_PREPARE:
	case mysql.COM_STMT_EXECUTE:
	case mysql.COM_STMT_CLOSE:
	case mysql.COM_STMT_SEND_LONG_DATA:
	case mysql.COM_STMT_RESET:
	case mysql.COM_SET_OPTION:
	default:
		return mysql.NewError()
	}
}
