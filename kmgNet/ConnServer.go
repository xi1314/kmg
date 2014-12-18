package kmgNet

import (
	"io"
	"net"
	"strings"
)

type ConnHandler interface {
	ConnHandle(conn net.Conn)
}
type ConnHandlerFunc func(conn net.Conn)

func (f ConnHandlerFunc) ConnHandle(conn net.Conn) {
	f(conn)
}

type ConnServer struct {
	Listener net.Listener
	Handler  ConnHandler
	Closer   io.Closer
}

func (server *ConnServer) Close() (err error) {
	err = server.Listener.Close()
	var err1 error
	if server.Closer != nil {
		err1 = server.Closer.Close()
	}
	if err != nil {
		return err
	}
	return err1
}

//start之前已经监听过了
func (server *ConnServer) Start() (err error) {
	go func() {
		defer server.Listener.Close()
		for {
			conn, err := server.Listener.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") ||
					strings.Contains(err.Error(), "accept on closed mux") {
					break
				}
				panic(err)
			}
			go server.Handler.ConnHandle(conn)
		}
	}()
	return nil
}

func (server *ConnServer) Addr() (net.Addr, error) {
	return server.Listener.Addr(), nil
}

//这一步会开始监听
func NewTCPServer(listenAddr string, hander ConnHandler, closer io.Closer) (s *ConnServer, err error) {
	s:=&ConnServer{}
	s.Listener, err = net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	s.Handler = hander
	s.Closer = closer
	return s, nil
}
