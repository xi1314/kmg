package main

import (
	"github.com/bronze1man/kmg/kmgConsole"
    _ "github.com/bronze1man/kmg/console/kmg/internal"
    _ "github.com/bronze1man/kmg/console/kmg/internal/goCommand"
    _ "github.com/bronze1man/kmg/console/kmg/internal/gitCommand"
)

func main() {
	kmgConsole.VERSION = "1.0"
	kmgConsole.Main()
}
