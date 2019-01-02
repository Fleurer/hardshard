package proxy

import (
	"fmt"
	"net"
	"os"

	"github.com/Fleurer/hardshard/pkg/mysql"
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
	if err := c.handshake(); err != nil {
		c.writeError(err)
		c.Close()
		return
	}

	c.loop()
}

func (c *Connection) handshake() error {
	return nil
}

func (c *Connection) loop() {
	for {
		payload, err := c.packetio.ReadPacket()
		if err != nil {
			log.Warn("connection.Run() readPacket error: %s", err.Error())
			return
		}

		err = c.handleRequestPacket(payload)
		if err != nil {
			fmt.Errorf("handleRequestPacket error: %s", err.Error())
			// c.packetio.WriteErrorPacket(err)
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

func (c *Connection) handleRequestPacket(payload []byte) error {
	cmd := payload[0]
	body := payload[1:]
	fmt.Printf("cmd: %v body: %v", cmd, body)
	os.Exit(1)

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
		// return mysql.NewError()
	}
	return nil
}

func (c *Connection) writeError(error) error {
	return nil
}
