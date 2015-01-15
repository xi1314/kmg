package netTester

import (
	"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgNet"
	"io"
	"net"
	"time"
)

func RunTcpListenerDialerTest(listener net.Listener, Dialer func() (net.Conn, error)) {
	defer listener.Close()
	go func() {
		for i := 0; ; i++ {
			c, err := listener.Accept()
			if kmgNet.IsSocketCloseError(err) {
				return
			}
			if err != nil {
				panic(err)
				return
			}
			go func(c net.Conn, i int) {
				defer c.Close()
				_, err := io.Copy(c, c)
				if kmgNet.IsSocketCloseError(err) {
					return
				}
				if err != nil {
					panic(err)
					return
				}
			}(c, i)
		}
	}()
	toWrite := []byte("hello world")
	buf := make([]byte, 8*1024)
	for i := 0; i < 3; i++ {
		//fmt.Println("request 0", i)
		func() {
			conn1, err := Dialer()
			if err != nil {
				panic(err)
				return
			}
			//fmt.Println("request 1", i, conn1)
			defer conn1.Close()
			for i := 0; i < 3; i++ {
				_, err = conn1.Write(toWrite)
				if err != nil {
					panic(err)
					return
				}
				n, err := conn1.Read(buf)
				if err != nil {
					panic(err)
					return
				}
				//fmt.Println("request 2", i)
				if !bytes.Equal(buf[:n], toWrite) {
					panic(fmt.Errorf("read write data not match"))
					return
				}
			}
			conn1.Close()
			time.Sleep(time.Microsecond)
		}()
	}
	listener.Close()
	time.Sleep(time.Microsecond)
}
