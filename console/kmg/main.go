package main

import (
	_ "github.com/bronze1man/kmg/console/command"
	"github.com/bronze1man/kmg/kmgConsole"
)

func main() {
	kmgConsole.VERSION = "0.1"
	kmgConsole.Main()
}
