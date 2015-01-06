package kmgNet

import (
	"net"
	"time"
)

//拨号的时候有一个timeout,每次读和写的时候也有一个timeout
func NewTimeoutDialer(timeout time.Duration) func(network, addr string) (net.Conn, error) {
	return func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, timeout)
		if err != nil {
			return nil, err
		}
		return &timeoutConn{
			Conn:    conn,
			Timeout: timeout,
		}, nil
	}
}

func TimeoutConn(conn net.Conn, timeout time.Duration) net.Conn {
	return &timeoutConn{
		Conn:    conn,
		Timeout: timeout,
	}
}

type timeoutConn struct {
	net.Conn
	Timeout time.Duration
}

func (c *timeoutConn) Read(b []byte) (n int, err error) {
	c.Conn.SetReadDeadline(time.Now().Add(c.Timeout))
	return c.Conn.Read(b)
}

func (c *timeoutConn) Write(b []byte) (n int, err error) {
	c.Conn.SetWriteDeadline(time.Now().Add(c.Timeout))
	return c.Conn.Write(b)
}
