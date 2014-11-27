package kmgNet

import (
	"io"
	"net"
	"github.com/bronze1man/kmg/kmgIo"
)


// a net.Conn with Reader Writer Closer override
type connRwcer struct {
	net.Conn
	rwc io.ReadWriteCloser
}

func (conn connRwcer) Read(p []byte) (n int, err error) {
	return conn.rwc.Read(p)
}
func (conn connRwcer) Write(p []byte) (n int, err error) {
	return conn.rwc.Write(p)
}
func (conn connRwcer) Close() (err error) {
	return conn.rwc.Close()
}

func NewDebugConn(conn net.Conn, name string) net.Conn {
	return connRwcer{
		Conn: conn,
		rwc:  kmgIo.NewDebugRwc(conn, name),
	}
}