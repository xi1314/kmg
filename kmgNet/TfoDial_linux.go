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
	if c.isClosed {
		c.lock.Unlock()
		return 0, ErrClosing
	}
	if len(b) <= tfoFirstSize {
		c.Conn, err = TfoDial(c.nextAddr, b)
		if err != nil {
			c.lock.Unlock()
			return 0, err
		}
		c.lock.Unlock()
		return len(b), err
	} else {
		c.Conn, err = TfoDial(c.nextAddr, b[:tfoFirstSize])
		if err != nil {
			c.lock.Unlock()
			return 0, err
		}
		c.lock.Unlock()
		_, err = c.Conn.Write(b[tfoFirstSize:])
		if err != nil {
			return 0, err
		}
		return len(b), err
	}
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

const tfoFirstSize = 1000

// dial tcp with tcp fastopen
// 第一个包体积不要太大,需要小于一定数量,否则会被吃掉(正确性问题),
// 如果过大,此处会在连接时发送前一部分,连接后又发送一部分
func TfoDial(nextAddr string, data []byte) (conn net.Conn, err error) {
	s, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(s)
	sa, err := IPv4TcpAddrToUnixSocksAddr(nextAddr)
	if err != nil {
		return nil, err
	}
	if len(data) <= tfoFirstSize {
		err = unix.Sendto(s, data, unix.MSG_FASTOPEN, sa)
		if err != nil {
			return
		}
	} else {
		err = unix.Sendto(s, data[:tfoFirstSize], unix.MSG_FASTOPEN, sa)
		if err != nil {
			return
		}
	}
	f := os.NewFile(uintptr(s), "TFODial")
	defer f.Close()
	conn, err = net.FileConn(f)
	if err != nil {
		return
	}
	if len(data) > tfoFirstSize {
		_, err = conn.Write(data[tfoFirstSize:])
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}
