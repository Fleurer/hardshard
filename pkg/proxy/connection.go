package proxy

import (
	"fmt"
	"net"
	"os"
	"sync/atomic"

	"github.com/Fleurer/hardshard/pkg/mysql"
	"github.com/siddontang/go-log/log"
)

var connectionIdCounter uint32 = 10000

type Connection struct {
	conn         net.Conn
	packetIO     *mysql.PacketIO
	packetCoder  *mysql.PacketCoder
	isClosed     bool
	connectionId uint32
	capability   uint32
}

func NewConnection(conn net.Conn) *Connection {
	c := &Connection{}
	c.conn = conn
	c.packetIO = mysql.NewPacketIOByConn(conn)
	c.packetCoder = mysql.NewPacketCoder()
	c.isClosed = false
	c.connectionId = atomic.AddUint32(&connectionIdCounter, 1)
	c.capability = 0
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
		payload, err := c.packetIO.ReadPacket()
		if err != nil {
			log.Warn("connection.Run() readPacket error=%s", err.Error())
			return
		}

		err = c.handleRequestPacket(payload)
		if err != nil {
			log.Warn("handleRequestPacket error=%s", err.Error())
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

func (c *Connection) writeInitialHandshake() error {
	h := mysql.InitialHandshake{}
	payload := c.packetCoder.EncodeInitialHandshake(h)
	return c.packetIO.WritePacket(payload)
}

func (c *Connection) readHandshakeResponse() error {
	payload, err := c.packetIO.ReadPacket()
	if err != nil {
		return err
	}
}
