package main

import (
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgDebug"
	"github.com/bronze1man/kmg/third/kmgRadius"
)

func main() {
	// run the server in a new thread.
	// 在一个新线程里面一直运行服务器.
	kmgRadius.RunServer(":1812", []byte("sEcReT"), kmgRadius.Handler{
		Auth: func(username string) (password string, exist bool) {
			if username != "a" {
				return "", false
			}
			return "b", true
		},
		AcctStart: func(req kmgRadius.AcctRequest) {
			kmgDebug.Println("start", req)
		},
		AcctUpdate: func(req kmgRadius.AcctRequest) {
			kmgDebug.Println("update", req)
		},
		AcctStop: func(req kmgRadius.AcctRequest) {
			kmgDebug.Println("stop", req)
		},
	})
	// wait for the system sign or ctrl-c to close the process.
	// 等待系统信号或者ctrl-c 关闭进程
	kmgConsole.WaitForExit()
}
