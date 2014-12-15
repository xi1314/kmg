package kmgNet

import (
	"errors"
	"io"
	"net"
	"time"
)

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
