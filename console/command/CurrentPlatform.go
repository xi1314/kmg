package command

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgConsole"
	"runtime"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "CurrentPlatform",
		Desc:   "get current platform(from this binary)",
		Runner: runCurrentPlatform,
	})
}

func runCurrentPlatform() {
	fmt.Println(runtime.GOOS + "_" + runtime.GOARCH)
}
