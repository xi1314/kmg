package kmgNet

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgIo"
	"io"
	"net"
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

// a net.Conn with Reader Writer Closer override
type sizeConnRwcer struct {
	net.Conn
	name string
}

func (c sizeConnRwcer) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if err != nil {
		fmt.Println("[debugConn]", c.name, "Write len:", n, "err:", err)
	} else {
		fmt.Println("[debugConn]", c.name, "Write len:", n)
	}
	return n, err
}
func (c sizeConnRwcer) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if err != nil {
		fmt.Println("[debugConn]", c.name, "Read len:", n, "iLen", len(b), "err:", err)
	} else {
		fmt.Println("[debugConn]", c.name, "Read len:", n, "iLen", len(b))
	}
	return n, err
}
func (c sizeConnRwcer) Close() (err error) {
	err = c.Conn.Close()
	fmt.Println("[debugConn]", c.name, "Close err:", err)
	return err
}

func NewSizeDebugConn(conn net.Conn, name string) net.Conn {
	return sizeConnRwcer{
		Conn: conn,
		name: name,
	}
}

// a net.Conn with Reader Writer Closer override
type stringConnRwcer struct {
	net.Conn
	name string
}

func (c stringConnRwcer) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if err != nil {
		fmt.Println("[debugConn]", c.name, "Write len:", n, "err:", err)
	} else {
		fmt.Printf("[debugConn] [%s] Write len: %d content: %s<EndOfContent>\n", c.name, n, string(b[:n]))
	}
	return n, err
}
func (c stringConnRwcer) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if err != nil {
		fmt.Println("[debugConn]", c.name, "Read len:", n, "iLen", len(b), "err:", err)
	} else {
		fmt.Printf("[debugConn] [%s] Read len: %d iLen: %d content: %s<EndOfContent>\n", c.name, n, len(b), string(b[:n]))
	}
	return n, err
}
func (c stringConnRwcer) Close() (err error) {
	err = c.Conn.Close()
	fmt.Println("[debugConn]", c.name, "Close err:", err)
	return err
}

func NewStringDebugConn(conn net.Conn, name string) net.Conn {
	return stringConnRwcer{
		Conn: conn,
		name: name,
	}
}
