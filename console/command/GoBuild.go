package command

import (
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
)

type GoBuild struct {
}

func (command *GoBuild) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoBuild", Short: "build some golang code in current project"}
}
func (command *GoBuild) Execute(context *console.Context) (err error) {
	args := append([]string{"build"}, context.Args[2:]...)
	cmd := kmgCmd.NewStdioCmd(context, "go", args...)
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
