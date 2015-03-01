package kmgConsole

import (
	"os"
	"os/signal"
	"syscall"
)

//等用户按CTRL+c退出,或者等到被kill掉
func WaitForExit() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}
