package main

import (
	_ "github.com/bronze1man/kmg/kmg/internal"
	//_ "github.com/bronze1man/kmg/kmg/internal/gitCmd"
	_ "github.com/bronze1man/kmg/kmg/internal/InstallCmd"
	_ "github.com/bronze1man/kmg/kmg/internal/goCmd"
	_ "github.com/bronze1man/kmg/kmg/internal/serviceCmd"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
)

// kmg make BuildAndUploadKmg
func main() {
	kmgHttp.AddCommandList()
	kmgConsole.Main()
}
