package mysql

import "errors"

var (
	ErrBadConn       = errors.New("connection was bad")
	ErrMalformPacket = errors.New("Malform packet error")
)
