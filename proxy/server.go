package proxy

import (
	"log"
	"net"
	"runtime"
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

func (s *Server) Run() {
	s.isRunning = true

	for s.isRunning {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal("accept error %s", err.Error())
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) Close() {
	s.isRunning = false

	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *Server) handleConn(conn net.Conn) {
	c := NewConnection(conn)

	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			log.Print("handleConn panic %v: %v\n%s", conn.RemoteAddr().String(), err, buf)
		}

		c.Close()
	}()

	c.Run()
}
