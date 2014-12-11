package kmgNet

import (
	"bytes"
	"errors"
	"net"
	"time"
)

func BufioConn() net.Conn {
	return bufioConn{
		Buffer: &bytes.Buffer{},
	}
}

type bufioConn struct {
	*bytes.Buffer
}

func (c bufioConn) Close() error {
	return nil
}

func (c bufioConn) LocalAddr() net.Addr {
	return bufioAddr{}
}

func (c bufioConn) RemoteAddr() net.Addr {
	return bufioAddr{}
}

func (c bufioConn) SetDeadline(t time.Time) error {
	return errors.New("kmgNet.BufioConn does not support deadlines")
}

func (c bufioConn) SetReadDeadline(t time.Time) error {
	return errors.New("kmgNet.BufioConn does not support deadlines")
}

func (c bufioConn) SetWriteDeadline(t time.Time) error {
	return errors.New("kmgNet.BufioConn does not support deadlines")
}

type bufioAddr struct{}

func (a bufioAddr) Network() string {
	return "bufioAddr"
}

func (a bufioAddr) String() string {
	return "bufioAddr"
}
