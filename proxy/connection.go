package proxy

import "net"

type Connection struct {
	conn     *net.Conn
	packetio *PacketIO

	isClosed bool
}

func NewConnection(conn *net.Conn) *Connection {
	c := &Connection{}
	c.conn = conn
	c.isClosed = false
	return conn
}

func (c *Connection) Run() {
	for {
		packet, err := c.packetio.readPacket()
		if err != nil {
			log.Warn("connection.Run() readPacket error: %s", err.Error())
			return
		}

		err := c.handlePacket(packet)
		if err != nil {
			log.Error("dispatch error %s", err.Error())
			c.packetio.writeErrorPacket(err)
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
