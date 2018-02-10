package mysql

import (
	"errors"
	"fmt"
)

var (
	ErrBadConn       = errors.New("connection was bad")
	ErrMalformPacket = errors.New("Malform packet error")
)

type MySqlError struct {
	Code    uint16
	Message string
	State   string
}

func (e *MySqlError) Error() string {
	return fmt.Sprintf("MYSQL ERROR %d (%s): %s", e.Code, e.State, e.Message)
}

func NewError(code uint16, message string) MySqlError {
	e := MySqlError{}
	e.Code = code
	e.Message = message
	if s, ok := MySQLState[code]; ok {
		e.State = s
	} else {
		e.State = DEFAULT_MYSQL_STATE
	}
	return e
}
