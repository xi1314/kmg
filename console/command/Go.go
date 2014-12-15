package command

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
		Runner: newGoCommand(""),
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GoBuild",
		Desc:   "run go build in current project",
		Runner: newGoCommand("build"),
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GoRun",
		Desc:   "run go run in current project",
		Runner: newGoCommand("run"),
	})
}

func newGoCommand(command string) func() {
	return func() {
		args := []string{}
		if command != "" {
			args = append(args, command)
		}
		args = append(args, os.Args[1:]...)
		cmd := kmgCmd.NewOsStdioCmd("go", args...)
		kmgc, err := kmgConfig.LoadEnvFromWd()
		exitOnErr(err)
		err = kmgCmd.SetCmdEnv(cmd, "GOPATH", kmgc.GOPATHToString())
		exitOnErr(err)
		err = cmd.Run()
		exitOnErr(err)
	}
}
