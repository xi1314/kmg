package netTester

import (
	"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgTime"
	"io"
	"net"
	"time"
)

type DirectDialer func() (net.Conn, error)
type ListenNewer func() (net.Listener, error)

func RunTcpListenerDialerTest(listenerNewer ListenNewer,
	Dialer DirectDialer,
	debug bool) {
	writeRead(listenerNewer, Dialer, debug)
	time.Sleep(time.Microsecond)
	readWrite(listenerNewer, Dialer, debug)
	time.Sleep(time.Microsecond)
	readOnly(listenerNewer, Dialer, debug)
	time.Sleep(time.Microsecond)
	writeOnly(listenerNewer, Dialer, debug)
}

func writeRead(listenerNewer ListenNewer,
	Dialer DirectDialer,
	debug bool) {
	listener := runEchoServer(listenerNewer)
	defer listener.Close()

	toWrite := []byte("hello world")

	kmgTime.MustNotTimeout(func() {
		conn1, err := Dialer()
		mustNotError(err)

		if debug {
			conn1 = kmgNet.NewDebugConn(conn1, "writeRead")
		}
		defer conn1.Close()
		for i := 0; i < 2; i++ {
			_, err = conn1.Write(toWrite)
			mustNotError(err)
			time.Sleep(time.Microsecond)
			mustReadSame(conn1, toWrite)
			time.Sleep(time.Microsecond)
		}
		conn1.Close()
	}, time.Second)
	listener.Close()

}

func readWrite(listenerNewer ListenNewer, Dialer DirectDialer, debug bool) {
	listener := runEchoServer(listenerNewer)
	defer listener.Close()

	toWrite := []byte("hello world")

	kmgTime.MustNotTimeout(func() {
		conn1, err := Dialer()
		mustNotError(err)
		if debug {
			conn1 = kmgNet.NewDebugConn(conn1, "readWrite")
		}
		defer conn1.Close()
		for i := 0; i < 2; i++ {
			go func() {
				time.Sleep(time.Microsecond)
				_, err = conn1.Write(toWrite)
				mustNotError(err)
			}()
			mustReadSame(conn1, toWrite)
			time.Sleep(time.Microsecond)
		}
		conn1.Close()
	}, time.Second)

	listener.Close()
}

func readOnly(listenerNewer ListenNewer, Dialer DirectDialer, debug bool) {
	toWrite := []byte("hello world")
	listener := listenAccept(listenerNewer, func(c net.Conn) {
		defer c.Close()
		for i := 0; i < 2; i++ {
			_, err := c.Write(toWrite)
			mustNotError(err)
			time.Sleep(time.Microsecond)
		}
	})
	defer listener.Close()
	kmgTime.MustNotTimeout(func() {
		conn1, err := Dialer()
		mustNotError(err)
		if debug {
			conn1 = kmgNet.NewDebugConn(conn1, "readOnly")
		}
		defer conn1.Close()
		for i := 0; i < 2; i++ {
			mustReadSame(conn1, toWrite)
			time.Sleep(time.Microsecond)
		}
		conn1.Close()
	}, time.Second)

	listener.Close()
}

func writeOnly(listenerNewer ListenNewer, Dialer DirectDialer, debug bool) {
	toWrite := []byte("hello world")
	listener := listenAccept(listenerNewer, func(c net.Conn) {
		defer c.Close()
		for i := 0; i < 2; i++ {
			mustReadSame(c, toWrite)
			time.Sleep(time.Microsecond)
		}
	})
	defer listener.Close()
	kmgTime.MustNotTimeout(func() {
		conn1, err := Dialer()
		mustNotError(err)
		if debug {
			conn1 = kmgNet.NewDebugConn(conn1, "writeOnly")
		}
		defer conn1.Close()
		for i := 0; i < 2; i++ {
			_, err = conn1.Write(toWrite)
			mustNotError(err)
			time.Sleep(time.Microsecond)
		}
		conn1.Close()
	}, time.Second)

	listener.Close()
}

func mustNotError(err error) {
	if err != nil {
		panic(err)
	}
}

func mustReadSame(r io.Reader, toWrite []byte) {
	buf := make([]byte, len(toWrite))
	n, err := io.ReadAtLeast(r, buf, len(toWrite))
	mustNotError(err)
	if !bytes.Equal(buf[:n], toWrite) {
		panic(fmt.Errorf("read write data not match"))
		return
	}
}

func runEchoServer(listenerNewer ListenNewer) net.Listener {
	return listenAccept(listenerNewer, func(c net.Conn) {
		defer c.Close()
		_, err := io.Copy(c, c)
		if kmgNet.IsSocketCloseError(err) {
			return
		}
		mustNotError(err)
	})
}

func listenAccept(listenerNewer ListenNewer, handler func(c net.Conn)) net.Listener {
	listener, err := listenerNewer()
	mustNotError(err)
	go func() {
		for {
			c, err := listener.Accept()
			if kmgNet.IsSocketCloseError(err) {
				return
			}
			mustNotError(err)
			go handler(c)
		}
	}()
	return listener
}
