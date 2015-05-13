package main

import (
	"github.com/bronze1man/kmg/kmg/SubCommand"
	"github.com/bronze1man/kmg/kmg/SubCommand/InstallCmd"
	"github.com/bronze1man/kmg/kmg/SubCommand/gitCmd"
	"github.com/bronze1man/kmg/kmg/SubCommand/goCmd"
	"github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
)

// kmg make upload
func main() {
	kmgHttp.AddCommandList()
	SubCommand.AddCommandList()
	InstallCmd.AddCommandList()
	gitCmd.AddCommandList()
	goCmd.AddCommandList()
	serviceCmd.AddCommandList()

	SubCommand.AddSelfInstallCommand()
	SubCommand.AddSelfUpdate()

	kmgConsole.Main()
}
