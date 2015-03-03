package main

import (
	_ "github.com/bronze1man/kmg/kmg/internal"
	_ "github.com/bronze1man/kmg/kmg/internal/gitCommand"
	_ "github.com/bronze1man/kmg/kmg/internal/goCommand"
	"github.com/bronze1man/kmg/kmgConsole"
)

func main() {
	kmgConsole.VERSION = "1.0"
	kmgConsole.Main()
}
