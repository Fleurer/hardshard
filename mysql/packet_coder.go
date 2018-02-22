package mysql

import "encoding/binary"

type PacketCoder struct {
	capabilities uint32
}

type HandshakeResponse struct {
	capabilityFlags uint32
	maxPacketSize   uint32
	charsetId       uint8
	userName        string
	authResponse    string
	databaseName    string
	connectAttrs    map[string]string
}

func (c *PacketCoder) EncodeOK(status uint16, affectedRows uint64, insertId uint64) []byte {
	// OK_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_ok_packet.html
	// As of MySQL 5.7.5, OK packes are also used to indicate EOF, and EOF packets are deprecated.
	// These rules distinguish whether the packet represents OK or EOF:
	// - OK: header = 0 and length of packet > 7
	// - EOF: header = 0xfe and length of packet < 9
	payload := make([]byte, 0, 32)
	payload = append(payload, OK_HEADER)
	payload = append(payload, PutLengthEncodedInt(affectedRows)...)
	payload = append(payload, PutLengthEncodedInt(insertId)...)

	if c.capabilities&CLIENT_PROTOCOL_41 > 0 {
		// if CLIENT_PROTOCOL_41 is set, the packet contains a warning count.
		payload = append(payload, Uint16ToBytes(status)...) // status_flags
		payload = append(payload, Uint16ToBytes(0)...)      // number of warnings
	}
	return payload
}

func (c *PacketCoder) EncodeEOF(warnings uint16, statusFlags uint16) []byte {
	// EOF_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_eof_packet.html
	payload := make([]byte, 0, 5)
	payload = append(payload, EOF_HEADER)
	if c.capabilities&CLIENT_PROTOCOL_41 > 0 {
		payload = append(payload, Uint16ToBytes(warnings)...)    // number of warnings
		payload = append(payload, Uint16ToBytes(statusFlags)...) // SERVER_STATUS_flags_enum
	}
	return payload
}

func (c *PacketCoder) EncodeError(e error) []byte {
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
	payload = append(payload, Uint16ToBytes(m.Code)...)

	if c.capabilities&CLIENT_PROTOCOL_41 > 0 {
		// It contains a SQL state value if CLIENT_PROTOCOL_41 is enabled.
		payload = append(payload, '#')
		payload = append(payload, m.State...) // FixedLengthString, always 5 bytes length
	}

	payload = append(payload, m.Message...) // RestOfPacketString
	return payload
}

func (c *PacketCoder) EncodeInitialHandshake(capabilityFlags uint32, statusFlags uint16, connectionId uint32, collationId uint8, authPluginData []byte) []byte {
	// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html
	payload := make([]byte, 0, 128)
	// protocol version: 10
	payload = append(payload, 10)
	// string[NULL], server version
	payload = append(payload, SERVER_VERSION...)
	// int32, connection id
	payload = append(payload, Uint32ToBytes(connectionId)...)
	// string[8], auth-plugin-data, part-1
	payload = append(payload, authPluginData[0:8]...)
	// 1 byte, [00] filter
	payload = append(payload, 0)
	// 2 bytes, capabilities lower 2 bytes
	payload = append(payload, byte(capabilityFlags), byte(capabilityFlags>>8))
	// 1 byte, charset, default utf-8
	payload = append(payload, byte(collationId))
	// 2 bytes, statusFlags
	payload = append(payload, Uint16ToBytes(statusFlags)...)
	// 2 bytes, capability_flags_2, upper 2 bytes of the capabilityFlags
	payload = append(payload, byte(capabilityFlags>>16), byte(capabilityFlags>>24))
	// 1 byte, length of auth-plugin-data or 0
	if capabilityFlags&CLIENT_PLUGIN_AUTH > 0 {
		payload = append(payload, byte(len(authPluginData)))
	} else {
		panic("please CLIENT_PLUGIN_AUTH")
	}
	// 10 bytes, reserved
	payload = append(payload, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	// string[$len]   auth-plugin-data-part-2
	// $len=MAX(13, length of auth-plugin-data - 8)
	if capabilityFlags&CLIENT_SECURE_CONNECTION > 0 {
		l := len(authPluginData) - 8
		if l < 13 {
			l = 13
		}
		payload = append(payload, authPluginData[8:]...)
	} else {
		panic("please CLIENT_SECURE_CONNECTION")
	}
	// string[NUL] auth-plugin name, if capabilities & CLIENT_PLUGIN_AUTH
	payload = append(payload, 0)
	return payload
}

func (c *PacketCoder) DecodeHandshakeResponse(payload []byte) (*HandshakeResponse, error) {
	// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeResponse
	pos := 0
	capabilityFlags := binary.LittleEndian.Uint32(payload[pos : pos+4])
	pos += 4
	maxPacketSize := binary.LittleEndian.Uint32(payload[4:8])
	charsetId := uint8(payload[8])
	return nil
}
