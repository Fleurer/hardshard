package mysql

import (
	"fmt"
	"net"
	"os"
	"sync/atomic"

	"github.com/siddontang/go-log/log"
)

var connectionIdCounter uint32 = 10000

type Connection struct {
	conn         net.Conn
	isClosed     bool
	packetIO     *PacketIO
	connectionId uint32
	capabilities uint32
	status       uint16
	salt         []byte
	collationId  uint8
}

type handkshakeResponse struct {
	capabilities   uint32
	charset        byte
	maxPacketSize  uint32
	user           []byte
	authData       []byte
	db             []byte
	authPluginName []byte
	attrs          map[string]string
}

func NewConnection(conn net.Conn) *Connection {
	c := &Connection{
		conn:         conn,
		packetIO:     NewPacketIOByConn(conn),
		isClosed:     false,
		connectionId: atomic.AddUint32(&connectionIdCounter, 1),
		capabilities: CLIENT_PLUGIN_AUTH | CLIENT_SECURE_CONNECTION | CLIENT_CONNECT_WITH_DB | CLIENT_CONNECT_ATTRS | CLIENT_PROTOCOL_41,
		status:       SERVER_STATUS_AUTOCOMMIT,
		salt:         GenerateRandBuf(20),
		collationId:  DEFAULT_COLLATION_ID,
	}
	return c
}

func (c *Connection) Run() {
	defer func() {
		c.Close()
	}()
	if err := c.handshake(); err != nil {
		c.writeError(err)
		return
	}
	c.loop()
}

func (c *Connection) handshake() error {
	if err := c.writeInitialHandshake(); err != nil {
		log.Error("handshake: writeInitialHandshake fail: err=%s", err)
		return err
	}
	handshake, err := c.readHandshakeResponse()
	fmt.Printf("client handshake: %v", handshake)
	if err != nil {
		log.Error("handshake: readHandshakeResponse fail: err=%s", err)
		return err
	}
	if err := c.writeOK(0, 0, 0); err != nil {
		log.Error("handshake: writeOK fail: err=%s", err)
		return err
	}
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
	case COM_QUIT:
	case COM_QUERY:
	case COM_INIT_DB:
	case COM_FIELD_LIST:
	case COM_STMT_PREPARE:
	case COM_STMT_EXECUTE:
	case COM_STMT_CLOSE:
	case COM_STMT_SEND_LONG_DATA:
	case COM_STMT_RESET:
	case COM_SET_OPTION:
	default:
		// return NewError()
	}
	return nil
}

func (c *Connection) writeOK(status uint16, affectedRows uint64, insertId uint64) error {
	// OK_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_ok_packet.html
	// As of MySQL 5.7.5, OK packes are also used to indicate EOF, and EOF packets are deprecated.
	// These rules distinguish whether the packet represents OK or EOF:
	// - OK: header = 0 and length of packet > 7
	// - EOF: header = 0xfe and length of packet < 9
	payload := make([]byte, 0, 32)
	payload = append(payload, OK_HEADER)
	payload = append(payload, EncodeLencInt(affectedRows)...)
	payload = append(payload, EncodeLencInt(insertId)...)
	if c.capabilities&CLIENT_PROTOCOL_41 > 0 {
		// if CLIENT_PROTOCOL_41 is set, the packet contains a warning count.
		payload = append(payload, EncodeUint16(status)...) // status_flags
		payload = append(payload, EncodeUint16(0)...)      // number of warnings
	}
	return c.packetIO.WritePacket(payload)
}

func (c *Connection) writeEOF(warnings uint16, status uint16) error {
	// EOF_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_eof_packet.html
	payload := make([]byte, 0, 5)
	payload = append(payload, EOF_HEADER)
	if c.capabilities&CLIENT_PROTOCOL_41 > 0 {
		payload = append(payload, EncodeUint16(warnings)...) // number of warnings
		payload = append(payload, EncodeUint16(status)...)   // SERVER_STATUS_flags_enum
	}
	return c.packetIO.WritePacket(payload)
}

func (c *Connection) writeError(e error) error {
	// ERR_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_err_packet.html
	// https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_dt_strings.html
	var m *MySqlError
	var ok bool

	m, ok = e.(*MySqlError)
	if !ok {
		m = NewMySqlError(ER_UNKNOWN_ERROR, e.Error())
	}

	payload := make([]byte, 0, 16+len(m.Message))
	payload = append(payload, ERR_HEADER)
	payload = append(payload, EncodeUint16(m.Code)...)

	if c.capabilities&CLIENT_PROTOCOL_41 > 0 {
		// It contains a SQL state value if CLIENT_PROTOCOL_41 is enabled.
		payload = append(payload, '#')
		payload = append(payload, m.State...) // FixedLengthString, always 5 bytes length
	}

	payload = append(payload, m.Message...) // RestOfPacketString
	return c.packetIO.WritePacket(payload)
}

func (c *Connection) writeInitialHandshake() error {
	// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html
	payload := make([]byte, 0, 128)
	// protocol version: 10
	payload = append(payload, 10)
	// string[NULL], server version
	payload = append(payload, SERVER_VERSION...)
	// int32, connection id
	payload = append(payload, EncodeUint32(c.connectionId)...)
	// string[8], auth-plugin-data, part-1
	payload = append(payload, c.salt[0:8]...)
	// 1 byte, [00] filter
	payload = append(payload, 0)
	// 2 bytes, capabilities lower 2 bytes
	payload = append(payload, byte(c.capabilities), byte(c.capabilities>>8))
	// 1 byte, charset, default utf-8
	payload = append(payload, byte(c.collationId))
	// 2 bytes, status
	payload = append(payload, EncodeUint16(c.status)...)
	// 2 bytes, capability_flags_2, upper 2 bytes of the capabilities
	payload = append(payload, byte(c.capabilities>>16), byte(c.capabilities>>24))
	// 1 byte, length of auth-plugin-data or 0
	if c.capabilities&CLIENT_PLUGIN_AUTH > 0 {
		payload = append(payload, byte(len(c.salt)))
	} else {
		panic("please CLIENT_PLUGIN_AUTH")
	}
	// 10 bytes, reserved
	payload = append(payload, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	// string[$len]   auth-plugin-data-part-2
	// $len=MAX(13, length of auth-plugin-data - 8)
	if c.capabilities&CLIENT_SECURE_CONNECTION > 0 {
		if len(c.salt[8:]) > 13 {
			panic("please len(salt[8:]) <= 13")
		}
		payload = append(payload, c.salt[8:]...)
	} else {
		panic("please CLIENT_SECURE_CONNECTION")
	}
	// string[NUL] auth-plugin name, if capabilities & CLIENT_PLUGIN_AUTH
	payload = append(payload, 0)
	return c.packetIO.WritePacket(payload)
}

// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeResponse
func (c *Connection) readHandshakeResponse() (*handkshakeResponse, error) {
	pr, err := c.packetIO.NewPacketReader()
	h := handkshakeResponse{}
	if err != nil {
		return nil, err
	}
	if h.capabilities, err = pr.ReadUint32(); err != nil {
		return nil, err
	}
	if h.maxPacketSize, err = pr.ReadUint32(); err != nil {
		return nil, err
	}
	if h.charset, err = pr.ReadByte(); err != nil {
		return nil, err
	}
	// reserved 23 bytes
	pr.Next(23)
	if h.user, err = pr.ReadBytes('\x00'); err != nil {
		return nil, err
	}
	if h.capabilities&CLIENT_SECURE_CONNECTION > 0 {
		n, err := pr.ReadByte()
		if err != nil {
			return nil, err
		}
		h.authData = pr.Next(int(n))
	}
	if h.capabilities&CLIENT_CONNECT_WITH_DB > 0 {
		fmt.Printf("readHandshakeResponse: user:%s\n", h.user)
		pr.PrintPacket()
		if h.db, err = pr.ReadBytes('\x00'); err != nil {
			return nil, err
		}
	}
	if h.capabilities&CLIENT_PLUGIN_AUTH > 0 {
		if h.authPluginName, err = pr.ReadBytes('\x00'); err != nil {
			return nil, err
		}
	}
	if h.capabilities&CLIENT_CONNECT_ATTRS > 0 {
		n, _, err := pr.ReadLencInt()
		if err != nil {
			return nil, err
		}
		for i := uint64(0); i < n; i++ {
			key, err := pr.ReadLencString()
			if err != nil {
				return nil, err
			}
			val, err := pr.ReadLencString()
			if err != nil {
				return nil, err
			}
			h.attrs[key] = val
		}
	}
	return &h, nil
}
