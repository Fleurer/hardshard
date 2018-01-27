package mysql

import (
	"bytes"
	"testing"

	"github.com/Fleurer/hardshard/mysql"
)

func TestNewPacketIO(t *testing.T) {
	r := bytes.NewBufferString("hello\n")
	w := bytes.NewBufferString("")

	mysql.NewPacketIO(r, w)
}
