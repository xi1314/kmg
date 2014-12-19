package kmgNet

import (
	"net"
	"strconv"
)

type Server interface {
	//异步开启服务器
	Start() error
	//关闭服务器,通常需要等一段时间(1ms)来等服务器确实关闭了
	Close() error
	//监听地址,如果可能请返回net库里面有的Addr,如果还没有开始监听会panic
	Addr() (net.Addr, error)
}

func MustGetServerAddrString(s Server) string {
	addr, err := s.Addr()
	if err != nil {
		panic(err)
	}
	return addr.String()
}

func MustGetServerLocalAddrString(s Server) string {
	addr, err := s.Addr()
	if err != nil {
		panic(err)
	}
	port, err := PortFromNetAddr(addr)
	if err != nil {
		panic(err)
	}
	return "127.0.0.1:" + strconv.Itoa(port)
}

func MustServerStart(s Server) {
	err := s.Start()
	if err != nil {
		panic(err)
	}
}

type FuncServer struct {
	StartFunc func() error
	CloseFunc func() error
	AddrFunc  func() (net.Addr, error)
	ExistAddr net.Addr
}

func (s *FuncServer) Start() error {
	return s.StartFunc()
}

//关闭服务器,通常需要等一段时间(1ms)来等服务器确实关闭了
func (s *FuncServer) Close() error {
	return s.CloseFunc()
}

//监听地址,如果可能请返回net库里面有的Addr,如果还没有开始监听会panic
func (s *FuncServer) Addr() (net.Addr, error) {
	if s.ExistAddr != nil {
		return s.ExistAddr, nil
	}
	return s.AddrFunc()
}
