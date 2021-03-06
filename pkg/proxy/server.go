package proxy

import (
	"net"
	"runtime"

	"github.com/Fleurer/hardshard/pkg/mysql"
	"github.com/siddontang/go-log/log"
)

type Server struct {
	addr     string
	listener net.Listener

	isRunning bool
}

func NewServer(addr string) (*Server, error) {
	s := &Server{}
	s.addr = addr

	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s.isRunning = false
	return s, nil
}

func (s *Server) Run() {
	s.isRunning = true

	for s.isRunning {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Error("accept error %s", err.Error())
			continue
		}
		log.Info("Accept %s", conn.RemoteAddr())

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
	myconn := mysql.NewConnection(conn)

	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			log.Error("handleConn panic %v: %v\n%s", conn.RemoteAddr().String(), err, buf)
		}

		myconn.Close()
	}()

	myconn.Run()
}
