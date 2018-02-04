package mysql

import "io"

type ResponsePacket interface {
	WritePayload(w io.Writer) error
}

type OkPacket struct {
}
