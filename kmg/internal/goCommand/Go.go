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
		Runner: goCommand,
	})

}

func goCommand() {
	args := []string{}
	args = append(args, os.Args[1:]...)
	cmd := kmgCmd.NewOsStdioCmd("go", args...)
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	err = kmgCmd.SetCmdEnv(cmd, "GOPATH", kmgc.GOPATHToString())
	kmgConsole.ExitOnErr(err)
	err = cmd.Run()
	kmgConsole.ExitOnErr(err)
}
