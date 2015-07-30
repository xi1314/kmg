package netTester

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgTime"
)

//只dial,没有传输任何数据就挂掉了.
func dialOnly(listenerNewer ListenNewer, Dialer DirectDialer, debug bool) {
	kmgTime.MustNotTimeout(func() {
		listener := listenAccept(listenerNewer, func(c net.Conn) {
			time.Sleep(time.Microsecond) //防止出现reset
			c.Close()
		})
		defer listener.Close()
		//client 双方主动关闭
		for i := 0; i < 2; i++ {
			conn1, err := Dialer()
			mustNotError(err)
			if debug {
				conn1 = kmgNet.NewDebugConn(conn1, fmt.Sprintf("dialOnly 1 %d", i))
			}
			conn1.Close()
		}

		//client 被动关闭
		buf := make([]byte, 1024)
		for i := 0; i < 2; i++ {
			conn1, err := Dialer()
			mustNotError(err)
			if debug {
				conn1 = kmgNet.NewDebugConn(conn1, fmt.Sprintf("dialOnly 2 %d", i))
			}
			defer conn1.Close()
			_, err = conn1.Read(buf)
			if err != io.EOF {
				panic("remote close, local not get io.EOF")
			}
			conn1.Close()
		}

		listener.Close()
	}, time.Second)

	time.Sleep(time.Microsecond)
	kmgTime.MustNotTimeout(func() {
		listener := listenAccept(listenerNewer, func(c net.Conn) {
			buf := make([]byte, 1024)
			defer c.Close()
			_, err := c.Read(buf)
			if err != io.EOF {
				panic("local close, remote not get io.EOF")
			}
		})
		defer listener.Close()
		//client 主动关闭
		for i := 0; i < 2; i++ {
			conn1, err := Dialer()
			mustNotError(err)
			if debug {
				conn1 = kmgNet.NewDebugConn(conn1, fmt.Sprintf("dialOnly 1 %d", i))
			}
			time.Sleep(time.Microsecond) //防止出现reset
			conn1.Close()
		}

		listener.Close()
	}, time.Second)

}
