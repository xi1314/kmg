package goCmd

import (
	"fmt"

	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
)

// use `kmg goenv` to setup current work project GOPATH
func GoEnvCmd() {
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	fmt.Println("export", "GOPATH="+kmgc.GOPATHToString()) //TODO 解决转义问题?
}
