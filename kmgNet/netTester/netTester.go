package netTester

import (
	"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgNet"
	"io"
	"net"
	"time"
)

type DirectDialer func() (net.Conn, error)
type ListenNewer func() (net.Listener, error)

/*
func RunTcpTestWithListenAddr(listenAddr string,Dialer DirectDialer,debug bool){
	return RunTcpListenerDialerTest(
		func()(net.Listener,error){
			return net.Listen("tcp",listenAddr)
		},
		Dialer,debug)
}
*/
func RunTcpTestWithNetDialAndNetListener(listenAddr string, dialAddr string, debug bool) {
	RunTcpListenerDialerTest(
		func() (net.Listener, error) {
			return net.Listen("tcp", listenAddr)
		},
		func() (net.Conn, error) {
			return net.Dial("tcp", dialAddr)
		}, debug)
}

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
	time.Sleep(time.Microsecond)
	thread(listenerNewer, Dialer, debug)
	time.Sleep(time.Microsecond)
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
