package main

import (
	_ "github.com/bronze1man/kmg/kmg/internal"
	//_ "github.com/bronze1man/kmg/kmg/internal/gitCommand"
	_ "github.com/bronze1man/kmg/kmg/internal/goCommand"
	_ "github.com/bronze1man/kmg/kmg/internal/serviceCommand"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
)

// kmg GoCrossCompile -v 'v1.6' github.com/bronze1man/kmg/kmg
func main() {
	kmgConsole.VERSION = "1.6"
	kmgHttp.AddCommand()
	kmgConsole.Main()
}
