package main

import (
	_ "github.com/bronze1man/kmg/kmg/internal"
	//_ "github.com/bronze1man/kmg/kmg/internal/gitCommand"
	_ "github.com/bronze1man/kmg/kmg/internal/goCommand"
	"github.com/bronze1man/kmg/kmgConsole"
)

// kmg GoCrossCompile -v 'v1.3' github.com/bronze1man/kmg/kmg
func main() {
	kmgConsole.VERSION = "1.3"
	kmgConsole.Main()
}
