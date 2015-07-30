package netTester

import (
	"bytes"
	"sync"
	"time"

	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgTime"
)

//一次写入大量数据(超过mtu即可)
func bigWrite(listenerNewer ListenNewer, Dialer DirectDialer, debug bool) {
	kmgTime.MustNotTimeout(func() {
		listener := runEchoServer(listenerNewer)
		defer listener.Close()

		toWrite := bytes.Repeat([]byte{0}, 1024*100)

		func() {
			//先写后读
			wg := sync.WaitGroup{}
			conn1, err := Dialer()
			mustNotError(err)
			if debug {
				conn1 = kmgNet.NewDebugConn(conn1, "bigWrite1")
			}
			defer conn1.Close()
			wg.Add(1)
			go func() {
				for i := 0; i < 2; i++ {
					mustReadSame(conn1, toWrite)
				}
				wg.Done()
			}()
			for i := 0; i < 2; i++ {
				_, err = conn1.Write(toWrite)
				mustNotError(err)
				time.Sleep(time.Microsecond)
			}
			wg.Wait()
			conn1.Close()
		}()

		func() {
			//先读后写
			conn1, err := Dialer()
			mustNotError(err)
			if debug {
				conn1 = kmgNet.NewDebugConn(conn1, "bigWrite2")
			}
			defer conn1.Close()
			go func() {
				for i := 0; i < 2; i++ {
					_, err := conn1.Write(toWrite)
					mustNotError(err)
					time.Sleep(time.Microsecond)
				}
			}()
			for i := 0; i < 2; i++ {
				mustReadSame(conn1, toWrite)
			}
			conn1.Close()
		}()

		listener.Close()
	}, time.Second)
}
