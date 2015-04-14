package goCommand

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"os"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Go",
		Desc:   "run go command in current project",
		Runner: GoCommand,
	})

}

// run go command in current project
// 1.go build -i github.com/xxx/xxx use to get fastest speed of build.
// 2.try remove pkg directory if you found you change is ignore.
func GoCommand() {
	cmd := kmgCmd.CmdSlice(append([]string{"go"}, os.Args[1:]...)).GetExecCmd()
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	err = kmgCmd.SetCmdEnv(cmd, "GOPATH", kmgc.GOPATHToString())
	kmgConsole.ExitOnErr(err)
	err = cmd.Run()
	kmgConsole.ExitOnErr(err)
}
