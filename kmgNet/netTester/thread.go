package netTester

import (
	"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgTask"
	"github.com/bronze1man/kmg/kmgTime"
	"time"
)

//有9MB数据 3线程 一共10个任务
func thread(listenerNewer ListenNewer, Dialer DirectDialer, debug bool) {
	listener := runEchoServer(listenerNewer)
	defer listener.Close()

	task := kmgTask.NewLimitThreadTaskManager(3)
	content := bytes.Repeat([]byte("Hello world"), 1024*30)
	kmgTime.MustNotTimeout(func() {
		for i := 0; i < 10; i++ {
			task.AddFunc(func() {
				if debug {
					fmt.Println("[thread] start", i)
				}
				conn, err := Dialer()
				mustNotError(err)
				defer conn.Close()
				go func() {
					for i := 0; i < 3; i++ {
						_, err = conn.Write(content)
						mustNotError(err)
						time.Sleep(time.Microsecond)
					}
				}()
				for i := 0; i < 3; i++ {
					mustReadSame(conn, content)
				}
			})
		}
	}, 10*time.Second)
	task.Close()
}
