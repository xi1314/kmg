package netTester

import (
	"time"

	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgTime"
)

//client 先写后读
func writeRead(listenerNewer ListenNewer,
	Dialer DirectDialer,
	debug bool) {
	kmgTime.MustNotTimeout(func() {
		listener := runEchoServer(listenerNewer)
		defer listener.Close()

		toWrite := []byte("hello world")

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
		listener.Close()
	}, time.Second)
}
