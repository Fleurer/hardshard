package mysql

// Protocol Basics: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_packets.html
// ERR_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_err_packet.html
// EOF_PACKET: https://dev.mysql.com/doc/dev/mysql-server/8.0.0/page_protocol_basic_eof_packet.html

type PacketCoder struct {
	capabilities uint32
}

func (c *PacketCoder) encodeOK(status uint, affectedRows uint64, insertId uint64) []byte {
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
		payload = append(payload, byte(status), byte(status>>8)) // status_flags
		payload = append(payload, 0, 0)                          // number of warnings
	}
	return payload
}
