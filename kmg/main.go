package main

import (
	"github.com/bronze1man/kmg/kmg/SubCommand"
	"github.com/bronze1man/kmg/kmg/SubCommand/InstallCmd"
	"github.com/bronze1man/kmg/kmg/SubCommand/gitCmd"
	"github.com/bronze1man/kmg/kmg/SubCommand/goCmd"
	"github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgSys"
)

// curl http://kmgtools.qiniudn.com/v1/installKmg.bash?v=$RANDOM | bash
// kmg make upload
func main() {
	kmgSys.RecoverPath()
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
