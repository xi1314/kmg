package main

import (
	"github.com/bronze1man/kmg/kmg/internal"
	"github.com/bronze1man/kmg/kmg/internal/InstallCmd"
	"github.com/bronze1man/kmg/kmg/internal/gitCmd"
	"github.com/bronze1man/kmg/kmg/internal/goCmd"
	"github.com/bronze1man/kmg/kmg/internal/serviceCmd"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
)

// kmg make upload
func main() {
	kmgHttp.AddCommandList()
	internal.AddCommandList()
	InstallCmd.AddCommandList()
	gitCmd.AddCommandList()
	goCmd.AddCommandList()
	serviceCmd.AddCommandList()

	internal.AddSelfInstallCommand()
	internal.AddSelfUpdate()

	kmgConsole.Main()
}
