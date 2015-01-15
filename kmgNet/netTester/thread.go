package netTester

import (
	//"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgTask"
	"github.com/bronze1man/kmg/kmgTime"
	"time"
)

//有9MB数据 3线程 一共10个任务
func thread(listenerNewer ListenNewer, Dialer DirectDialer, debug bool) {
	kmgTime.MustNotTimeout(func() {
		listener := runEchoServer(listenerNewer)
		defer listener.Close()

		task := kmgTask.NewLimitThreadTaskManager(10)
		content := []byte("Hello world")
		for i := 0; i < 30; i++ {
			task.AddFunc(func() {
				if debug {
					fmt.Println("[thread] start", i)
				}
				conn, err := Dialer()
				mustNotError(err)
				defer conn.Close()
				go func() {
					for i := 0; i < 10; i++ {
						_, err = conn.Write(content)
						mustNotError(err)
						time.Sleep(time.Microsecond)
					}
				}()
				for i := 0; i < 10; i++ {
					mustReadSame(conn, content)
				}
			})
		}
		task.Close()
	}, 5*time.Second) //在rt很高的环境下可能会花费较长时间
}
