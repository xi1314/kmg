package kmgConsole

import (
	"os"
	"os/signal"
	"syscall"
)

//等用户按CTRL+c退出,或者等到被kill掉
func WaitForExit() {
	//不要在这个地方检查WaitForExit和AddExitAction一起使用,因为程序自身会进行调用
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}
