package kmgConsole

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForSIGINT() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
}
