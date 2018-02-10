package mysql

type PacketCoder struct {
	capabilities uint32
}

func (c *PacketCoder) encodeOK(status uint16, affectedRows uint64, insertId uint64) []byte {
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

func (c *PacketCoder) encodeEOF(warnings uint16, statusFlags uint16) []byte {
	// EOF_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_eof_packet.html
	payload := make([]byte, 0, 5)
	payload = append(payload, EOF_HEADER)
	if c.capabilities&CLIENT_PROTOCOL_41 > 0 {
		payload = append(payload, Uint16ToBytes(warnings)...)    // number of warnings
		payload = append(payload, Uint16ToBytes(statusFlags)...) // SERVER_STATUS_flags_enum
	}
	return payload
}

func (c *PacketCoder) encodeError(e error) []byte {
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
