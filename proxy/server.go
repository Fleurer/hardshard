package proxy

import (
	"net"
	"runtime"

	"github.com/Fleurer/hardshard/connection"
)

type Server struct {
	host     string
	port     int
	listener net.Listener

	isRunning bool
}

func NewServer(host string, port int) *Server {
	s := &Server{}
	s.host = host
	s.port = port
	s.isRunning = false
	return s
}

func (s *Server) Run() error {
	s.isRunning = true

	for s.isRunning {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Error("accept error %s", err.Error())
			continue
		}

		go s.handleConn(conn)
	}

	return s
}

func (s *Server) Close() {
	s.isRunning = false

	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *Server) handleConn(conn *net.Conn) {
	conn := connection.NewConnection(c)

	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			log.Error("handleConn panic %v: %v\n%s", c.RemoteAddr().String(), err, buf)
		}

		conn.Close()
	}()

	conn.Run()
}
