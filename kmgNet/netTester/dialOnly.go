package netTester

import (
	"github.com/bronze1man/kmg/kmgTime"
	"io"
	"net"
	"time"
)

func dialOnly(listenerNewer ListenNewer, Dialer DirectDialer, debug bool) {
	kmgTime.MustNotTimeout(func() {
		listener := listenAccept(listenerNewer, func(c net.Conn) {
			c.Close()
		})
		defer listener.Close()

		for i := 0; i < 2; i++ {
			conn1, err := Dialer()
			mustNotError(err)
			conn1.Close()
		}

		buf := make([]byte, 1024)
		for i := 0; i < 2; i++ {
			conn1, err := Dialer()
			mustNotError(err)
			defer conn1.Close()
			_, err = conn1.Read(buf)
			if err != io.EOF {
				panic("remote close, local not get io.EOF")
			}
			conn1.Close()
		}

		listener.Close()
	}, time.Second)

}
