package kmgNet

import (
	"errors"
	"io"
	"net"
	"time"
)

func NewRwcOverConn(rwc io.ReadWriteCloser, conn net.Conn) net.Conn {
	return &RwcOverConn{
		Reader: rwc,
		Writer: rwc,
		Closer: rwc,
		Conn:   conn,
	}
}

type RwcOverConn struct {
	io.Reader
	io.Writer
	io.Closer
	net.Conn
}

func (c *RwcOverConn) Read(p []byte) (n int, err error) {
	return c.Reader.Read(p)
}

func (c *RwcOverConn) Write(p []byte) (n int, err error) {
	return c.Writer.Write(p)
}
func (c *RwcOverConn) Close() (err error) {
	return c.Closer.Close()
}
func (c *RwcOverConn) GetUnderlyingConn() net.Conn {
	return c.Conn
}

func RwcConn(rwc io.ReadWriteCloser) net.Conn {
	return rwcConn{
		ReadWriteCloser: rwc,
	}
}

type rwcConn struct {
	io.ReadWriteCloser
}

func (c rwcConn) LocalAddr() net.Addr {
	return FakeAddr
}

func (c rwcConn) RemoteAddr() net.Addr {
	return FakeAddr
}

func (c rwcConn) SetDeadline(t time.Time) error {
	return errors.New("kmgNet.rwcConn does not support deadlines")
}

func (c rwcConn) SetReadDeadline(t time.Time) error {
	return errors.New("kmgNet.rwcConn does not support deadlines")
}

func (c rwcConn) SetWriteDeadline(t time.Time) error {
	return errors.New("kmgNet.rwcConn does not support deadlines")
}
