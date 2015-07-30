package kmgHttp

import (
	"fmt"
	"net"
	"net/http"

	"github.com/bronze1man/kmg/kmgNet"
)

//一个http的满足 kmgNet.Server接口的服务器
func NewHttpNetServer(Addr string, handler http.Handler) kmgNet.Server {
	return &httpNetServer{
		Server: &http.Server{
			Handler: handler,
		},
		addr: Addr,
	}
}

type httpNetServer struct {
	*http.Server
	addr     string
	listener net.Listener
}

func (s *httpNetServer) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	go func() {
		err := s.Server.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

func (s *httpNetServer) Close() error {
	return s.listener.Close()
}

func (s *httpNetServer) Addr() (net.Addr, error) {
	if s.listener == nil {
		return nil, fmt.Errorf("[httpNetServer] you should start server first")
	}
	return s.listener.Addr(), nil
}
