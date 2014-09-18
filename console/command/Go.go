package command

import (
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
)

type Go struct {
}

func (command *Go) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "Go", Short: "run go command in current project"}
}
func (command *Go) Execute(context *console.Context) (err error) {
	cmd := kmgCmd.NewStdioCmd(context, "go", context.Args[2:]...)
	kmgc, err := kmgConfig.LoadEnvFromWd()
	if err != nil {
		return
	}
	err = kmgCmd.SetCmdEnv(cmd, "GOPATH", kmgc.GOPATHToString())
	if err != nil {
		return err
	}
	return cmd.Run()
}
