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
	return FakeAddr
}

func (c bufioConn) RemoteAddr() net.Addr {
	return FakeAddr
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

var FakeAddr = fakeAddr{}

type fakeAddr struct{}

func (a fakeAddr) Network() string {
	return "fakeAddr"
}

func (a fakeAddr) String() string {
	return "fakeAddr"
}
