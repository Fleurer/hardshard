package mysql

import "errors"

const (
	ErrBadConn       = errors.New("connection was bad")
	ErrMalformPacket = errors.New("Malform packet error")
)
