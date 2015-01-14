package kmgNet

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"sync"
)

// dial tcp with tcp fastopen
// you should use echo 3 > /proc/sys/net/ipv4/tcp_fastopen to enable it
// you should use linux kernel version >3.7
// you should write something before read to use this function.
// network is useless, it always use tcp4
func TfoLazyDial(network string, nextAddr string) (conn net.Conn, err error) {
	return &tfoLazyConn(network, nextAddr), nil
}

type tfoLazyConn struct {
	net.Conn
	nextAddr string
	dialLock sync.Mutex
}

func (c *tfoLazyConn) Read(b []byte) (n int, err error) {
	//fast path
	if c.Conn != nil {
		return c.Conn.Read(b)
	}
	c.dialLock.Lock()
	defer c.dialLock.Unlock()
	if c.Conn == nil {
		c.Conn, err = net.Dial("tcp", c.nextAddr)
		if err != nil {
			return
		}
	}
	return c.Conn.Read(b)
}

func (c *tfoLazyConn) Write(b []byte) (n int, err error) {
	//fast path
	if c.Conn != nil {
		return c.Conn.Write(b)
	}
	c.dialLock.Lock()
	defer c.dialLock.Unlock()
	if c.Conn != nil {
		return c.Conn.Write(b)
	}
	c.Conn, err = TfoDial(c.nextAddr, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

//dial tcp with tcp fastopen
func TfoDial(nextAddr string, firstData []byte) (conn net.Conn, err error) {
	s, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(s)
	tcpAddr, err := net.ResolveTCPAddr("tcp", nextAddr)
	if err != nil {
		return nil, err
	}
	ip := tcpAddr.IP.To4()
	if ip == nil {
		return nil, fmt.Errorf("[TProxyListen] only support tcp4 right now.")
	}
	var ipA [4]byte
	copy(ipA[:], ip[:4])
	sockAddr := &unix.SockaddrInet4{
		Port: tcpAddr.Port,
		Addr: ipA,
	}
	err = unix.Sendto(s, firstData, unix.MSG_FASTOPEN, sockAddr)
	if err != nil {
		return
	}
	f := os.NewFile(uintptr(s), "TFODial")
	defer f.Close()
	return net.FileConn(f)
}
