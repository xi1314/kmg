package kmgNet

import (
	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgIo"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"sync"
	"time"
)

func (conn *fasterTcpConn) Read(b []byte) (nr int, err error) {
	err = conn.TCPConn.SetReadDeadline(time.Now().Add(10 * time.Minute))
	if err != nil {
		return 0, err
	}
	nr, err = conn.TCPConn.Read(b)
	if err != nil {
		return nr, err
	}
	conn.closeLock.Lock()
	defer conn.closeLock.Unlock()
	fdNum := int(conn.fd.Fd())
	if fdNum == -1 { //已经关闭过了,此处不管了.
		return nr, err
	}
	err = unix.SetsockoptInt(fdNum, unix.IPPROTO_TCP, unix.TCP_QUICKACK, 1)
	if err != nil { //TODO 此处总是会爆bad file descriptor,原因不明.
		kmgErr.LogErrorWithStack(err)
		return
	}
	return
}
func (conn *fasterTcpConn) Write(b []byte) (nr int, err error) {
	err = conn.TCPConn.SetWriteDeadline(time.Now().Add(10 * time.Minute))
	if err != nil {
		return 0, err
	}
	return conn.TCPConn.Write(b)
}

func (conn *fasterTcpConn) Close() (err error) {
	conn.closeLock.Lock()
	defer conn.closeLock.Unlock()
	//此处可以添加本连接是否关闭的信息,然后给unix.SetsockoptInt使用.
	return kmgIo.MultiErrorHandle(conn.TCPConn.CloseRead, conn.TCPConn.Close, conn.fd.Close)
}

func (conn *fasterTcpConn) GetUnderlyingConn() net.Conn {
	return conn.TCPConn
}

type fasterTcpConn struct {
	*net.TCPConn
	fd        *os.File
	closeLock sync.Mutex
}

func LessDelayTcpConn(conn *net.TCPConn) (connOut net.Conn, err error) {
	//err = conn.SetKeepAlive(true)
	//if err!=nil{
	//	kmgErr.LogErrorWithStack(err)
	//	return nil,err
	//}
	//err = conn.SetKeepAlivePeriod(5*time.Second) //5s太小,耗流量非常凶残.
	//if err!=nil{
	//	kmgErr.LogErrorWithStack(err)
	//	return nil,err
	//}
	fd, err := conn.File()
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		return
	}
	conn1, err := net.FileConn(fd)
	if err != nil {
		fd.Close()
		kmgErr.LogErrorWithStack(err)
		return
	}
	conn.Close()
	//尝试将连接重新设置回 block 模式,减少cpu占用,此方案不稳定,并且不知道如何解决不稳定的问题.
	//err = unix.SetNonblock(int(fd.Fd()),true)
	//if err!=nil{
	//	fd.Close()
	//	kmgErr.LogErrorWithStack(err)
	//	return nil,err
	//}
	//return NewDebugConn(fasterTcpConn{TCPConn: conn, fd: fd},conn.LocalAddr().String()+"_"+conn.RemoteAddr().String()), nil
	return &fasterTcpConn{TCPConn: conn1.(*net.TCPConn), fd: fd}, nil
}

func LessDelayDial(network string, address string) (conn net.Conn, err error) {
	conn, err = net.Dial(network, address)
	if err != nil {
		return conn, err
	}
	return LessDelayTcpConn(conn.(*net.TCPConn))
}

type fasterTcpListener struct {
	*net.TCPListener
}

func MustLessDelayListen(network string, address string) net.Listener {
	return fasterTcpListener{MustListen(network, address).(*net.TCPListener)}
}

func (l fasterTcpListener) Accept() (outC net.Conn, err error) {
	c, err := l.TCPListener.AcceptTCP()
	if err != nil {
		return c, err
	}
	return LessDelayTcpConn(c)
}
