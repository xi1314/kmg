package kmgNet

import (
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
	return &tfoLazyConn{nextAddr: nextAddr}, nil
}

type tfoLazyConn struct {
	net.Conn
	nextAddr string
	lock     sync.Mutex
	isClosed bool
}

func (c *tfoLazyConn) Read(b []byte) (n int, err error) {
	c.lock.Lock()
	if c.Conn != nil && !c.isClosed {
		c.lock.Unlock()
		return c.Conn.Read(b)
	}
	if c.isClosed {
		c.lock.Unlock()
		return 0, ErrClosing
	}
	c.Conn, err = net.Dial("tcp", c.nextAddr)
	if err != nil {
		c.lock.Unlock()
		return
	}
	c.lock.Unlock()
	return c.Conn.Read(b)
}

func (c *tfoLazyConn) Write(b []byte) (n int, err error) {
	c.lock.Lock()
	if c.Conn != nil && !c.isClosed {
		c.lock.Unlock()
		return c.Conn.Write(b)
	}
	defer c.lock.Unlock()
	if c.isClosed {
		return 0, ErrClosing
	}
	c.Conn, err = TfoDial(c.nextAddr, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *tfoLazyConn) Close() error {
	c.lock.Lock()
	if c.isClosed {
		c.lock.Unlock()
		return ErrClosing
	}
	c.isClosed = true
	if c.Conn != nil {
		c.lock.Unlock()
		return c.Conn.Close()
	}
	c.lock.Unlock()
	return nil
}

//dial tcp with tcp fastopen
func TfoDial(nextAddr string, firstData []byte) (conn net.Conn, err error) {
	s, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(s)
	sa, err := IPv4TcpAddrToUnixSocksAddr(nextAddr)
	if err != nil {
		return nil, err
	}
	err = unix.Sendto(s, firstData, unix.MSG_FASTOPEN, sa)
	if err != nil {
		return
	}
	f := os.NewFile(uintptr(s), "TFODial")
	defer f.Close()
	return net.FileConn(f)
}
